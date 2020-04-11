package main

import (
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

	var url string = fmt.Sprintf("http://%s/api/v1/commands/?cmd=toggle", domain)
	log.Print(fmt.Sprintf("Toggling Volumio on %s", url))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Printf("%s Status is in the 2xx range", url)
	} else {
		log.Fatalf("Can't toggle Volumio via %s", url)
	}
}
