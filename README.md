# Exago [![Circle CI](https://circleci.com/gh/exago/svc.svg?style=svg)](https://circleci.com/gh/exago/svc) [![](https://badge.imagelayers.io/jgautheron/exago-service:latest.svg)](https://imagelayers.io/?images=jgautheron/exago-service:latest 'Get your own badge on imagelayers.io')

Exago is a code quality tool that inspects your Go repository and reports on what could be improved. The dashboard displays metrics that we consider as your application pillars, you can dive deeper and browse directly the recommandations in the code.

This is the API backend consumed by the [app](https://github.com/exago/app).

## Getting started

### Development

This tool assumes you are working in a standard Go 1.5 workspace, as described in http://golang.org/doc/code.html.

#### Configuration

The configuration is managed exclusively through environment variables.

Variable               | Description | Mandatory
---------------- | ------ | ------------
GITHUB_ACCESS_TOKEN       | Necessary to consume GitHub's API | Yes
AWS_ACCESS_KEY_ID        | Required for AWS Lambda | Yes
AWS_SECRET_ACCESS_KEY     | Required for AWS Lambda | Yes
RUNNER_IMAGE_NAME      | Required for running tests | Yes
HTTP_PORT      | HTTP port to bind | Yes
REDIS_HOST      | Redis host | No
ALLOW_ORIGIN   | Origin allowed for API calls (CORS) | Yes
LOG_LEVEL   | Log level (debug, info, warn, error, fatal) | Yes

#### Redis

Redis is used as cache backend.  
The service will degrade gracefully if there's no Redis available for the given `REDIS_HOST`.

## Contributing

See the [dedicated page](CONTRIBUTING.md).
