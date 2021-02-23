package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brutella/dnssd"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var instanceFlag = flag.String("n", "croquette", "Service Name")
var serviceFlag = flag.String("t", "_http._tcp", "Service type")
var domainFlag = flag.String("d", "local", "Browsing domain")
var modeFlag = flag.String("m", "brute", "Resolution mode (brute or vanilla)")

func main() {
	flag.Parse()

	if *modeFlag == "vanilla" {
		dom := fmt.Sprintf("%s.%s", strings.Trim(*instanceFlag, "."), strings.Trim(*domainFlag, "."))
		ips, err := net.LookupIP(dom)
		if err != nil {
			fmt.Printf("FATAL: %s: %s", dom, err)
			os.Exit(1)
		}
		jj, _ := json.Marshal(ips)
		fmt.Printf(string(jj))
		return
	}

	service := fmt.Sprintf("%s.%s.", strings.Trim(*serviceFlag, "."), strings.Trim(*domainFlag, "."))
	// instance := fmt.Sprintf("%s.%s.%s.", strings.Trim(*instanceFlag, "."), strings.Trim(*serviceFlag, "."), strings.Trim(*domainFlag, "."))

	addFn := func(srv dnssd.Service) {
		j, _ := json.Marshal(srv)

		fmt.Println(string(j))
		//		fmt.Println(srv.ServiceInstanceName())
		if srv.Host == strings.Trim(*instanceFlag, ".") {
			//    if srv.ServiceInstanceName() == instance {
			j, _ := json.Marshal(srv)

			fmt.Println("FOUND ------->")
			fmt.Println(string(j))

			time.Sleep(1)
			os.Exit(0)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := dnssd.LookupType(ctx, service, addFn, func(dnssd.Service) {}); err != nil {
		fmt.Println(err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	select {
	case <-stop:
		cancel()
	}
}
