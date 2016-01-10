package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
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

type lambdaResponse struct {
	Data     *json.RawMessage       `json:"data"`
	Metadata map[string]interface{} `json:"_metadata"`
}

func callLambdaFn(fn string, ctxt lambdaContext) (resp lambdaResponse, err error) {
	creds := credentials.NewStaticCredentials(cfg.awsAccessKeyID, cfg.awsSecretAccessKey, "")
	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lambdaRegion).WithCredentials(creds))

	payload, _ := json.Marshal(ctxt)

	params := &lambda.InvokeInput{
		FunctionName: aws.String("exago-" + fn),
		Payload:      payload,
	}

	out, err := svc.Invoke(params)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(out.Payload, &resp)
	if err != nil {
		return resp, err
	}

	return resp, err
}
