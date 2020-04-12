package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ping sqs queue to see if status has changed, if it has then toggle volumio
func main() {
	domain, exists := os.LookupEnv("DOMAIN")
	if !exists {
		domain = "volumio"
	}

	// sqs endpoint
	sqsEndpoint, exists := os.LookupEnv("SQS_ENDPOINT")
	if !exists {
		log.Fatal("You must export a SQS_ENDPOINT")
	}

	// poll sqs for message
	sqsMsg, sqsErr := pollSqs(sqsEndpoint)
	if sqsErr != nil {
		log.Fatal(sqsErr)
	}
	if sqsMsg == "" {
		log.Print("No SQS messages, exiting...")
		os.Exit(0)
	}
	log.Print(sqsMsg)

	// call volumio
	volumioErr := callURL(domain)
	if volumioErr != nil {
		log.Fatal(volumioErr)
	}
	log.Printf("Toggled Volumio")
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

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &sqsURL,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return "", err
	}

	log.Printf("Received %d messages.\n", len(result.Messages))
	if len(result.Messages) == 0 {
		return "", nil
	}

	return *result.Messages[0].Body, nil
}
