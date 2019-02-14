package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	for i := 10; i <= 100; i += 10 {
		str := "[" + bar(i/10, 10) + "] " + strconv.Itoa(i) + "%"
		fmt.Printf("\r%s", str)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("")
}

func bar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return str
}
