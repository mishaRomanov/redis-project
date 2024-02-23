package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mishaRomanov/redis-project/internal/entities"
	"github.com/mishaRomanov/redis-project/internal/handlers"
	"github.com/sirupsen/logrus"
)

// this is client that receives and displays orders
func main() {
	logrus.Println("client up....")

	//here we send a request to create a jwt token that client side will use
	handlers.SendRequestToAuthAndWriteToken()

	client := echo.New()
	logrus.Infof("Here is your token: %s\n", entities.Token)

	//handles order placement
	client.POST("/add", handlers.ClientHandlerAdd)

	//handles the delete functionality
	client.GET("/close/:id", handlers.ClientHandlerDelete)

	//starting client
	clientError := client.Start(":3030")
	if clientError != nil {
		logrus.Fatalf("%v\n", clientError)
	}
}
