module github.com/dubo-dubon-duponey/goello

replace github.com/brutella/dnssd v1.2.0 => github.com/dubo-dubon-duponey/dnssd v1.2.1-0.20210609210013-5ca5fb2acdd2

require (
	github.com/brutella/dnssd v1.2.0
	// github.com/grandcat/zeroconf v1.0.0
	github.com/grandcat/zeroconf v1.0.1-0.20210929195321-a393c0e41e54
	github.com/miekg/dns v1.1.42
	github.com/pion/mdns v0.0.5
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
)

go 1.16
