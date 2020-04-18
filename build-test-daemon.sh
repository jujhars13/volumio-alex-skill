#!/bin/bash

(cd sqs-poll-daemon && go fmt main.go && go install)
(cd sqs-poll-daemon && go fmt main.go && GOOS=linux GOARCH=arm GOARM=5 go build -o volumio-sqs-poll-daemon-raspi)

export SQS_ENDPOINT="https://sqs.eu-west-1.amazonaws.com/753637769290/jujhar-test-volumio"
export DOMAIN="volumio"

aws sqs send-message \
    --queue-url="${SQS_ENDPOINT}" \
    --message-body="test-play" 

if [[ "${BUILD}" != "false" ]]; then
    volumio-sqs-poll-daemon
fi
