# Exago [![CircleCI](https://circleci.com/gh/jgautheron/exago.svg?style=svg)](https://circleci.com/gh/jgautheron/exago) [![](https://badge.imagelayers.io/jgautheron/exago:latest.svg)](https://imagelayers.io/?images=jgautheron/exago:latest 'Get your own badge on imagelayers.io')

Exago is a code quality tool that inspects your Go repository and reports on what could be improved. The dashboard displays metrics that we consider as your application pillars, you can dive deeper and browse directly the recommandations in the code.

This is the API backend consumed by the [app](https://github.com/jgautheron/exago-app).

## How it works

Exago "outsources" the entire code processing to two dedicated AWS Lambda functions that take care of pulling and processing the code concurrently.  

Once both functions are done processing the code, this service retrieves the KPIs and outputs them as JSON formatted output, which is then displayed by the [frontend](https://github.com/jgautheron/exago-app). Every successful processing is cached in LevelDB.

For those familiar with Lambda, there is a default limit of 100 concurrent running functions, which means that we can process about 50 projects simultaneously. Reaching that limit won't cause any dysfunction as Exago is relying on a routine pool to execute orderly repository checks.

## Rank calculation

All metrics are factored with the amount of lines of code, so that bigger projects can also get into high ranks. Actually the bigger the project the more we try to be permissive.

The amount of tests is not factored in the calculation of the rank, the dot only shows two colors: green or red, depending if the tests ran successfully or not.

We spent months tuning the [algorithms](https://docs.google.com/spreadsheets/d/150xwGQVrY-3qH8-VNDqCDzcQRC0wsIG7-uMps-fXQB0/edit?usp=sharing) to make sure that Exago is fair with every project size. Evaluating a project's quality is tricky and metrics are not everything, but there are a few indicators that can give us hints.

1. If the README is missing or the code is not gofmt'd it's usually a bad sign
2. Having many third parties is a liability but it's important to have good ones: preferably higher rank & maintained (we show all third parties, not only the direct ones - also if you click on the detailed view you will see their rank)
3. No test is a bad sign but high code coverage is not necessarily a must, what's important is testing what really matters, not reaching 99%.
4. Each linter is not equal and has a different weight

Exago tries its best to show you relevant KPIs but ultimately it's up to you to make your own opinion.

## Known limits

- The repository analysis will fail if processing the code exceeds AWS Lambda's 5 mins time limit
- Every project that relies on `CGO` will fail since it's disabled
- Not `go-get`table projects will fail

## Building blocks

- [exago-app](https://github.com/jgautheron/exago-app) - The exago client
- [exago-runner](https://github.com/jgautheron/exago-runner) - Runs all quality checks concurrently
- [exago-lambda-project-runner](https://github.com/jgautheron/exago-lambda-project-runner) - Function that runs `exago-runner` in Amazon Lambda
- [cov](https://github.com/chreble/cov) - The coverage tool generator which includes every possible detail
- [go](https://github.com/chreble/go) - Go fork that speeds up `go get` by implementing the option to do shallow clones

## Getting started

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
POOL_SIZE   | Processing pool size | Yes

## Contributing

See the [dedicated page](CONTRIBUTING.md).

## Contributors

- Karol Górecki [@karolgorecki](https://twitter.com/karolgorecki)
- Christophe Eble [@christopheeble](https://twitter.com/christopheeble)
