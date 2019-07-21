package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/latikawadhwani/simple-service-broker/pkg/broker"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

func statusAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	logger := lager.NewLogger("simple-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	servicebroker := broker.NewBrokerImpl(logger)

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: os.Getenv("AUTH_USER"),
		Password: os.Getenv("AUTH_PASSWORD"),
	}
	brokerAPI := brokerapi.New(servicebroker, logger, brokerCredentials)

	http.HandleFunc("/health", statusAPI)
	http.Handle("/", brokerAPI)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("\n\nStarting Stupid Service Broker on 0.0.0.0:" + port)
	logger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
