package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//用户名
var userName string

//本机连接
var selfConn *net.TCPConn

//读取行文本
var reader = bufio.NewReader(os.Stdin)

//建立连接
func connect(addr string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		conn.Close()
		os.Exit(1)
	}
	selfConn = conn
	go msgSender()
	go msgReceiver()

}

//消息接收器
func msgReceiver() {
	buff := make([]byte, 2048)
	for {
		len, _ := selfConn.Read(buff)
		fmt.Println(string(buff[:len]))
	}
}

//消息发送器
func msgSender() {
	for {
		readLineMsg, _, _ := reader.ReadLine()
		readLineMsg = []byte(userName + " : " + string(readLineMsg))
		selfConn.Write(readLineMsg)
	}
}

//主函数
func main() {
	fmt.Println("请问你怎么称呼？")
	name, _, _ := reader.ReadLine()
	userName = string(name)
	connect("127.0.0.1:8099")
	select {}
}
