package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	//MaxCaches allowed how many caches
	MaxCaches = 128
	//Timeout set the limit of time out
	Timeout = 4
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
