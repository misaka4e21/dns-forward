package main

import (
	"net"
	"os"
	"strconv"
)

func main() {
	ip := net.ParseIP(os.Getenv("UPSTREAM_DNS_IP"))
	port, _ := strconv.ParseInt(os.Getenv("UPSTREAM_DNS_PORT"), 10, 32)
	local_port, _ := strconv.ParseInt(os.Getenv("LOCAL_DNS_PORT"), 10, 32)
	NewDNSService(net.UDPAddr{IP: ip, Port: int(port)}).Listen(int(local_port))
}
