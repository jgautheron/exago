// Package leveldb is a thin simplicity layer over the LevelDB database.
package leveldb

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	. "github.com/hotolab/exago-svc/config"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	db     LevelDB
	once   sync.Once
	logger = log.WithField("prefix", "leveldb")

	// Make sure it satisfies the interface.
	_ Database = (*LevelDB)(nil)
)

type LevelDB struct {
	conn *ldb.DB
}

func GetInstance() LevelDB {
	once.Do(func() {
		conn, err := ldb.OpenFile(Config.DatabasePath, &opt.Options{})
		if err != nil {
			logger.Fatalf("An error occurred while trying to open the DB: %s", err)
		}
		db = LevelDB{conn}
	})
	return db
}

func (l LevelDB) DeleteAllMatchingPrefix(prefix []byte) error {
	iter := l.conn.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()
	for iter.Next() {
		// Get the key
		key := iter.Key()
		ckey := make([]byte, len(key))
		copy(ckey, key)

		if err := l.conn.Delete(ckey, nil); err != nil {
			return err
		}
	}
	return nil
}

func (l LevelDB) Put(key []byte, data []byte) error {
	return l.conn.Put(key, data, nil)
}

func (l LevelDB) Get(key []byte) ([]byte, error) {
	return l.conn.Get(key, nil)
}

func (l LevelDB) Delete(key []byte) error {
	return l.conn.Delete(key, nil)
}

type Database interface {
	DeleteAllMatchingPrefix(prefix []byte) error
	Delete(key []byte) error
	Put(key []byte, data []byte) error
	Get(key []byte) ([]byte, error)
}
