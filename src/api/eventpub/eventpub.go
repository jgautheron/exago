package eventpub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/jgautheron/exago/src/api/config"
	"github.com/pkg/errors"
)

var _ EventPublisher = (*EventPub)(nil)

type EventPub struct {
	client *pubsub.Client
	config *config.GoogleCloudConfig
}

func NewFromConfig(ctx context.Context, gcCfg *config.GoogleCloudConfig) (*EventPub, error) {
	client, err := pubsub.NewClient(ctx, gcCfg.GoogleProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize pubsub client")
	}
	return &EventPub{client, gcCfg}, nil
}

func (evp *EventPub) sendEvent(typ string, event interface{}, topic string) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.Wrapf(err, "Could not compose the payload of the `%s` event for queue", typ)
	}
	err = evp.publish(typ, payload, topic)
	if err != nil {
		return errors.Wrapf(err, "Could not publish an `%s` event with payload ```%s```", typ, payload)
	}
	return nil
}

func (evp *EventPub) publish(typ string, msg []byte, topic string) error {
	ctx := context.Background()
	t := evp.client.Topic(topic)

	// https://github.com/googleapis/google-cloud-go/issues/1552
	t.PublishSettings.NumGoroutines = 1

	result := t.Publish(ctx, &pubsub.Message{
		Attributes: map[string]string{
			"type": typ,
		},
		Data: msg,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err := result.Get(ctx)
	if err != nil {
		return err
	}
	return nil
}
