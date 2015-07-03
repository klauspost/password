// Copyright 2015, Klaus Post, see LICENSE for details.

// Driver for BoltDB
package boltpw

import (
	"github.com/boltdb/bolt"
)

type BoltDB struct {
	DB     *bolt.DB
	Bucket []byte
}

// New will return a new Database interface that read/writes
// items to a single bucket.
func New(db *bolt.DB, bucket string) (*BoltDB, error) {
	b := &BoltDB{DB: db, Bucket: []byte(bucket)}

	// Execute several commands within a write transaction.
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(b.Bucket))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return b, nil

}

// Has satisfies the password.DB interface
func (b BoltDB) Has(s string) (bool, error) {
	var res bool
	err := b.DB.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(b.Bucket).Get([]byte(s))
		res = v != nil
		return nil
	})
	return res, err
}

// Has satisfies the password.DbWriter interface.
// It writes a single password to the database
func (b BoltDB) Add(s string) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(b.Bucket)
		b.Put([]byte(s), []byte{})
		return nil
	})
}

// AddMultiple satisfies the password.BulkWriter interface.
// It writes a number of passwords to the database
func (b BoltDB) AddMultiple(s []string) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(b.Bucket)
		for _, key := range s {
			err := b.Put([]byte(key), []byte{})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
