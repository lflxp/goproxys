package main

import (
	"net"

	"github.com/lflxp/goproxys/protocol"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Debug = kingpin.Flag(
		"debug",
		"debug log level",
	).Bool()

	Http = kingpin.Flag(
		"http",
		"http proxy代理，无加密",
	).Bool()

	Socket5 = kingpin.Flag(
		"socket5",
		"socket5 proxy代理，无验证",
	).Bool()
)

func init() {
	kingpin.Version("goproxys v0.1")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	if *Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}
	log.Println("Started Proxy")

	if *Http {
		log.Println("Http Proxy Listening port: 8081")
	} else if *Socket5 {
		log.Println("Socket5 Proxy Listening port: 8081")
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		if *Http {
			go protocol.HandleHttpRequestTCP(client)
		} else if *Socket5 {
			go protocol.HandleSocket5RequestTCP(client)
		} else {
			go protocol.HandleSocket5RequestTCP(client)
		}
	}
}
