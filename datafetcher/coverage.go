package datafetcher

import (
	"bytes"
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/exago/svc/config"
	"github.com/fsouza/go-dockerclient"
	"github.com/twinj/uuid"
)

var (
	errEmptyOutput = errors.New("The container output is empty")
)

func GetCoverage(log *log.Entry, repository string) (interface{}, error) {
	resp := make(chan interface{})
	go runContainer(log, repository, resp)

	out := <-resp
	switch out.(type) {
	case error:
		return nil, out.(error)
	case *json.RawMessage:
		return out, nil
	}

	return nil, nil
}

func runContainer(log *log.Entry, repository string, resp chan interface{}) {
	client, _ := docker.NewClientFromEnv()
	cfg := docker.Config{
		Image: config.Get("RunnerImageName"),
		Cmd:   []string{repository},
	}

	// Prevent collisions
	cn := "exago-" + uuid.NewV4().String()
	opts := docker.CreateContainerOptions{Name: cn, Config: &cfg}
	container, err := client.CreateContainer(opts)
	if err != nil {
		resp <- err
		return
	}
	log.WithField("container", cn).Info("Container created")

	defer func() {
		err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
		if err != nil {
			log.WithField("container", cn).Error(err)
			return
		}
	}()

	err = client.StartContainer(container.ID, &docker.HostConfig{})
	if err != nil {
		resp <- err
		return
	}

	var buf bytes.Buffer
	err = client.Logs(docker.LogsOptions{
		Container:    container.ID,
		OutputStream: &buf,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
		Timestamps:   false,
	})

	if buf.Len() == 0 {
		resp <- errEmptyOutput
		return
	}

	var msg *json.RawMessage
	if err = json.Unmarshal(buf.Bytes(), &msg); err != nil {
		resp <- err
		return
	}

	resp <- msg
}
