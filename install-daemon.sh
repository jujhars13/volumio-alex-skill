#!/bin/bash
# install daemon to raspi

# build raspi binary
(cd sqs-poll-daemon && \
go fmt main.go && \
GOOS=linux GOARCH=arm GOARM=5 go build -o volumio-sqs-poll-daemon-raspi)

# copy required files over to volumio instance
# replace aws secrets on copy
< deploy/sqs-poll-daemon-systemd.service \
    sed s/__AWS_ACCESS_KEY_ID__/$(< "${SECRETS_DIR}/jujhar/aws-jujhar-volumio-sqs-access-key")/ | \
    sed s/__AWS_SECRET_ACCESS_KEY__/$(< "${SECRETS_DIR}/jujhar/aws-jujhar-volumio-sqs-access-secret")/ | \
    ssh volumio "cat >/home/jujhar/sqs-poll-daemon-systemd.service"
scp sqs-poll-daemon/volumio-sqs-poll-daemon-raspi volumio:/home/jujhar/volumio-sqs-poll-daemon-raspi

# execute commands on raspi
ssh volumio << "EOF"
    date
    sudo systemctl stop sqs-poll-daemon-systemd
    sudo cp /home/jujhar/sqs-poll-daemon-systemd.service /etc/systemd/system/sqs-poll-daemon-systemd.service
    sudo chmod 644 /etc/systemd/system/sqs-poll-daemon-systemd.service
    sudo cp /home/jujhar/volumio-sqs-poll-daemon-raspi /opt/volumio-sqs-poll-daemon-raspi
    sudo chmod 777 /opt/volumio-sqs-poll-daemon-raspi
    sudo systemctl enable sqs-poll-daemon-systemd
    sudo systemctl start sqs-poll-daemon-systemd
EOF
