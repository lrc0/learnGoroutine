package main

import (
	"io"
	"net"
	"os"

	"gopkg.in/logger.v1"
)

var connPool = make(map[string]*net.Conn)

func main() {
	// host := os.Args[1]
	// port := os.Args[2]
	listener, err := net.Listen("tcp", "0.0.0.0:8001")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Infof("listen on %s", "0.0.0.0:8001")
	log.Info("waiting for connect ....")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		connPool[conn.RemoteAddr().String()] = &conn
		log.Infof("%s connected", conn.RemoteAddr().String())

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 4096)
	defer conn.Close()

	//读取客户端传过来的文件名
	n, err := conn.Read(buf)
	if err != nil {
		log.Error(err)
		return
	}

	file, err := os.Create(string(buf[:n]))
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	//读取文件名成功之后，返回OK给客户端
	conn.Write([]byte("ok"))

	//开始正式接收客户端传过来的文件
	for {
		n, err := conn.Read(buf)
		if err == io.EOF && n == 0 {
			delete(connPool, conn.RemoteAddr().String())
			break
		}
		if err != nil {
			log.Error(err)
			return
		}

		file.Write([]byte(buf[:n]))
	}
}
