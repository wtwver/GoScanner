package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/getResult", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "result",
		})
	})

	r.Run(":80") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
