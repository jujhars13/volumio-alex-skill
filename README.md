# Volumio Alexa Skill (via SQS)

![volumio logo](logo.png)

Set of lambdas and daemons allow Alexa voice control over [Volumio](https://volumio.org/) audio playback.

Alexa cannot call urls/ip's that are on a private network (`192.168.x.x`) and we don't really want to expose our Volumio box to whole 'tinternet or try some brittle port-forwarding/proxying.  

So we'll get Alexa to plop a message onto an SQS queue and get a daemon on Volumio to poll the queue periodically and control Volumio.

## Architecture

Daemon [polls SQS](poll/) for a message in an SQS queue and then sends an appropriate http request to [Volumio's API](https://volumio.github.io/docs/API/REST_API.html)

a bit like this:

```bash
curl http://volumio/api/v1/commands/\?cmd=toggle
```

SQS queue has message put there via Alexa skill.

### volumio-sqs-poll-daemon

```bash
# env vars
export DOMAIN="volumio" # can be localhost if run on Raspberry pi itself
export SQS_ENDPOINT="sqs://my-aws-sqs-endpoint"
volumio-sqs-poll
```

### Alexa skill

TODO

## TODO

### v1

- [x] write poller
- [x] create manual SQS queue and command line bash addMessage
- [x] test poller
- [x] test poller on RasPi
- [x] deploy daemon to RasPi
- [ ] write Alexa skill
- [ ] create sqs queue
- [ ] write CF for AWS deployment

### v2

- [ ] respond to content of SQS message (play,pause,next track etc)

## Licence

[MIT](LICENCE.txt)
