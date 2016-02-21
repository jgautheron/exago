package leveldb

import (
	"gopkg.in/vmihailenco/msgpack.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/config"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	db *ldb.DB
)

// connect initiates a LevelDB connection.
func connect() *ldb.DB {
	var err error
	db, err = ldb.OpenFile(config.Get("DatabasePath"), &opt.Options{
		Filter: filter.NewBloomFilter(10),
	})
	if err != nil {
		log.Fatal(err)
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

func FindForRepositoryCmd(key []byte) ([]byte, error) {
	b, err := Instance().Get(key, nil)
	if err != nil {
		return nil, err
	}
	var out []byte
	err = msgpack.Unmarshal(b, &out)
	return out, err
}

func Save(key, data []byte) error {
	b, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	err = Instance().Put(key, b, nil)
	return err
}
