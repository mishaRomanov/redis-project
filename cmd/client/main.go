package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Order struct {
	OrderID   string `json:"order-id"`
	OrderDESC string `json:"order-desc"`
}

var (
	orders []Order
)

// func that listens to input
func listen() {
	for {
		var a string
		fmt.Scanln(&a)
		switch {
		case a == "/all":
			logrus.Println(orders)
		case a == "/done":
			logrus.Println("ok. which one is done? the id")
		}
	}
}

// this is client that receives and displays orders
func main() {
	logrus.Println("type /all to see all orders.")
	client := echo.New()
	//listen to input
	go listen()

	//handles order placement
	client.POST("/place-order", func(ctx echo.Context) error {
		logrus.Infof("New order received.")
		var req Order
		defer ctx.Request().Body.Close()

		//reading json from request body
		r, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			logrus.Errorf("error while reading json body: %v\n", err)
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		//parsing json
		err = json.Unmarshal(r, &req)
		if err != nil {
			logrus.Errorf("error while parsing json: %v\n", err)
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		//adding new order to the list
		orders = append(orders, req)

		return ctx.String(http.StatusOK, "OK")
	})

	//starting client
	clientError := client.Start(":3030")
	if clientError != nil {
		logrus.Fatalf("%v\n", clientError)
	}
}
