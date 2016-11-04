package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	MAX_CACHES = 8192
	TIMEOUT    = 30
)

var (
	PROTO   = "tcp"
	NSADDRS = []string{
		"208.67.222.222:443",
		"208.67.220.220:443",
		"216.146.35.35:53",
		"216.146.36.36:53",
	}
)

func main() {
	var host string
	var port int
	args := os.Args[1:]
	argslen := len(args)
	fmt.Println(argslen)
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
		PROTO = args[2]
	}

	server := &Server{
		host:     host,
		port:     port,
		rTimeout: TIMEOUT * time.Second,
		wTimeout: TIMEOUT * time.Second,
	}

	server.Run()
	fmt.Println("DNS server start")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			fmt.Println("Signal recieved, now stop and exit")
			break forever
		}
	}

}
