package main

import (
	"flag"
	"github.com/ztrue/tracerr"

	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"time"
)

const (
	//MaxCaches allowed how many caches
	MaxCaches = 64
	//Timeout set the limit of time out
	Timeout = 4
)

var (
	//Proto set the protocol type using
	Proto = flag.String("protocol", "tcp", "查询协议")
	//NsAddrs Nameserver addresses
	NsAddrs = []string{
		"223.5.5.5:53",
		"223.6.6.6:53",
	}
	CnNsAddrs = []string{
		"223.5.5.5:53",
		"223.6.6.6:53",

	}
	WorldNsAddrs = []string{
		"8.8.8.8:53",
		"8.8.4.4:53",
	}

	AreaZone = flag.String("z", "cn", "区域：cn代表国内，world代表世界")
	ListenOn = flag.String("l", "127.0.0.1:53", "监听地址和端口")
	Help = flag.Bool("h", false, "帮助")
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
	flag.Parse()
	if *Help == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// check network  is online

	if *AreaZone != "cn" {
		NsAddrs = WorldNsAddrs
	}
	server := &Server{
		listenOn:     *ListenOn,
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
