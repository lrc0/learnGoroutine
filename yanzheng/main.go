package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)
	go setData(ch)
	fmt.Println("1:", <-ch)
	fmt.Println("2:", <-ch)
	fmt.Println("3:", <-ch)
	fmt.Println("4:", <-ch)
	fmt.Println("5:", <-ch)
}
func setData(ch chan string) {
	ch <- "test"
	ch <- "hello wolrd"
	ch <- "123"
	ch <- "456"
	ch <- "789"
}
