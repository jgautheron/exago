package rank

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/jgautheron/exago-service/redis"
)

var (
	errMissingData = errors.New("Not enough data to calculate the rank")

	required = map[string]bool{
		"loc":     true,
		"imports": true,
		"test":    true,
	}
)

type Rank struct {
	conn       redigo.Conn
	repository string

	// Unmarshaled data retrieved from the DB
	loc     map[string]int
	imports []string
	tests   testRunner

	// Output
	Score Score
}

func New() *Rank {
	return &Rank{
		conn: redis.GetConn(),
	}
}

// SetRepository sets the repository name.
func (rk *Rank) SetRepository(repo string) {
	rk.repository = repo
}

func (rk *Rank) Get() (interface{}, error) {
	data, err := rk.loadData()
	if err != nil {
		return nil, err
	}

	if err := rk.deserialize(data); err != nil {
		return nil, err
	}

	rk.calcScore()
	return rk, err
}

func (rk *Rank) loadData() (map[string][]byte, error) {
	data := map[string][]byte{}
	for idfr := range required {
		o, err := rk.conn.Do("HGET", rk.repository, idfr)
		if o == nil {
			return nil, errMissingData
		}
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		if _, err := buf.Write(o.([]byte)); err != nil {
			return nil, err
		}
		rd, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(rd)
		if err != nil {
			return nil, err
		}

		data[idfr] = b
	}

	return data, nil
}

func (rk *Rank) deserialize(data map[string][]byte) error {
	if err := json.Unmarshal(stripEnvelope(data["loc"]), &rk.loc); err != nil {
		return err
	}
	if err := json.Unmarshal(stripEnvelope(data["imports"]), &rk.imports); err != nil {
		return err
	}
	if err := json.Unmarshal(stripEnvelope(data["test"]), &rk.tests); err != nil {
		return err
	}
	return nil
}

// stripEnvelope removes the JSend enveloppe to simplify the JSON processing.
func stripEnvelope(data []byte) []byte {
	d := string(data)
	d = strings.Replace(d, `{"data":`, "", 1)
	d = strings.Replace(d, `,"status":"success"}`, "", 1)
	return []byte(d)
}

type testRunner struct {
	Checklist struct {
		Failed []struct {
			Category string `json:"Category"`
			Desc     string `json:"Desc"`
			Name     string `json:"Name"`
		} `json:"Failed"`
		Passed []struct {
			Category string `json:"Category"`
			Desc     string `json:"Desc"`
			Name     string `json:"Name"`
		} `json:"Passed"`
	} `json:"checklist"`
	Packages []struct {
		Coverage      float64 `json:"coverage"`
		ExecutionTime float64 `json:"execution_time"`
		Name          string  `json:"name"`
		Success       bool    `json:"success"`
	} `json:"packages"`
}
