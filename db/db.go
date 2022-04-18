package db

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/neosouler7/GObserver/tg"
)

var (
	db   *bolt.DB
	once sync.Once
)

const (
	dbName     = "gobserver"
	dataBucket = "data"
)

type data struct {
	createdAt     []byte
	lastUpdatedAt []byte
	payload       []byte
}

func DB() *bolt.DB {
	if db == nil {
		once.Do(func() {
			dbPointer, err := bolt.Open(fmt.Sprintf("%s.db", dbName), 0600, nil)
			db = dbPointer
			tg.HandleErr(err)

			err = db.Update(func(tx *bolt.Tx) error {
				_, err = tx.CreateBucketIfNotExists([]byte(dataBucket))
				tg.HandleErr(err)
				return err
			})
			tg.HandleErr(err)
		})
	}
	return db
}

func Start() {
	fmt.Println("db called")
	db = DB()
}
