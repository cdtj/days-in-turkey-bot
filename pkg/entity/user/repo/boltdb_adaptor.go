package repo

import (
	"context"
	"encoding/binary"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/model"
)

type UserBoltDBAdaptor struct {
	boltdb *db.BoltDB
}

func NewUserBoltDBAdaptor(boltdb *db.BoltDB) *UserBoltDBAdaptor {
	return &UserBoltDBAdaptor{
		boltdb: boltdb,
	}
}

func (a *UserBoltDBAdaptor) Load(ctx context.Context, key any) (any, error) {
	result := new(model.User)
	if err := a.boltdb.Load(ctx, keyToByte(key.(int64)), result); err != nil {
		return nil, err
	}
	return result, nil
}
func (a *UserBoltDBAdaptor) Save(ctx context.Context, key any, data any) error {
	return a.boltdb.Save(ctx, keyToByte(key.(int64)), data)
}

func keyToByte(id int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(id))
	return b
}
