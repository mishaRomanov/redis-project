package main

import (
	"encoding/json"
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/dialer"
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
	//creating a map storage for orders
	ordersMap = make(map[string]string)
)

// func that listens to input
func listen() {
	for {
		var a string
		fmt.Scanln(&a)
		switch {
		//list all orders
		case a == "/all":
			if len(ordersMap) == 0 {
				logrus.Println("no active orders found.")
			}
			for k, v := range ordersMap {
				//printing all orders
				logrus.Printf("Order %s: %s\n", k, v)
			}
		//close the order
		case a == "/done":
			logrus.Print("ok. which one is done? the id:")
			fmt.Scan(&a)
			err := dialer.CloseOrder(a)
			if err != nil {
				logrus.Errorf("error while closing given order: %v\n", err)
			}
			//deleting the order from map
			delete(ordersMap, a)
		}
	}
}

// this is client that receives and displays orders
func main() {
	logrus.Println("client up....")
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
		ordersMap[req.OrderID] = req.OrderDESC

		return ctx.String(http.StatusOK, "OK")
	})

	//starting client
	clientError := client.Start(":3030")
	if clientError != nil {
		logrus.Fatalf("%v\n", clientError)
	}
}
