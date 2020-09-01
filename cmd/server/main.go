package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brutella/dnssd"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var instanceFlag = flag.String("name", "Service", "Service name")
var serviceFlag = flag.String("type", "_asdf._tcp", "Service type")
var domainFlag = flag.String("domain", "local", "domain")
var portFlag = flag.Int("port", 12345, "Port")

var timeFormat = "15:04:05.000"

var hostFlag = flag.String("host", "", "Host")

func main() {
	flag.Parse()
	if len(*instanceFlag) == 0 || len(*serviceFlag) == 0 || len(*domainFlag) == 0 {
		flag.Usage()
		return
	}

	instance := fmt.Sprintf("%s.%s.%s.", strings.Trim(*instanceFlag, "."), strings.Trim(*serviceFlag, "."), strings.Trim(*domainFlag, "."))

	fmt.Printf("Registering Service %s port %d\n", instance, *portFlag)
	fmt.Printf("DATE: –––%s–––\n", time.Now().Format("Mon Jan 2 2006"))
	fmt.Printf("%s	...STARTING...\n", time.Now().Format(timeFormat))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if resp, err := dnssd.NewResponder(); err != nil {
		fmt.Println(err)
	} else {
		cfg := dnssd.Config{
			Name:   *instanceFlag,
			Type:   *serviceFlag,
			Domain: *domainFlag,
			Port:   *portFlag,
		}

		if *hostFlag != "" {
			cfg.Host = *hostFlag
		}

		srv, err := dnssd.NewService(cfg)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			select {
			case <-stop:
				cancel()
			}
		}()

		go func() {
			time.Sleep(1 * time.Second)
			handle, err := resp.Add(srv)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%s	Got a reply for service %s: Name now registered and active\n", time.Now().Format(timeFormat), handle.Service().ServiceInstanceName())
			}
		}()
		err = resp.Respond(ctx)

		if err != nil {
			fmt.Println(err)
		}
	}
}
