package main

import (
	"net"

	"github.com/lflxp/goproxys/protocol"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	// "github.com/eahydra/socks/cmd/socksd"
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

	Mysql = kingpin.Flag(
		"mysql",
		"mysql proxy代理，无加密,负载均衡",
	).Bool()

	Socket5 = kingpin.Flag(
		"socket5",
		"socket5 proxy代理，无验证",
	).Bool()

	Socket5Cipher = kingpin.Flag(
		"sc",
		"socket5 cipher加密代理",
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
	// cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// config := &tls.Config{Certificates: []tls.Certificate{cer}}

	l, err := net.Listen("tcp", ":8081")
	// l, err := tls.Listen("tcp", ":8081", config)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Started Proxy")

	if *Http {
		log.Println("Http Proxy Listening port: 8081")
	} else if *Socket5 {
		log.Println("Socket5 Proxy Listening port: 8081")
	} else if *Mysql {
		log.Println("Mysql Proxy Listening port: 8081")
	} else {
		log.Println("Socket5 Proxy Listening port: 8081")
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		log.Println(client.RemoteAddr().String())

		if *Http {
			go protocol.HandleHttpRequestTCP(client)
		} else if *Socket5 {
			go protocol.HandleSocket5RequestTCP(client)
		} else if *Mysql {
			go protocol.HandleMysqlRequestTCP(client)
		} else {
			go protocol.HandleSocket5RequestTCP(client)
		}
		// else if *Socket5Cipher {
		// 	cipherDecorator := socksd.NewCipherConnDecorator(conf.Crypto, conf.Password)
		// 	listener = socksd.NewDecorateListener(listener, cipherDecorator)
		// 	socks5Svr, err := socks.NewSocks5Server(forward)
		// 	if err != nil {
		// 		listener.Close()
		// 		ErrLog.Println("socks.NewSocks5Server failed, err:", err)
		// 		return
		// 	}
		// 	go func() {
		// 		defer listener.Close()
		// 		socks5Svr.Serve(listener)
		// 	}()
		// }
	}
}
