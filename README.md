# Volumio Alexa Skill (via SQS)

Set of lambdas and daemons allow voice control over [Volumio](https://volumio.org/) audio playback.

## Architecture

Daemon [polls SQS](poll/) for a message in an SQS queue and then sends an appropriate http request to [Volumio's API](https://volumio.github.io/docs/API/REST_API.html)

SQS queue has message put there via Alexa skill.

### volumio-sqs-poll

TODO

```bash
# env vars
export DOMAIN="volumio" # can be localhost if run on Raspberry pi itself

volumio-sqs-poll

```

### Alexa skill

TODO

## Licence

[MIT](LICENCE.txt)