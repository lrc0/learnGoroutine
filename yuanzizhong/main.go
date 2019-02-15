package main

import (
	"fmt"
	"gopkg.in/logger.v1"
	"io"
	"net/http"
)

func main() {
	res, err := http.Get("https://time.is")
	if err != nil {
		log.Error(err)
		return
	}

	defer res.Body.Close()
	for {
		buf := make([]byte, 4096)
		n, err := res.Body.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println("result: ", string(buf[:n]))
	}
}
