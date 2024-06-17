package main

import (
	"golang-authentication/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(helpers.CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.Run()
}
