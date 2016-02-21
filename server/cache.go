package server

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/redis"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func getCacheIdentifier(r *http.Request) string {
	ps := context.Get(r, "params").(httprouter.Params)
	sp := strings.Split(r.URL.String(), "/")
	switch action := sp[4]; {
	case ps.ByName("linter") != "":
		return action + ":" + ps.ByName("linter")
	case ps.ByName("path") != "":
		return action + ":" + ps.ByName("path")
	default:
		return action
	}
}

func cacheOutput(r *http.Request, output []byte) {
	const timeout = 2 * 3600

	ps := context.Get(r, "params").(httprouter.Params)
	lgr := context.Get(r, "logger").(*log.Entry)
	c := redis.GetConn()

	idfr := getCacheIdentifier(r)
	k := fmt.Sprintf("%s/%s/%s", ps.ByName("registry"), ps.ByName("owner"), ps.ByName("repository"))
	if _, err := c.Do("HMSET", k, idfr, output); err != nil {
		lgr.Error(err)
		return
	}
	// if _, err := c.Do("EXPIRE", k, timeout); err != nil {
	// 	lgr.Error(err)
	// 	return
	// }
}
