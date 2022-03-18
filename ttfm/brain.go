package ttfm

import (
	"github.com/andreapavoni/ttfm_bot/db"
)

type Brain struct {
	Db *db.Repo
}

func NewBrain(dir string) *Brain {
	db, err := db.New(dir, &db.Options{})
	if err != nil {
		panic(err)
	}

	return &Brain{Db: db}
}

func (b *Brain) Get(bucket, key string, value interface{}) error {
	if err := b.Db.Read(bucket, key, value); err != nil {
		return err
	}
	return nil
}

func (b *Brain) Put(bucket, key string, value interface{}) error {
	return b.Db.Write(bucket, key, value)
}
