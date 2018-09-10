package main

import (
	"fmt"
	"net"
	"time"

	"gopkg.in/logger.v1"
)

//客户端
var clientMap = make(map[string]*net.TCPConn)

func listenClient(addr string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	tcpListen, _ := net.ListenTCP("tcp", tcpAddr)
	for { //不停地接收
		clientConn, _ := tcpListen.AcceptTCP()                   //监听请求连接
		clientMap[clientConn.RemoteAddr().String()] = clientConn //
		go addReceiver(clientConn)
		fmt.Println("用户: ", clientConn.RemoteAddr().String(), " 已连接")
	}
}

//向连接添加接收器
func addReceiver(conn *net.TCPConn) {
	for {
		bytes := make([]byte, 2048)
		len, err := conn.Read(bytes)
		if err != nil {
			log.Error(err)
			conn.Close()
			return
		}
		fmt.Println(string(bytes[:len]))
		msgBroadcast(bytes[:len], conn.RemoteAddr().String())
	}
}

//广播给所有 client
func msgBroadcast(msg []byte, name string) {
	for k, conn := range clientMap {
		if k != name {
			conn.Write(msg)
		}
	}
}

//主函数
func main() {
	log.Info("服务启动...")
	time.Sleep(1 * time.Second)
	log.Info("等待客户端连接...")
	go listenClient("127.0.0.1:8099")
	select {}
}
