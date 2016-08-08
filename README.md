# Exago [![Circle CI](https://circleci.com/gh/hotolab/exago-svc.svg?style=svg)](https://circleci.com/gh/hotolab/exago-svc) [![](https://badge.imagelayers.io/jgautheron/exago-service:latest.svg)](https://imagelayers.io/?images=jgautheron/exago-service:latest 'Get your own badge on imagelayers.io')

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
AWS_REGION     | Required for AWS Lambda | Yes
HTTP_PORT      | HTTP port to bind | Yes
DATABASE_PATH      | Path to the database | No
ALLOW_ORIGIN   | Origin allowed for API calls (CORS) | Yes
LOG_LEVEL   | Log level (debug, info, warn, error, fatal) | Yes

#### Third Parties

All third parties are in the `vendor` folder.  
:warning: Don't forget to set `GO15VENDOREXPERIMENT` to `1` if you're not running at least Go 1.6.

## Contributing

See the [dedicated page](CONTRIBUTING.md).
