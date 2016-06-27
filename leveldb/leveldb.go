// Package leveldb is a thin simplicity layer over the LevelDB database.
package leveldb

import (
	"fmt"

	. "github.com/exago/svc/config"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	// db corresponds to the database connection.
	// Internal functions directly use db instead of instance() so remember to
	// always Init() the connection before.
	db *ldb.DB
)

func Init() (err error) {
	_, err = instance()
	return err
}

func FindForRepositoryCmd(key []byte) (b []byte, err error) {
	b, err = db.Get(key, nil)
	if err != nil {
		if err == ldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

func FindAllForRepository(prefix []byte) (map[string][]byte, error) {
	m := map[string][]byte{}
	iter := db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	for iter.Next() {
		// Get the key
		key := iter.Key()
		ckey := make([]byte, len(key))
		copy(ckey, key)

		// Get the value
		val := iter.Value()
		cval := make([]byte, len(val))
		copy(cval, val)

		m[string(ckey)] = cval
	}
	return m, iter.Error()
}

func DeleteAllMatchingPrefix(prefix []byte) error {
	iter := db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	for iter.Next() {
		// Get the key
		key := iter.Key()
		ckey := make([]byte, len(key))
		copy(ckey, key)

		if err := db.Delete(ckey, nil); err != nil {
			return err
		}
	}
	return nil
}

func Save(key []byte, data []byte) error {
	return db.Put(key, data, nil)
}

func Get(key []byte) ([]byte, error) {
	b, err := db.Get(key, nil)
	if err != nil {
		if err == ldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return b, err
}

// connect initiates a LevelDB connection.
func connect() (*ldb.DB, error) {
	var err error
	db, err = ldb.OpenFile(Config.DatabasePath, &opt.Options{})
	if err != nil {
		return nil, fmt.Errorf("An error occurred while trying to open the DB: %s", err)
	}
	return db, nil
}

// Instance returns the current LevelDB instance.
func instance() (*ldb.DB, error) {
	if db != nil {
		return db, nil
	}
	return connect()
}
