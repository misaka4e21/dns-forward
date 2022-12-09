// Copyright 2018 The Rind Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

// DNSServer will do Listen, Query and Send.
type DNSServer interface {
	Listen()
	Query(Packet)
}

// DNSService is an implementation of DNSServer interface.
type DNSService struct {
	conn      *net.UDPConn
	clients   clientList
	forwarder net.UDPAddr
}

// Packet carries DNS packet payload and sender address.
type Packet struct {
	addr    net.UDPAddr
	message dnsmessage.Message
}

const (
	// DNS packet max length
	packetLen int = 512
)

// Listen starts a DNS server on port 53
func (s *DNSService) Listen(port int) {
	var err error
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer s.conn.Close()

	for {
		buf := make([]byte, packetLen)
		_, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		var m dnsmessage.Message
		err = m.Unpack(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(m.Questions) == 0 {
			continue
		}
		go s.Query(Packet{*addr, m})
	}
}

// Query lookup answers for DNS message.
func (s *DNSService) Query(p Packet) {
	// got response from forwarder, send it back to client
	if p.message.Header.Response {
		pKey := pString(p)
		if addrs, ok := s.clients.get(pKey); ok {
			for _, addr := range addrs {
				go s.sendPacket(p.message, addr)
			}
			s.clients.remove(pKey)
		}
		return
	}

	// forwarding
	s.clients.set(pString(p), p.addr)
	go s.sendPacket(p.message, s.forwarder)
}

func (s *DNSService) sendPacket(message dnsmessage.Message, addr net.UDPAddr) {
	packed, err := message.Pack()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = s.conn.WriteToUDP(packed, &addr)
	if err != nil {
		log.Println(err)
	}
}

// New setups a DNSService, rwDirPath is read-writable directory path for storing dns records.
func NewDNSService(forwarder net.UDPAddr) *DNSService {
	return &DNSService{
		clients:   clientList{data: make(map[string][]net.UDPAddr)},
		forwarder: forwarder,
	}
}
