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
)

var instanceFlag = flag.String("n", "croquette", "Service Name")
var serviceFlag = flag.String("t", "_http._tcp", "Service type")
var domainFlag = flag.String("d", "local.", "Browsing domain")
var modeFlag = flag.String("m", "brute", "Resolution mode (brute or vanilla)")

var (
	version = flag.Bool("version", false, "print version")
)

// var debug = flag.Bool("debug", true, "Dump")

func vanilla() {
	dom := fmt.Sprintf("%s.%s", strings.Trim(*instanceFlag, "."), strings.Trim(*domainFlag, "."))
	ips, err := net.LookupIP(dom)
	if err != nil {
		fmt.Printf("FATAL: %s: %s", dom, err)
		os.Exit(1)
	}
	jj, _ := json.Marshal(ips)
	fmt.Println("----- Got vanilla entry as follow -----")
	fmt.Printf(string(jj))
	os.Exit(0)
}

var timeFormat = "15:04:05.000"

func main() {
	flag.Parse()

	if *version != false {
		fmt.Println("unversioned")
		os.Exit(0)
	}

	if len(*instanceFlag) == 0 || len(*serviceFlag) == 0 || len(*domainFlag) == 0 {
		flag.Usage()
		return
	}

	if *modeFlag == "vanilla" {
		vanilla()
	}

	service := fmt.Sprintf("%s.%s.", strings.Trim(*serviceFlag, "."), strings.Trim(*domainFlag, "."))
	//instance := fmt.Sprintf("%s.%s.%s.", strings.Trim(*instanceFlag, "."), strings.Trim(*serviceFlag, "."), strings.Trim(*domainFlag, "."))

	addFn := func(e dnssd.BrowseEntry) {
		m, _ := json.Marshal(e)
		fmt.Println(string(m))
		if e.Host == strings.Trim(*instanceFlag, ".") {
			//text := ""
			//for key, value := range e.Text {
			//	text += fmt.Sprintf("%s=%s", key, value)
			//}
			os.Exit(0)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := dnssd.LookupType(ctx, service, addFn, func(dnssd.BrowseEntry) {}); err != nil {
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
