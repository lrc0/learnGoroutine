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
	listener, err := net.Listen("tcp", "localhost:"+*p)
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

// func usa() {
// 	u, err := time.LoadLocation("US/Eastern")
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	usa := now.In(u).Format(layout)
// 	// fmt.Println("usa: ", usa)
// 	ch <- usa
// }

// func london() {
// 	t, err := time.LoadLocation("Europe/London")
// 	london := now.In(t).Format(layout)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}

// 	// fmt.Println("london: ", london)
// 	ch <- london
// }
