package handlers

import (
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/storage"
	"net/http"
	"strings"
	//
	"github.com/labstack/echo/v4"
	"github.com/mishaRomanov/redis-project/internal/order"
	"github.com/sirupsen/logrus"
)

// Handler struct stores redis in it
type Handler struct {
	redis storage.Storager
}

// NewHandler  creates handler instance
func NewHandler(redisStorager storage.Storager) *Handler {
	instance := Handler{
		redis: redisStorager,
	}
	return &instance
}

// GeneralHandler is a general handler for all types of requests, it decides what to do with the request
func (h *Handler) GeneralHandler(ctx echo.Context) error {
	//determining the type of the request
	switch ctx.Request().Method {
	case http.MethodPost:
		return h.NewOrder(ctx)
	case http.MethodDelete:
		return h.CloseOrder(ctx)
	}
	return ctx.String(http.StatusMethodNotAllowed, "Try another method or visit /info endpoint")
}

// Info handles /about request
func (h *Handler) Info(ctx echo.Context) error {
	logrus.Infoln("New request")
	return ctx.String(http.StatusOK,
		`Hello! This service lets you create and track orders. 
Make a POST request to /order with a JSON body to create a new order. 
You can also make a DELETE request to /order with a JSON body to close an order.
For more information visit https://github.com/mishaRomanov/redis-project`)
}

// Creates a new order
func (h *Handler) NewOrder(ctx echo.Context) error {
	logrus.Infof("New order POST request")
	//creating a request body struct piece
	data, body := order.ParseBody(ctx)

	orderID, err := h.redis.NewOrder(data.Description)
	if err != nil {
		logrus.Error(err)
		return ctx.String(http.StatusInternalServerError, "Error while writing values to redis")
	}

	//re-send the json body to client side and checking whether we get an error or not
	err = order.SendOrder(strings.NewReader(string(body)))
	if err != nil {
		logrus.Errorf("%v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Order %s created.", orderID))
}

// Deletes the order
func (h *Handler) CloseOrder(ctx echo.Context) error {
	logrus.Println("Received request to close order...")
	data, _ := order.ParseBody(ctx)
	if err := h.redis.CloseOrder(data.ID); err != nil {
		logrus.Errorf("%v\n", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Order %s closed.", data.ID))
}

// client add
func ClientHandlerAdd(ctx echo.Context) error {
	data, body := order.ParseBody(ctx)
	if body == nil {
		return ctx.String(http.StatusBadRequest, "Error while reading json body")
	}
	logrus.Infof("New order received: %s: %s\n", data.ID, data.Description)
	//adding new order to the orders map
	order.OrdersMap[data.ID] = data.Description
	logrus.Println("All orders: %v", order.OrdersMap)
	return ctx.String(http.StatusOK, "OK")
}

// client delete
func ClientHandlerDelete(ctx echo.Context) error {
	data, body := order.ParseBody(ctx)
	if body == nil {
		return ctx.String(http.StatusBadRequest, "Error while reading json body")
	}
	if err := order.CloseOrder(data.ID); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, "OK")
}
