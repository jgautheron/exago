package leveldb

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/exago/svc/config"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	db *ldb.DB
)

// connect initiates a LevelDB connection.
// connect cannot be called concurrently or it will crash with the infamous error:
// "resource temporarily unavailable", which is why the middleware initDB
// ensures the DB connection is instantiated before anything happens.
func connect() *ldb.DB {
	var err error
	db, err = ldb.OpenFile(Config.DatabasePath, &opt.Options{
		Filter: filter.NewBloomFilter(10),
	})
	if err != nil {
		log.Fatalf("An error occurred while trying to open the DB: %s", err)
	}
	return db
}

// Instance returns the current LevelDB instance.
func Instance() *ldb.DB {
	if db != nil {
		return db
	}
	return connect()
}

func FindForRepositoryCmd(key []byte) (b []byte, err error) {
	b, err = Instance().Get(key, nil)
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
	iter := Instance().NewIterator(util.BytesPrefix(prefix), nil)
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
	iter := Instance().NewIterator(util.BytesPrefix(prefix), nil)
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
	return Instance().Put(key, data, nil)
}
