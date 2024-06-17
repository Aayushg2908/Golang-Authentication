package main

import (
	"golang-authentication/handlers"
	"golang-authentication/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(helpers.CORSMiddleware())

	r.POST("/api/register", handlers.Register)
	r.POST("/api/login", handlers.Login)
	r.GET("/api/me", helpers.AuthMiddleware(), handlers.Me)

	defer helpers.CloseClient(helpers.Client)

	r.Run()
}
