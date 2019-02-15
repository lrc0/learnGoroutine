package main

import (
	"fmt"
	"io"
	"net"
	"os"
	// "strconv"
	// "strings"

	"gopkg.in/logger.v1"
)

func main() {
	ch := make(chan int)
	host := os.Args[1]
	port := os.Args[2]
	filename := os.Args[3]
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
		return
	}
	go sendPackage(conn, filename, ch)
	<-ch
}

func sendPackage(conn net.Conn, filename string, ch chan int) {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Error(err)
		return
	}

	//发送文件名到服务端
	buf := make([]byte, 4096)
	conn.Write([]byte(fileInfo.Name()))

	//接收服务端返回的OK
	n, err := conn.Read(buf)
	if err != nil {
		log.Error(err)
		return
	}
	if string(buf[:n]) != "ok" {
		log.Fatal("Server is not receive the filename")
	}

	var total int
	//开始正式传输文件内容
	for {
		buf := make([]byte, 4096)
		n, err := file.Read(buf)
		if err == io.EOF && n == 0 {
			ch <- 1
			break
		}
		if err != nil {
			log.Error(err)
			return
		}
		total += n

		per := (float64(total) / float64(fileInfo.Size())) * 100
		fmt.Printf("\r[%s] %s%s", bar(per, 100), fmt.Sprintf("%.2f", per), fmt.Sprintf("%s", "% "))
		conn.Write(buf[:n])
	}
	fmt.Println()
}

func bar(count, size float64) string {
	str := ""
	for i := float64(0); i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return str
}
