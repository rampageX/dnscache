package main

import (
	"github.com/ztrue/tracerr"

	"os"
	"os/signal"
	"strconv"
	"time"
	"github.com/rs/zerolog/log"

)

const (
	//MaxCaches allowed how many caches
	MaxCaches = 64
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


/**
 * Logging error
 */
func LogError(err error) {
	log.Error().Msg(tracerr.Sprint(tracerr.Wrap(err)))
}

/**
 * Logging info
 */
func LogInfo(msg string) {
	log.Info().Msg(msg)
}
/**
 * Logging info sprintf
 */
func LogInfoF(fms string, msg ...interface{}) {
	log.Info().Msgf(fms, msg)
}

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
	LogInfo("DNS server started")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, os.Kill)

forever:
	for {
		select {
		case <-sig:
			LogInfo("Signal received, now stop and exit")
			break forever
		}
	}

}
