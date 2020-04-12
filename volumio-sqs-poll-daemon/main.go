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
	domain, exists := os.LookupEnv("domain")
	if !exists {
		domain = "volumio"
	}
	success, err := callURL(domain)
}

func callURL(domain) (bool, err) {

	var url string = fmt.Sprintf("http://%s/api/v1/commands/?cmd=toggle", domain)
	log.Print(fmt.Sprintf("Toggling Volumio on %s", url))
	resp, err := http.Get(url)
	if err != nil {
		log.Printf(err)
		return 0, err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		log.Printf("%s Status is not in the 2xx range", url)
		log.Printf("Can't toggle Volumio via %s", url)
		return 0, errors.New("http response error")
	}

	return true, nil
}
