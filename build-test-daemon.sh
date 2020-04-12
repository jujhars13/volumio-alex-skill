#!/bin/bash

(cd volumio-sqs-poll-daemon && go fmt main.go && go install)

export SQS_ENDPOINT="https://sqs.eu-west-1.amazonaws.com/753637769290/jujhar-test-volumio"
export DOMAIN="volumio"

aws sqs send-message \
    --queue-url="${SQS_ENDPOINT}" \
    --message-body="test-play" 
    
volumio-sqs-poll-daemon