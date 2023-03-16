package main

import (
	"ddd-golang-microservices/recommendation/internal/recommendation"
	"ddd-golang-microservices/recommendation/internal/transport"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	c := retryablehttp.NewClient()
	c.RetryMax = 10

	partnerAdaptor, err := recommendation.NewPartnershipAdaptor(c.StandardClient(), "http://localhost:3031")
	if err != nil {
		log.Fatal("failed to create partnerAdaptor: ", err)
	}

	svc, err := recommendation.NewService(partnerAdaptor)
	if err != nil {
		log.Fatal("failed to create a service: ", err)
	}

	handler, err := recommendation.NewHandler(*svc)
	if err != nil {
		log.Fatal("failed to create a handler: ", err)
	}

	m := transport.NewMux(*handler)

	port := "4040"
	if err := http.ListenAndServe(":"+port, m); err != nil {
		log.Fatal("server errored:", err)
	}

	log.Println("server listening on port", port)
}
