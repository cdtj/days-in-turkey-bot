package db

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log/slog"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
	"github.com/boltdb/bolt"
)

type BoltDB struct {
	name   string
	path   string
	bucket []byte
	db     *bolt.DB
}

func NewBoltDB(path, bucket string) *BoltDB {
	return &BoltDB{
		path:   path,
		name:   bucket,
		bucket: []byte(bucket),
	}
}

func (db *BoltDB) Serve(ctx context.Context) error {
	bdb, err := bolt.Open(db.path+".bolt", 0600, nil)
	if err != nil {
		return err
	}
	db.db = bdb
	<-ctx.Done()
	return nil
}

func (db *BoltDB) Shutdown(ctx context.Context) {
	shutdownCtx, shutdownStopCtx := context.WithTimeout(ctx, 15*time.Second)
	defer shutdownStopCtx()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			slog.Error("boltdb", "path", db.path, "bucket", db.bucket, "msg", "unable to gracefully stop boltdb", "error", shutdownCtx.Err())
			return
		}
	}()
	slog.Info("boltdb", "path", db.path, "bucket", db.bucket, "status", "stopping")
	if err := db.db.Close(); err != nil {
		slog.Error("boltdb", "path", db.path, "bucket", db.bucket, "msg", "unable to gracefully stop boltdb", "error", err)
		return
	}
	slog.Info("boltdb", "path", db.path, "bucket", db.bucket, "status", "stopped")
}

func (db *BoltDB) Load(ctx context.Context, key any) (any, error) {
	var result interface{}
	switch db.name {
	case "users":
		result = new(model.User)
		key = keyToByte(key)
	default:
		return nil, ErrDBUnknownEntity
	}
	err := db.db.View(func(tx *bolt.Tx) error {
		slog.Info("loading", "bucket", db.bucket, "type", fmt.Sprintf("%T", result))
		bucket := tx.Bucket(db.bucket)
		if bucket == nil {
			return ErrDBBucketNotFound
		}
		data := bucket.Get(key.([]byte))
		if data == nil {
			return ErrDBEntryNotFound
		}
		return gob.NewDecoder(bytes.NewBuffer(data)).Decode(result)
	})
	return result, err
}

func (db *BoltDB) Save(ctx context.Context, key any, data any) error {
	switch db.name {
	case "users":
		data = data.(*model.User)
		key = keyToByte(key)
	default:
		return ErrDBUnknownEntity
	}
	return db.db.Update(func(tx *bolt.Tx) error {
		slog.Info("saving", "bucket", db.bucket, "data", data, "type", fmt.Sprintf("%T", data))
		bucket, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		err = gob.NewEncoder(buf).Encode(data)
		if err != nil {
			return err
		}

		return bucket.Put(key.([]byte), buf.Bytes())
	})
}

func keyToByte(id interface{}) []byte {
	switch v := id.(type) {
	case int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(v))
		return b
	}
	return nil
}
