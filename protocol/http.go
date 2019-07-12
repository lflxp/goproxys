package protocol

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// http protocol handler
// 参考资料：https://www.flysnow.org/2016/12/24/golang-http-proxy.html
func HandleHttpRequestTCP(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	// 创建空字节数组
	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("Byte Origin: %s", string(b[:]))
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Error(err)
		return
	}

	if hostPortURL.Opaque == "443" { // https 访问
		address = hostPortURL.Scheme + ":443"
	} else { // http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { // host 不带端口，默认80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	log.Debugf("结果: %s %s %s\n", address, method, host)
	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Error(err)
		return
	}

	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(b[:n])
	}

	// 进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
}
