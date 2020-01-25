package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/jgautheron/exago/internal/config"
	"github.com/pkg/errors"
)

type Firestore struct {
	client *firestore.Client
}

func NewFromConfig(ctx context.Context, gcCfg *config.GoogleCloudConfig) (*Firestore, error) {
	client, err := firestore.NewClient(ctx, gcCfg.GoogleProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "Could not initialize firestore client")
	}
	return &Firestore{client}, nil
}
