package config

import (
	"log"
	"os"
)

var data map[string]*configItem

type configItem struct {
	value, defaultValue string
	required            bool
}

func SetUp() {
	data = map[string]*configItem{
		"GithubAccessToken": {
			os.Getenv("GITHUB_ACCESS_TOKEN"), "", true,
		},
		"AwsRegion": {
			os.Getenv("AWS_REGION"), "eu-west-1", false,
		},
		"AwsAccessKeyID": {
			os.Getenv("AWS_ACCESS_KEY_ID"), "", true,
		},
		"AwsSecretAccessKey": {
			os.Getenv("AWS_SECRET_ACCESS_KEY"), "", true,
		},
		"RunnerImageName": {
			os.Getenv("RUNNER_IMAGE_NAME"), "jgautheron/exago-runner", false,
		},
		"HttpPort": {
			os.Getenv("HTTP_PORT"), "8080", false,
		},
		"DatabasePath": {
			os.Getenv("DATABASE_PATH"), "/data/exago.db", false,
		},
		"AllowOrigin": {
			os.Getenv("ALLOW_ORIGIN"), "", true,
		},
		"LogLevel": {
			os.Getenv("LOG_LEVEL"), "info", false,
		},
	}

	for k, m := range data {
		if m.required && len(m.value) == 0 {
			log.Fatalf("Missing value for %s", k)
		}
		if len(m.value) == 0 && len(m.defaultValue) > 0 {
			m.value = m.defaultValue
		}
	}
}

func Get(key string) string {
	if _, exists := data[key]; !exists {
		return ""
	}
	return data[key].value
}
