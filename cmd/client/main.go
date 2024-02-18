package main

import (
	"github.com/mishaRomanov/redis-project/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// this is client that receives and displays orders
func main() {
	logrus.Println("client up....")
	client := echo.New()

	//handles order placement
	client.POST("/add", handlers.ClientHandlerAdd)

	client.POST("/close", handlers.ClientHandlerDelete)

	//starting client
	clientError := client.Start(":3030")
	if clientError != nil {
		logrus.Fatalf("%v\n", clientError)
	}
}
