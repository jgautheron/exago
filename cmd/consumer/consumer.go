package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jgautheron/exago/internal/database/firestore"
	"github.com/jgautheron/exago/internal/eventpub"
	exago "github.com/jgautheron/exago/pkg"
	"github.com/jgautheron/exago/pkg/analysis/task"
	"github.com/sirupsen/logrus"
)

// PubSubMessage is the payload of a Pub/Sub event
type PubSubMessage struct {
	Message struct {
		ID         string            `json:"id"`
		Data       []byte            `json:"data,omitempty"`
		Attributes map[string]string `json:"attributes"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type Consumer struct {
	db *firestore.Firestore
}

// New creates new Consumer
func New(db *firestore.Firestore) (*Consumer, error) {
	return &Consumer{db}, nil
}

// ProcessRecord handles data from a single record
func (c *Consumer) ProcessRecord(ctx context.Context, r PubSubMessage) error {
	eventType, _ := r.Message.Attributes["type"]
	switch eventType {
	case eventpub.TypeRepositoryAdded:
		var ev eventpub.RepositoryAddedEvent
		if err := json.Unmarshal(r.Message.Data, &ev); err != nil {
			logrus.WithError(err).Error("Cannot unmarshal JSON payload")
			return err
		}

		return c.HandleRepositoryAddedEvent(ctx, ev)
	}
	return nil
}

func (c *Consumer) HandleRepositoryAddedEvent(ctx context.Context, ev eventpub.RepositoryAddedEvent) error {
	m := task.NewManager(ev.Repository)
	res := m.ExecuteRunners()
	if res.Success {
		out, err := json.Marshal(res.Runners)
		if err != nil {
			return err
		}

		var foo exago.Results
		err = json.Unmarshal(out, &foo)
		if err != nil {
			return err
		}

		// forward to score
		// save to firestore
		return nil
	}

	return fmt.Errorf("%#v", res.Errors)
}
