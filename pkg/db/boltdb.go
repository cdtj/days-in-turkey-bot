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
	path   string
	bucket string
	db     *bolt.DB
}

func NewBoltDB(path, bucket string) *BoltDB {
	return &BoltDB{
		path:   path,
		bucket: bucket,
	}
}

func (db *BoltDB) Serve(ctx context.Context) error {
	bdb, err := bolt.Open(db.path, 0600, nil)
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
			slog.Error("unable to gracefully stop boltdb", "error", shutdownCtx.Err())
			return
		}
	}()

	slog.Info("closing boltdb", "path", db.path, "bucket", db.bucket)
	if err := db.db.Close(); err != nil {
		slog.Error("unable to gracefully stop boltdb", "error", err)
	}
}

func (db *BoltDB) Load(ctx context.Context, id interface{}) (interface{}, error) {
	var result interface{}
	switch db.bucket {
	case "users":
		result = new(model.User)
	}
	err := db.db.View(func(tx *bolt.Tx) error {
		slog.Info("loading", "bucket", db.bucket, "type", fmt.Sprintf("%T", result))
		bucket := tx.Bucket([]byte(db.bucket))
		if bucket == nil {
			return ErrDBBucketNotFound
		}
		data := bucket.Get(keyToByte(id))
		if data == nil {
			return ErrDBEntryNotFound
		}
		return gob.NewDecoder(bytes.NewBuffer(data)).Decode(result)
	})
	return result, err
}

func (db *BoltDB) Save(ctx context.Context, id interface{}, data interface{}) error {
	switch db.bucket {
	case "users":
		data = data.(*model.User)
	}
	return db.db.Update(func(tx *bolt.Tx) error {
		slog.Info("saving", "bucket", db.bucket, "data", data, "type", fmt.Sprintf("%T", data))
		bucket, err := tx.CreateBucketIfNotExists([]byte(db.bucket))
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		err = gob.NewEncoder(buf).Encode(data)
		if err != nil {
			return err
		}

		return bucket.Put((keyToByte(id)), buf.Bytes())
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
