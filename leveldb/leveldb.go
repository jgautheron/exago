// Package leveldb is a thin simplicity layer over the LevelDB database.
package leveldb

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	. "github.com/exago/svc/config"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	db   LevelDB
	once sync.Once

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
			log.Fatal("An error occurred while trying to open the DB: %s", err)
		}
		db = LevelDB{conn}
	})
	return db
}

func (l LevelDB) FindForRepositoryCmd(key []byte) (b []byte, err error) {
	b, err = l.conn.Get(key, nil)
	if err != nil {
		// if err == ldb.ErrNotFound {
		// 	return nil, nil
		// }
		return nil, err
	}
	return
}

func (l LevelDB) FindAllForRepository(prefix []byte) (map[string][]byte, error) {
	m := map[string][]byte{}
	iter := l.conn.NewIterator(util.BytesPrefix(prefix), nil)
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
	b, err := l.conn.Get(key, nil)
	if err != nil {
		if err == ldb.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return b, err
}

type Database interface {
	FindForRepositoryCmd(key []byte) (b []byte, err error)
	FindAllForRepository(prefix []byte) (map[string][]byte, error)
	DeleteAllMatchingPrefix(prefix []byte) error
	Put(key []byte, data []byte) error
	Get(key []byte) ([]byte, error)
}
