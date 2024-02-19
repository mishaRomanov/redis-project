package main

import (
	"github.com/mishaRomanov/redis-project/internal/config"
	"github.com/mishaRomanov/redis-project/internal/handlers"
	"github.com/mishaRomanov/redis-project/internal/storage"

	//
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Println("server up....")
	//creating a service
	service := echo.New()

	//parsing env variables into config
	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	//creating a redis instance
	dbredis, redisError := storage.NewInstance(cfg.Port, cfg.Password, cfg.DB)
	if redisError != nil {
		logrus.Fatalf("%v", redisError)
	}

	//creating handler instance by inserting database object inside
	handlerService := handlers.NewHandler(dbredis)

	//handles get /info
	service.GET("/info", handlerService.Info)

	//handles new order insertion
	service.POST("/order", handlerService.NewOrder)

	//handles order deletion
	service.DELETE("/order/:id", handlerService.CloseOrder)

	//starting a service and catching error
	serviceError := service.Start(":8080")
	if serviceError != nil {
		logrus.Fatalf("%v", serviceError)
	}
}
