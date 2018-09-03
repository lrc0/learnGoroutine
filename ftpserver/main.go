package main

import (
	"bufio"
	"fmt"
	// "io"
	"net"
	"os"
	"time"

	"gopkg.in/logger.v1"
)

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

		go handleConn(conn)
		processSendData(conn)
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
			os.Exit(-1)
			break
		}

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

//写信息到conn
func processSendData(conn net.Conn) {
	for {
		buf := bufio.NewReader(os.Stdin)
		mas, err := buf.ReadString('\n')
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		conn.Write([]byte(mas))
	}
}
