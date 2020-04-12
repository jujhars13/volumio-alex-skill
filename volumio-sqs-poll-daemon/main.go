package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
		log.Print("No SQS message, exiting...")
		os.Exit(0)
	}

	// call volumio
	volumioErr := callURL(domain)
	if volumioErr != nil {
		log.Fatal(volumioErr)
	}
	log.Printf("Toggled Volumio")
}

func callURL(domain string) error {

	var url string = fmt.Sprintf("http://%s/api/v1/commands/?cmd=toggle", domain)
	log.Print(fmt.Sprintf("Toggling Volumio on %s", url))
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

func pollSqs(url string) (string, error) {

	log.Print(fmt.Sprintf("Polling %s for a sqs message", url))

	return "", nil
}
