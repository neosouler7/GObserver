// Package db deals with database to store informations.
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

// Create new bucket.
func createBucket() {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
		tg.HandleErr(err)
		return err
	})
	tg.HandleErr(err)
}

// Reset bucket.
func ResetBucket() {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	tg.HandleErr(err)
}

// Restores existing database.
func restore() {}

// Update database informations.
func Update() {}

// Handles database by conditions.
func Start() {
	fmt.Println("db called")

	dbPointer, err := bolt.Open(fmt.Sprintf("%s.db", dbName), 0600, nil)
	db = dbPointer
	tg.HandleErr(err)

	// if not exists, "create" & "reset"
	if db == nil {
		createBucket()
		ResetBucket()
	} else {
		// lastUpdatedAt :=
		// if exists and valid, "restore"
		// if exists but expired, "Reset"

	}
}
