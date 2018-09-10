package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"gopkg.in/logger.v1"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8001")
	if err != nil {
		log.Error(err)
		os.Exit(-1)
		return
	}
	go sendpackage(conn)
	go processRecvData(conn)
	processSendData(conn)
}

func processRecvData(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		b, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Error("ERR: ", err)
				conn.Close()
				os.Exit(1)
			}
		}
		fmt.Fprintf(os.Stdout, "%s: %s", login(), string(buffer[:b]))
	}
}

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

func sendpackage(conn net.Conn) {
	for i := 1; ; i++ {
		mass := strconv.Itoa(i) + " times" + "===> send heart beat from client: " + login() + "\n"
		conn.Write([]byte(mass))
		time.Sleep(55 * time.Second)
	}
}

func login() string {
	return os.Args[1]
}
