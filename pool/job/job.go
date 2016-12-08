package job

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	. "github.com/hotolab/exago-svc/config"
)

const (
	fnPrefix = "exago-"
)

var (
	svc       *lambda.Lambda
	ErrNoData = errors.New("Empty dataset")
)

// Response contains the generic JSend response sent by Lambda functions.
type Response struct {
	Success bool              `json:"success"`
	Data    *json.RawMessage  `json:"data"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type context struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch,omitempty"`
}

func New() {
	creds := credentials.NewStaticCredentials(
		Config.AwsAccessKeyID,
		Config.AwsSecretAccessKey,
		"",
	)
	svc = lambda.New(
		session.New(),
		aws.NewConfig().
			WithRegion(Config.AwsRegion).
			WithCredentials(creds),
	)
}

func CallLambdaFn(fn, repo, branch, goversion string) (lrsp Response, err error) {
	payload, _ := json.Marshal(context{
		Repository: repo,
		Branch:     branch,
	})
	params := &lambda.InvokeInput{
		FunctionName: aws.String(fnPrefix + fn),
		Payload:      payload,
		Qualifier:    aws.String(goversion),
	}

	out, err := svc.Invoke(params)
	if err != nil {
		return lrsp, err
	}

	var resp Response
	if err = json.Unmarshal(out.Payload, &resp); err != nil {
		return lrsp, err
	}

	// If the Lambda request failed, return the message as an error
	if !resp.Success {
		for _, msg := range resp.Errors {
			// Return the first error
			return lrsp, errors.New(msg)
		}
	}

	return resp, nil
}
