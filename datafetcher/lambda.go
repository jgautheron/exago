package datafetcher

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/exago/svc/config"
)

const (
	lambdaRegion = "eu-west-1"
)

type lambdaContext struct {
	Registry   string `json:"registry"`
	Username   string `json:"username"`
	Repository string `json:"repository"`
	Branch     string `json:"branch,omitempty"`
	Linters    string `json:"linters,omitempty"`
	Indexes    string `json:"indexes,omitempty"`
}

type LambdaResponse struct {
	Data     *json.RawMessage       `json:"data"`
	Metadata map[string]interface{} `json:"_metadata"`
}

func callLambdaFn(fn string, ctxt lambdaContext) (*json.RawMessage, error) {
	creds := credentials.NewStaticCredentials(
		config.Get("AwsAccessKeyID"),
		config.Get("AwsSecretAccessKey"),
		"",
	)
	svc := lambda.New(
		session.New(),
		aws.NewConfig().
			WithRegion(lambdaRegion).
			WithCredentials(creds),
	)

	payload, _ := json.Marshal(ctxt)

	params := &lambda.InvokeInput{
		FunctionName: aws.String("exago-" + fn),
		Payload:      payload,
	}

	out, err := svc.Invoke(params)
	if err != nil {
		return nil, err
	}

	var resp LambdaResponse
	err = json.Unmarshal(out.Payload, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, err
}
