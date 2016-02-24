package leveldb

import (
	"strings"

	"gopkg.in/vmihailenco/msgpack.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/config"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	db *ldb.DB
)

// connect initiates a LevelDB connection.
func connect() *ldb.DB {
	var err error
	db, err = ldb.OpenFile(config.Values.DatabasePath, &opt.Options{
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

func FindAllForRepository(prefix []byte) (map[CodeInfoKey][]byte, error) {
	m := map[CodeInfoKey][]byte{}
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

		out := []byte{}
		if err := msgpack.Unmarshal(cval, &out); err != nil {
			return nil, err
		}

		ks := extractKey(ckey)
		m[ks] = out
	}
	return m, iter.Error()
}

func Save(key, data []byte) error {
	b, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	err = Instance().Put(key, b, nil)
	return err
}

func extractKey(key []byte) CodeInfoKey {
	sp := strings.Split(string(key), "-")
	return CodeInfoKey{
		sp[1],
		sp[2],
		sp[3],
	}
}

type CodeInfoKey struct {
	Branch, Linter, Type string
}
