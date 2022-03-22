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

func (b *Brain) Get(key string, value interface{}) error {
	if err := b.Db.Read("brain", key, value); err != nil {
		return err
	}
	return nil
}

func (b *Brain) Put(key string, value interface{}) error {
	return b.Db.Write("brain", key, value)
}
