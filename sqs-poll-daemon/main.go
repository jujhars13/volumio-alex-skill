package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kardianos/service"
)

var logger service.Logger

// how many seconds between polls
const waitTime time.Duration = 62

// Program structures.
// Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	logger.Infof("I'm running %v.", service.Platform())
	// Do work here
	domain, exists := os.LookupEnv("DOMAIN")
	if !exists {
		domain = "volumio"
	}

	// sqs endpoint
	sqsEndpoint, exists := os.LookupEnv("SQS_ENDPOINT")
	if !exists {
		log.Fatal("You must export a SQS_ENDPOINT")
	}

	// repeat poll https://gist.github.com/ryanfitz/4191392
	logger.Infof("Pausing for %d seconds", waitTime)
	for range time.Tick(waitTime * time.Second) {

		// poll sqs for message
		sqsMsg, sqsErr := pollSqs(sqsEndpoint)
		if sqsErr != nil {
			log.Fatal(sqsErr)
		}
		if sqsMsg != "" {
			logger.Infof(sqsMsg)
			// call volumio
			volumioErr := callURL(domain)
			if volumioErr != nil {
				log.Fatal(volumioErr)
			}
			logger.Infof("Toggled Volumio")
		}

		logger.Infof("Pausing for %d seconds", waitTime)
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

// ping sqs queue to see if status has changed, if it has then toggle volumio
func main() {
	svcConfig := &service.Config{
		Name:        "volumio-sqs-poll-daemon",
		DisplayName: "Volumio SQS Poll Daemon",
		Description: "Poll SQS to check for volumio messages in SQS.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

// Call volumio's REST api
func callURL(domain string) error {

	var url string = fmt.Sprintf("http://%s/api/v1/commands/?cmd=toggle", domain)
	logger.Infof("Toggling Volumio on %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		logger.Info("HTTP Status is not in the 2xx range")
		logger.Infof("Can't toggle Volumio via %s", url)
		return errors.New("http response non 20x error")
	}

	return nil
}

// Check to see if we have a message
func pollSqs(sqsURL string) (string, error) {

	logger.Infof("Polling %s for a sqs message", sqsURL)

	awsSession := session.Must(session.NewSession())

	svc := sqs.New(awsSession)

	receiveMsgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &sqsURL,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return "", err
	}

	logger.Infof("Received %d messages.\n", len(receiveMsgResult.Messages))
	if len(receiveMsgResult.Messages) == 0 {
		return "", nil
	}

	// remove message from queue
	_, delErr := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &sqsURL,
		ReceiptHandle: receiveMsgResult.Messages[0].ReceiptHandle,
	})
	if delErr != nil {
		return "", delErr
	}

	return *receiveMsgResult.Messages[0].Body, nil
}
