# Volumio Alexa Skill (via SQS)

Set of lambdas and daemons allow Alexa voice control over [Volumio](https://volumio.org/) audio playback.

Alexa cannot call urls/ip's that are on a private network (`192.168.x.x`) and we don't really want to expose our Volumio box to 'tinternet.  So we'll get Alexa to plop a message onto an SQS queue and get a daemon on Volumio to poll the queue periodically and control Volumio.

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
