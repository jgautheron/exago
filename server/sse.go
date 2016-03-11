package server

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

type RepoEvent struct {
	eventType string
	data      interface{}
}

func (r RepoEvent) Id() string    { return r.eventType }
func (r RepoEvent) Event() string { return "message" }
func (r RepoEvent) Data() string {
	out, err := json.Marshal(r.data)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(out)
}

func (rc *RepositoryChecker) publishError(et string, err error) {
	log.Errorln("publishError", err)
	rc.srv.Publish([]string{rc.repository.Name}, RepoEvent{
		et, err,
	})
}

func (rc *RepositoryChecker) publishMessage(et string, message interface{}) {
	log.Infoln("publishMessage", message)
	rc.srv.Publish([]string{rc.repository.Name}, RepoEvent{
		et, message,
	})
}
