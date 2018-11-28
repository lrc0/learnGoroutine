package main

import (
	"fmt"

	"github.com/zhenorzz/snowflake"
)

func main() {

	// Create a new Node with a Node number of 1
	sf, err := snowflake.New(0)
	if err != nil {
		panic(err)
	}

	// Generate a snowflake ID.
	uuid, _ := sf.Generate()

	// Print
	fmt.Println(uuid)
}
