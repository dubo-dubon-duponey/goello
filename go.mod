module github.com/dubo-dubon-duponey/goello

replace github.com/brutella/dnssd v1.2.0 => github.com/dubo-dubon-duponey/dnssd v1.2.1-0.20210609210013-5ca5fb2acdd2

go 1.17

require (
	github.com/brutella/dnssd v1.2.0
	github.com/grandcat/zeroconf v1.0.1-0.20210929195321-a393c0e41e54
	github.com/miekg/dns v1.1.43
	github.com/pion/mdns v0.0.5
	golang.org/x/net v0.0.0-20211013171255-e13a2654a71e
)

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/pion/logging v0.2.2 // indirect
	golang.org/x/sys v0.0.0-20210426080607-c94f62235c83 // indirect
)
