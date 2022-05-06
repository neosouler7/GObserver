// Package db deals with database to store informations.
package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

const (
	dbName        = "gobserver"
	mold          = "mold"
	moldBucket    = "moldBucket"
	createdAt     = "createdAt"
	LastUpdatedAt = "lastUpdatedAt"
)

type MoldStruct struct {
	Payload []byte
}

// Returns of bucket exists.
func bucketExists() bool {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(moldBucket))
		if b == nil {
			return nil
		}
		return errors.New("bucket exists")
	})
	if err == nil {
		return false
	}
	return true
}

// Create new bucket.
func createBucket() error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(moldBucket))
		if err != nil {
			log.Fatalln(err)
		}
		return err
	})
	return err
}

// Clears existing bucket.
func clearBucket() error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(moldBucket))
		if err != nil {
			log.Fatalln(err)
		}
		return err
	})
	return err
}

// Saves checkpoints.
func SaveCheckPoint(checkpoint string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(moldBucket))
		checkPointAsByte := []byte(strconv.Itoa(int(time.Now().Unix())))
		err := bucket.Put([]byte(checkpoint), checkPointAsByte)
		if err != nil {
			log.Fatalln(err)
		}
		return err
	})
	return err
}

// Loads checkpoints.
func GetCheckPoint(checkpoint string) []byte {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(moldBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

// Update mold informations.
func UpdateMold(data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(moldBucket))
		err := bucket.Put([]byte(mold), data)
		return err
	})
	if err != nil {
		log.Fatalln(err)
	}
}

// Creates new bucket. (if exists, delete first)
func InitBucket() error {
	if bucketExists() {
		err := clearBucket()
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Cleared %s bucket.", dbName)
	}
	err := createBucket()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Created new %s bucket.", dbName)

	err = SaveCheckPoint(createdAt)
	return err
}

// Starts db package.
func Start() {
	dbPointer, err := bolt.Open(fmt.Sprintf("%s.db", dbName), 0600, nil)
	db = dbPointer
	if err != nil {
		log.Fatalln(err)
	}

	err = InitBucket()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%s db synced.", dbName)
}
