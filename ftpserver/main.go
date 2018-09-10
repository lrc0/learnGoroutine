package main

import (
	"fmt"
	// "io"
	"net"
	"os"
	"time"

	"gopkg.in/logger.v1"
)

var connPool = make(map[string]net.Conn)

func main() {
	listener, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Error(err)
		return
	}
	defer listener.Close()
	log.Info("waiting for client.....")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		conn.SetDeadline(time.Now().Add(time.Second * time.Duration(60)))
		log.Info(conn.RemoteAddr().String(), "=====>> connection success")
		connPool[conn.RemoteAddr().String()] = conn

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// defer conn.Close()
	buffer := make([]byte, 2048)

	//读信息从conn
	for {
		mess := make(chan int)
		go heartbeat(conn, mess)

		b, err := conn.Read(buffer)
		if err != nil {
			log.Error("ERR: ", err)
			conn.Close()
			break
		}
		broadCast(buffer[:b], conn.RemoteAddr().String())
		go getMassage(b, mess)

		fmt.Fprintf(os.Stdout, "%s: %s", conn.RemoteAddr().String(), string(buffer[:b]))
	}
}

func heartbeat(conn net.Conn, mess chan int) {
	select {
	case <-mess:
		conn.SetDeadline(time.Now().Add(time.Duration(60) * time.Second))
		// fmt.Println("Add more time")
	case <-time.After(59 * time.Second):
		log.Info("close connection for 60 second without input")
		conn.Close()
	}
}

func getMassage(bytes int, mess chan int) {
	if bytes != 0 {
		mess <- bytes
	}
	close(mess)
}

//广播
func broadCast(msg []byte, name string) {
	for k, con := range connPool {
		if k != name {
			con.Write([]byte(msg))
		}
	}
}
