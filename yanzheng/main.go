// package main

// import (
// 	"fmt"
// 	// "github.com/zhenorzz/snowflake"
// )

// func main() {

// 	// // Create a new Node with a Node number of 1
// 	// sf, err := snowflake.New(0)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// // Generate a snowflake ID.
// 	// uuid, _ := sf.Generate()

// 	// // Print
// 	// fmt.Println(uuid)
// }

package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")
		status := c.DefaultPostForm("status", "anonymous")

		c.JSON(200, gin.H{
			"status":  status,
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}
