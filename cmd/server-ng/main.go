package main

import (
	// "context"
	"encoding/json"
	"flag"
	"fmt"
	// "github.com/brutella/dnssd"
	"net"

	// "github.com/brutella/dnssd"
	"github.com/grandcat/zeroconf"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var versionFlag = flag.Bool("version", false, "print version")
var jsonFlag = flag.String("json", "", "[{\"Type\": \"_http._tcp\", \"Name\": \"Fancy Name\", \"Host\": \"myhost\", \"Port\": 1234, \"Text\": {\"foo\": \"bar\"}}]")

// var timeFormat = "15:04:05.000"

// var waitTime = flag.Int("wait", 10, "Duration in [s] to publish service for.")

/*
func mainOld() {
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
*/

type Config struct {
	// Name of the service
	Name string

	// Type is the service type, for example "_hap._tcp".
	Type string

	// Domain is the name of the domain, for example "local".
	// If empty, "local" is used.
	Domain string

	// Host is the name of the host (no trailing dot).
	// If empty the local host name is used.
	Host string

	// Txt records
	Text []string
	// map[string]string

	// IP addresses of the service.
	// This field is deprecated and should not be used.
	IPs []net.IP

	// Port is the port of the service.
	Port int

	// Interfaces at which the service should be registered
	Ifaces []string
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

/*func addrsForInterface(iface *net.Interface) ([]net.IP, []net.IP) {
	var v4, v6, v6local []net.IP
	addrs, _ := iface.Addrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				v4 = append(v4, ipnet.IP)
			} else {
				switch ip := ipnet.IP.To16(); ip != nil {
				case ip.IsGlobalUnicast():
					v6 = append(v6, ipnet.IP)
				case ip.IsLinkLocalUnicast():
					v6local = append(v6local, ipnet.IP)
				}
			}
		}
	}
	if len(v6) == 0 {
		v6 = v6local
	}
	return v4, v6
}
*/

func addrsForInterface(iface *net.Interface) []string {
	addrs, _ := iface.Addrs()
	var v4, v6, v6local []string
	// var v4, v6, v6local []net.IP
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
		// ips := servicesDeclaration[0].IPs
		port := serviceConfig.Port
		// ifaces := servicesDeclaration[0].Ifaces
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
	// Timeout timer.
	/*
		var tc <-chan time.Time
		if *waitTime > 0 {
			tc = time.After(time.Second * time.Duration(*waitTime))
		}
	*/
	/*
		select {
		case <-sig:
			// Exit by user
		case <-tc:
			// Exit by timeout
		}*/

	/* go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		select {
		// case <-tc:
		case <-stop:
			// cancel()
		}
	}()*/

	select {
	case <-sig:
	}

	log.Println("Shutting down.")
}
