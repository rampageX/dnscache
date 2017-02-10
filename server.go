package main

import (
	"strconv"
	"time"

	"fmt"

	"github.com/miekg/dns"
)

type Server struct {
	host     string
	port     int
	rTimeout time.Duration
	wTimeout time.Duration
}

func (s *Server) Addr() string {
	return s.host + ":" + strconv.Itoa(s.port)
}

func (s *Server) Run() {

	Handler := NewHandler()

	go Handler.PreparePool()

	//fmt.Println(Handler.resolver.NameserversPool)

	tcpHandler := dns.NewServeMux()
	tcpHandler.HandleFunc(".", Handler.DoTCP)

	udpHandler := dns.NewServeMux()
	udpHandler.HandleFunc(".", Handler.DoUDP)

	tcpServer := &dns.Server{Addr: s.Addr(),
		Net:          "tcp",
		Handler:      tcpHandler,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout,
	}

	udpServer := &dns.Server{Addr: s.Addr(),
		Net:          "udp",
		Handler:      udpHandler,
		UDPSize:      65535,
		ReadTimeout:  s.rTimeout,
		WriteTimeout: s.wTimeout,
	}

	go s.start(udpServer)
	go s.start(tcpServer)

}

func (s *Server) start(ds *dns.Server) {

	fmt.Println("Start listener on ", ds.Net, ":", s.Addr())
	err := ds.ListenAndServe()
	if err != nil {
		fmt.Println("Start listener failed:", ds.Net, ":", s.Addr(), err.Error())
	}

}
