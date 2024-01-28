package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mishaRomanov/redis-project/config"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

func main() {
	//creating a service
	service := echo.New()
	//creating context
	ctx := context.Background()

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	//creating a redis instance
	dbredis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	pong, err := dbredis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println(pong)

	service.GET("/about", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello world!")
	})
	fatal_error := service.Start(":8080")
	if fatal_error != nil {
		log.Fatalf("%v", fatal_error)
	}
}
