package lambda

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	. "github.com/exago/svc/config"
	"github.com/exago/svc/taskrunner"
)

const (
	fnPrefix = "exago-"
)

var (
	ErrNoData = errors.New("Empty dataset")

	// Make sure it satisfies the interface.
	_ taskrunner.TaskRunner = (*Runner)(nil)
)

type Runner struct {
	Repository    string
	ShouldCleanup bool
}

// Response contains the generic JSend response sent by Lambda functions.
type Response struct {
	Status   string                 `json:"status"`
	Data     *json.RawMessage       `json:"data"`
	Metadata map[string]interface{} `json:"_metadata"`
}

// context sent to the Lambda functions.
type context struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch,omitempty"`
	Linters    string `json:"linters,omitempty"`
	Cleanup    bool   `json:"cleanup,omitempty"`
}

type cmd struct {
	name      string
	ctxt      context
	data      interface{}
	unMarshal func(l *cmd, j []byte) (data interface{}, err error)
}

// Data returns the response from Lambda.
func (l *cmd) Data() (interface{}, error) {
	res, err := l.call()
	if err != nil {
		return nil, err
	}

	if l.data, err = l.unMarshal(l, *res.Data); err != nil {
		return nil, err
	}

	return l.data, nil
}

func (l *cmd) call() (lrsp Response, err error) {
	creds := credentials.NewStaticCredentials(
		Config.AwsAccessKeyID,
		Config.AwsSecretAccessKey,
		"",
	)
	svc := lambda.New(
		session.New(),
		aws.NewConfig().
			WithRegion(Config.AwsRegion).
			WithCredentials(creds),
	)

	payload, _ := json.Marshal(l.ctxt)
	params := &lambda.InvokeInput{
		FunctionName: aws.String(fnPrefix + l.name),
		Payload:      payload,
	}

	out, err := svc.Invoke(params)
	if err != nil {
		return lrsp, err
	}

	var resp struct {
		Status   string                 `json:"status"`
		Data     *json.RawMessage       `json:"data"`
		Metadata map[string]interface{} `json:"_metadata"`
	}
	if err = json.Unmarshal(out.Payload, &resp); err != nil {
		return lrsp, err
	}

	// Data is always expected from Lambda
	if resp.Data == nil {
		return lrsp, ErrNoData
	}

	// If the Lambda request failed, return the message as an error
	if resp.Status == "fail" {
		var msg struct {
			// Message is the only expected field in Data
			Message string `json:"message"`
		}
		if err = json.Unmarshal(*resp.Data, &msg); err != nil {
			return lrsp, err
		}
		return lrsp, errors.New(msg.Message)
	}

	return resp, err
}
