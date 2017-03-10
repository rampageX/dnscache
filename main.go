package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	//MaxCaches allowed how many caches
	MaxCaches = 8192
	//Timeout set the limit of time out
	Timeout = 30
)

var (
	//Proto set the protocol type using
	Proto = "tcp"
	//NsAddrs Nameserver addresses
	NsAddrs = []string{
		"208.67.222.222:443",
		"208.67.220.220:443",
	}
)

func main() {
	var host string
	var port int
	args := os.Args[1:]
	argslen := len(args)
	//fmt.Println(argslen)
	host = "127.0.0.1"
	port = 53
	if argslen >= 1 {
		host = args[0]
	}
	if argslen >= 2 {
		newport, err := strconv.Atoi(args[1])
		if err == nil {
			port = newport
		}
	}
	if argslen >= 3 {
		Proto = args[2]
	}
	// check network  is online

	for {
		conn, err := net.Dial("tcp", "208.67.222.222:443")
		if err != nil {
			fmt.Println("Failed to connect network, will sleep for 5s")
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		} else {
			conn.Close()
			fmt.Println("Success connect to 208.67.222.222:443")
			break
		}
	}
	server := &Server{
		host:     host,
		port:     port,
		rTimeout: Timeout * time.Second,
		wTimeout: Timeout * time.Second,
	}

	server.Run()
	fmt.Println("DNS server start")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			fmt.Println("Signal received, now stop and exit")
			break forever
		}
	}

}
