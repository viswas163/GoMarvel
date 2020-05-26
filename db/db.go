package db

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// Store : Store model for DB
type Store struct {
	db *bolt.DB
}

var (
	// ErrNotFound : No key found in db for get()
	ErrNotFound = errors.New("store: key not found")
	// ErrBadValue : Bad input value for put()
	ErrBadValue = errors.New("store: bad value")

	bucketName = []byte("MarvelDB")
	dbPath     = "marvel.db"
	dbInit     = false
	dbStore    *Store
)

// GetInstance : Returns the instance of the DB
func GetInstance() *Store {
	return dbStore
}

// Open : Opens a DB Store file
func Open(path string) (*Store, error) {
	if path == "" {
		path = dbPath
	}
	opts := &bolt.Options{
		Timeout: 50 * time.Millisecond,
	}

	// Open the bolt DB in path
	db, err := bolt.Open(path, 0640, opts)
	if err != nil {
		fmt.Println("Error opening Bolt DB")
		return nil, err
	}

	// Add a bucket if not exists already
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		fmt.Println("Error creating bucket")
		return nil, err
	}
	dbStore = &Store{db: db}
	dbInit = true
	return dbStore, nil
}

// Put : Adds an entry to the store with the given key and value
func (store *Store) Put(key string, value interface{}) error {
	if value == nil {
		return ErrBadValue
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return store.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), buf.Bytes())
	})
}

// Get : Returns the entry from store to the value reference with given key
func (store *Store) Get(key string, value interface{}) error {
	return store.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, v := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		} else if value == nil {
			return nil
		} else {
			d := gob.NewDecoder(bytes.NewReader(v))
			return d.Decode(value)
		}
	})
}

// Delete : Deletes the entry with given key
func (store *Store) Delete(key string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()
		if k, _ := c.Seek([]byte(key)); k == nil || string(k) != key {
			return ErrNotFound
		}
		return c.Delete()
	})
}

// Rename : Renames the entry of given key with new key
func (store *Store) Rename(key string, newKey string) error {
	var t []byte
	if err := store.Get(key, &t); err != nil {
		return err
	}
	if len(t) <= 0 {
		return nil
	}
	if err := store.Delete(key); err != nil {
		return err
	}
	if err := store.Put(newKey, t); err != nil {
		return err
	}
	return nil
}

// Close : closes the Store DB file
func (store *Store) Close() error {
	if !dbInit {
		return nil
	}
	return store.db.Close()
}
