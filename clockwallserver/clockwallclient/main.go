package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"gopkg.in/logger.v1"
)

var ch = make(chan string, 3)

func main() {
	slice := os.Args[1:]

	for _, zone := range slice {
		var p string
		switch {
		case zone == "Local":
			p = "8001"
		case zone == "US/Eastern":
			p = "8002"
		case zone == "Europe/London":
			p = "8003"
		}

		go func() {
			conn := connect(p)
			defer conn.Close()
			buf := bufio.NewScanner(conn)
			for buf.Scan() {
				ch <- buf.Text()
			}
		}()
	}

	for {
		fmt.Println("-----------------------------------------")
		fmt.Println("| ", <-ch, "|")
		fmt.Println("-----------------------------------------")
	}
}
func connect(p string) net.Conn {
	conn, err := net.Dial("tcp", "0.0.0.0:"+p)
	if err != nil {
		log.Error(err)
		return nil
	}
	return conn
}
