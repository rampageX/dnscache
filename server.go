package main

import (
	"time"

	"github.com/miekg/dns"
)

//Server the server struct
type Server struct {
	listenOn string
	rTimeout time.Duration
	wTimeout time.Duration
}

//Addr the server address
func (s *Server) Addr() string {
	return s.listenOn
}

//Run set up and running the server
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

	LogInfoF("Start listener on %s:%s", ds.Net, s.Addr())
	err := ds.ListenAndServe()
	if err != nil {
		LogInfoF("Start listener failed: %s:%s:%+v", ds.Net, s.Addr(), err)
	}

}
