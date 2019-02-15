package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"gopkg.in/logger.v1"
)

var layout = "2006-01-02 15:04:05"
var p = flag.String("port", os.Args[1], "port")
var zone = flag.String("zone", os.Args[2], "zone")

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", "0.0.0.0:"+*p)
	if err != nil {
		log.Error(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		go china(conn)
	}
}

func china(conn net.Conn) {
	c, err := time.LoadLocation(*zone)
	if err != nil {
		log.Error(err)
		return
	}
	for {
		fmt.Fprintln(conn, *zone+"====>"+time.Now().In(c).Format(layout))
		time.Sleep(1 * time.Second)
	}
}
