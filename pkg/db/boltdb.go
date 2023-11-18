package db

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"time"

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
			slog.Error("boltdb", "path", db.path, "bucket", db.name, "msg", "unable to gracefully stop boltdb", "error", shutdownCtx.Err())
			return
		}
	}()
	slog.Info("boltdb", "path", db.path, "bucket", db.name, "status", "stopping")
	if err := db.db.Close(); err != nil {
		slog.Error("boltdb", "path", db.path, "bucket", db.name, "msg", "unable to gracefully stop boltdb", "error", err)
		return
	}
	slog.Info("boltdb", "path", db.path, "bucket", db.name, "status", "stopped")
}

func (db *BoltDB) Load(ctx context.Context, key []byte, result any) error {
	return db.db.View(func(tx *bolt.Tx) error {
		slog.Debug("loading", "bucket", db.name, "type", fmt.Sprintf("%T", result))
		bucket := tx.Bucket(db.bucket)
		if bucket == nil {
			return ErrDBBucketNotFound
		}
		data := bucket.Get(key)
		if data == nil {
			return ErrDBEntryNotFound
		}
		return gob.NewDecoder(bytes.NewBuffer(data)).Decode(result)
	})
}

func (db *BoltDB) Save(ctx context.Context, key []byte, data any) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		slog.Debug("saving", "bucket", db.name, "data", data, "type", fmt.Sprintf("%T", data))
		bucket, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		err = gob.NewEncoder(buf).Encode(data)
		if err != nil {
			return err
		}

		return bucket.Put(key, buf.Bytes())
	})
}
