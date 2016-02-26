package datafetcher

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/exago/svc/config"
	"github.com/exago/svc/leveldb"
)

var (
	errNoData = errors.New("Empty dataset")
)

type lambdaContext struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch,omitempty"`
	Linter     string `json:"linters,omitempty"`
}

func (lc lambdaContext) String() string {
	return fmt.Sprintf("%s-%s-%s", lc.Repository, lc.Branch, lc.Linter)
}

type LambdaResponse struct {
	Status   string                 `json:"status"`
	Data     *json.RawMessage       `json:"data"`
	Metadata map[string]interface{} `json:"_metadata"`
}

func callLambdaFn(fn string, ctxt lambdaContext) (lrsp LambdaResponse, err error) {
	creds := credentials.NewStaticCredentials(
		config.Values.AwsAccessKeyID,
		config.Values.AwsSecretAccessKey,
		"",
	)
	svc := lambda.New(
		session.New(),
		aws.NewConfig().
			WithRegion(config.Values.AwsRegion).
			WithCredentials(creds),
	)

	payload, _ := json.Marshal(ctxt)
	params := &lambda.InvokeInput{
		FunctionName: aws.String("exago-" + fn),
		Payload:      payload,
	}

	out, err := svc.Invoke(params)
	if err != nil {
		return lrsp, err
	}

	var resp LambdaResponse
	if err = json.Unmarshal(out.Payload, &resp); err != nil {
		return lrsp, err
	}

	// Data is always expected from Lambda
	if resp.Data == nil {
		return lrsp, errNoData
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

type lambdaCmd struct {
	name      string
	ctxt      lambdaContext
	score     int
	data      interface{}
	unMarshal func(l *lambdaCmd, j []byte) (data interface{}, err error)
}

// Data retrieves the data for the given command transparently
// whether the data is already cached or not.
func (l *lambdaCmd) Data() (interface{}, error) {
	var data []byte

	// Check if the data is cached
	isCached, cacheKey := false, l.cacheKey()
	cres, err := leveldb.FindForRepositoryCmd(cacheKey)
	if len(cres) > 0 && err == nil {
		isCached = true
		data = cres
	}

	if !isCached {
		// Fetch the data with Lambda
		res, err := callLambdaFn(l.name, l.ctxt)
		if err != nil {
			return nil, err
		}
		data = *res.Data
	}

	if l.data, err = l.unMarshal(l, data); err != nil {
		return nil, err
	}

	if !isCached {
		// Cache the data
		if err := leveldb.Save(cacheKey, data); err != nil {
			return nil, err
		}
	}

	return l.data, err
}

// cacheKey generates a key that will be used to identify the entry in database.
func (l *lambdaCmd) cacheKey() []byte {
	ck := fmt.Sprintf("%s-%s", l.ctxt.String(), l.name)
	return []byte(ck)
}
