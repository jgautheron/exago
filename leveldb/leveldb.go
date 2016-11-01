// Package leveldb is a thin simplicity layer over the LevelDB database.
package leveldb

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	. "github.com/hotolab/exago-svc/config"
	"github.com/hotolab/exago-svc/repository/model"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	logger = log.WithField("prefix", "database")

	// Make sure it satisfies the interface.
	_ model.Database = (*LevelDB)(nil)
)

type LevelDB struct {
	conn *ldb.DB
}

func New() (model.Database, error) {
	conn, err := ldb.OpenFile(Config.DatabasePath, &opt.Options{})
	if err != nil {
		return nil, fmt.Errorf("An error occurred while trying to open the DB: %s", err)
	}
	return &LevelDB{conn}, nil
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
