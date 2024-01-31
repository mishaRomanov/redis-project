package main

import (
	"context"
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/handlers"
	//
	"github.com/labstack/echo/v4"
	"github.com/mishaRomanov/redis-project/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	//creating a service
	service := echo.New()

	//creating context
	ctx := context.Background()

	//parsing env variables into config
	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	//creating a redis instance
	dbredis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// testing whether redis connection was successful or not
	pong, err := dbredis.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	logrus.Println(pong)

	//creating handler instance by inserting database object inside
	handlerService := handlers.NewHandler(dbredis)
	//handle get /about
	service.GET("/about", handlerService.Info)

	//handle POST /add
	service.POST("/add", handlerService.InsertValue)

	//starting a service and catching error
	serviceError := service.Start(":8080")
	if serviceError != nil {
		logrus.Fatalf("%v", serviceError)
	}
}
