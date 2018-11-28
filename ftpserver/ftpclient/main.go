package main

import (
	"bufio"
	"errors"
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
	name, err := login()
	if err != nil {
		conn.Close()
		return
	}
	buffer := make([]byte, 2048)
	for {
		b, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				msg := fmt.Sprintf("%s 退出聊天室", name)
				conn.Write([]byte(msg))
				// conn.Close()
				os.Exit(1)
			}
		}
		fmt.Fprintf(os.Stdout, "%s: %s", name, string(buffer[:b]))
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
	name, err := login()
	if err != nil {
		conn.Close()
		return
	}
	for i := 1; ; i++ {
		mass := strconv.Itoa(i) + " times" + "===> send heart beat from client: " + name + "\n"
		conn.Write([]byte(mass))
		time.Sleep(55 * time.Second)
	}
}

func login() (string, error) {
	slice := []string{
		"ruicai",
		"baiwei",
		"xiaoyimei",
	}

	taget := os.Args[1]
	m := make(map[string]bool)

	for _, name := range slice {
		m[name] = true
	}

	if m[taget] == true {
		log.Info("登录成功")
	} else {
		log.Error("登录失败")
		os.Exit(-1)
		return "", errors.New("登录失败")
	}

	return taget, nil
}
