package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brutella/dnssd"
	"log"
	"os"
	"time"
)

var versionFlag = flag.Bool("version", false, "print version")
var jsonFlag = flag.String("json", "", "[{\"Type\": \"_http._tcp\", \"Name\": \"Fancy Name\", \"Host\": \"myhost\", \"Port\": 1234, \"Text\": {\"foo\": \"bar\"}}]")

var timeFormat = "15:04:05.000"

func main() {
	flag.Parse()

	// Show version and return if asked for
	if *versionFlag != false {
		fmt.Println("unversioned")
		os.Exit(0)
	}

	// Otherwise parse the json flag and unmarshall into a slice of configs
	var servicesDeclaration []dnssd.Config
	err := json.Unmarshal([]byte(*jsonFlag), &servicesDeclaration)
	if err != nil {
		log.Fatal(err)
	}

	// Get context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get our responder
	responder, err := dnssd.NewResponder()
	if err != nil {
		fmt.Println(err)
		return
	}

	var services []dnssd.Service
	for _, serviceConfig := range servicesDeclaration {
		srv, err := dnssd.NewService(serviceConfig)
		if err != nil {
			log.Fatal(err)
		}
		services = append(services, srv)
	}
	/*
		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			select {
			case <-stop:
				cancel()
			}
		}()
	*/

	// Register our services to the repsonder
	for _, srv := range services {
		fmt.Printf("Gonna register %s\n", srv.Name)
		time.Sleep(1 * time.Second)
		handle, err := responder.Add(srv)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s	Got a reply for service %s: Name now registered and active\n", time.Now().Format(timeFormat), handle.Service().ServiceInstanceName())
		}
	}

	err = responder.Respond(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
