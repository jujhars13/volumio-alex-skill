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

const waitTime time.Duration = 20

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
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
	log.Printf("Pausing for %d seconds", waitTime)
	for range time.Tick(waitTime * time.Second) {

		// poll sqs for message
		sqsMsg, sqsErr := pollSqs(sqsEndpoint)
		if sqsErr != nil {
			log.Fatal(sqsErr)
		}
		if sqsMsg != "" {
			log.Print(sqsMsg)
			// call volumio
			volumioErr := callURL(domain)
			if volumioErr != nil {
				log.Fatal(volumioErr)
			}
			log.Printf("Toggled Volumio")
		}

		log.Printf("Pausing for %d seconds", waitTime)
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

// ping sqs queue to see if status has changed, if it has then toggle volumio
func main() {
	svcConfig := &service.Config{
		Name:        "sqs-poll-daemon",
		DisplayName: "SQS Poll daemon",
		Description: "Poll SQS to check for volumio messages.",
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

func callURL(domain string) error {

	var url string = fmt.Sprintf("http://%s/api/v1/commands/?cmd=toggle", domain)
	log.Printf("Toggling Volumio on %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		log.Print("HTTP Status is not in the 2xx range")
		log.Printf("Can't toggle Volumio via %s", url)
		return errors.New("http response non 20x error")
	}

	return nil
}

func pollSqs(sqsURL string) (string, error) {

	log.Printf("Polling %s for a sqs message", sqsURL)

	sess := session.Must(session.NewSession())

	svc := sqs.New(sess)

	receiveMsgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &sqsURL,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return "", err
	}

	log.Printf("Received %d messages.\n", len(receiveMsgResult.Messages))
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
