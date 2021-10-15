package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"

	"github.com/grandcat/zeroconf"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var versionFlag = flag.Bool("version", false, "print version")
var jsonFlag = flag.String("json", "", "[{\"Type\": \"_http._tcp\", \"Name\": \"Fancy Name\", \"Host\": \"myhost\", \"Port\": 1234, \"Text\": {\"foo\": \"bar\"}}]")

type Config struct {
	// Name of the service
	Name string

	// Host is the name of the host (no trailing dot).
	Host string

	// Type is the service, for example "_hap._tcp".
	Type string

	// Txt records
	Text []string

	// IP addresses of the service.
	IPs []net.IP

	// Port is the port of the service.
	Port int
}

func listMulticastInterfaces() []net.Interface {
	var interfaces []net.Interface
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifi := range ifaces {
		if (ifi.Flags & net.FlagUp) == 0 {
			continue
		}
		if (ifi.Flags & net.FlagMulticast) > 0 {
			interfaces = append(interfaces, ifi)
		}
	}

	return interfaces
}

func addrsForInterface(iface *net.Interface) []string {
	addrs, _ := iface.Addrs()
	var v4, v6, v6local []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				v4 = append(v4, ipnet.IP.String())
			} else {
				switch ip := ipnet.IP.To16(); ip != nil {
				case ip.IsGlobalUnicast():
					v6 = append(v6, ipnet.IP.String())
				case ip.IsLinkLocalUnicast():
					v6local = append(v6local, ipnet.IP.String())
				}
			}
		}
	}
	if len(v6) == 0 {
		v6 = v6local
	}
	var ips []string
	ips = append(v4, v6...)
	return ips
}

func main() {
	flag.Parse()

	// Show version and return if asked for
	if *versionFlag != false {
		fmt.Println("unversioned")
		os.Exit(0)
	}

	var servicesDeclaration []Config
	err := json.Unmarshal([]byte(*jsonFlag), &servicesDeclaration)
	if err != nil {
		log.Fatal(err)
	}

	ifaces := listMulticastInterfaces()

	var ips []string

	for _, iface := range ifaces {
		v := addrsForInterface(&iface)
		ips = append(ips, v...)
	}

	for _, serviceConfig := range servicesDeclaration {
		name := serviceConfig.Name
		service := serviceConfig.Type
		host := serviceConfig.Host
		text := serviceConfig.Text
		port := serviceConfig.Port
		domain := "local."

		server, err := zeroconf.RegisterProxy(name, service, domain, port, host, ips, text, ifaces)

		if err != nil {
			panic(err)
		}
		defer server.Shutdown()

		log.Println("Published service:")
		log.Println("- Name:", name)
		log.Println("- Type:", service)
		log.Println("- Domain:", domain)
		log.Println("- Port:", port)
	}

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
	}

	log.Println("Shutting down.")
}
