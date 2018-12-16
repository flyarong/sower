package proxy

import (
	"net"
	"strings"

	"github.com/golang/glog"
	"github.com/wweir/sower/parse"
	"github.com/wweir/sower/proxy/quic"
)

type Server interface {
	Listen(port string) (<-chan net.Conn, error)
}

func StartServer(netType, port string) {
	var server Server
	switch netType {
	case QUIC.String():
		server = quic.NewServer()
	case KCP.String():
	}

	if port == "" {
		glog.Fatalln("port must set")
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	connCh, err := server.Listen(port)
	if err != nil {
		glog.Fatalf("listen %v fail: %s", port, err)
	}

	for {
		conn := <-connCh
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	conn, addr, err := parser.ParseAddr(conn)
	if err != nil {
		glog.Warningln(err)
		return
	}
	glog.V(1).Infof("new conn from %s to %s", conn.RemoteAddr(), addr)

	rc, err := net.Dial("tcp", addr)
	if err != nil {
		glog.Warningln(err)
		return
	}
	defer rc.Close()
	if err := rc.(*net.TCPConn).SetKeepAlive(true); err != nil {
		glog.Warningln(err)
	}
	relay(rc, conn)
}
