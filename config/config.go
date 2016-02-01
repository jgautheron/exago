package config

import (
	"errors"
	"os"
)

var data map[string]string

func SetUp() error {
	data = map[string]string{
		// Should be overridable later by a logged in user
		"GithubAccessToken":  os.Getenv("GITHUB_ACCESS_TOKEN"),
		"AwsAccessKeyID":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"AwsSecretAccessKey": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"RunnerImageName":    os.Getenv("RUNNER_IMAGE_NAME"),
		"HttpPort":           os.Getenv("HTTP_PORT"),
		"RedisHost":          os.Getenv("REDIS_HOST"),
		"PapertrailURL":      os.Getenv("PAPERTRAIL_URL"),
		"PapertrailPort":     os.Getenv("PAPERTRAIL_PORT"),
		"PapertrailHost":     os.Getenv("PAPERTRAIL_HOST"),
		"PapertrailApp":      os.Getenv("PAPERTRAIL_APP"),
		"Env":                os.Getenv("ENV"),
		"AllowOrigin":        os.Getenv("ALLOW_ORIGIN"),
		"LogLevel":           os.Getenv("LOG_LEVEL"),
	}

	// Basic validation
	if data["GithubAccessToken"] == "" ||
		data["AwsAccessKeyID"] == "" ||
		data["AwsSecretAccessKey"] == "" ||
		data["HttpPort"] == "" ||
		data["RedisHost"] == "" {
		return errors.New("Missing environment variable(s)")
	}

	if data["Env"] == "" {
		// Possible values: dev, prod
		data["Env"] = "dev"
	}

	return nil
}

func Get(key string) string {
	if _, exists := data[key]; !exists {
		return ""
	}
	return data[key]
}
