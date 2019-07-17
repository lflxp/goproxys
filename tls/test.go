package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	l, err := tls.Listen("tcp", ":33333", config)
	defer l.Close()
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleTCPRequest(client)
	}
}
func handleTCPRequest(client net.Conn) {
	remote, err := net.Dial("tcp", "localhost:12333") // connect to Redis
	if err != nil {
		log.Println(err)
		return
	}
	go io.Copy(remote, client) // copy client data to database
	io.Copy(client, remote)    // copy database data to client
}
