package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	//
	"github.com/labstack/echo/v4"
	"github.com/mishaRomanov/redis-project/internal/dialer"
	"github.com/mishaRomanov/redis-project/internal/storage"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	redis storage.Storager
}

// function that parses json body into struct. return struct and the body itself
func parseBody(ctx echo.Context) (RequestBody, []byte) {
	//creating a request body struct piece
	data := RequestBody{}
	//defer the closure of the body
	defer ctx.Request().Body.Close()
	//reading json body
	r, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorf("error while reading the body: %\n", err)
		return RequestBody{}, nil
	}
	//parsing json
	err = json.Unmarshal(r, &data)
	if err != nil {
		logrus.Errorf("error while parsing json: %v\n", err)
		return RequestBody{}, nil
	}
	return data, r
}

// поля делаем экспортируемые чтобы json.Unmarshall смог распарсить тело запроса
type RequestBody struct {
	OrderID   string `json:"order-id"`
	OrderDESC string `json:"order-desc,omitempty"`
}

// Info handles /about request
func (h *Handler) Info(ctx echo.Context) error {
	logrus.Infoln("New request")
	return ctx.String(http.StatusOK,
		`Hello! This service lets you create and track orders. 
Make a POST request to /new-order to create an order.
Use json and "order-id" and "description" fields.`)
}

// NewOrder handles /new-order POST request
func (h *Handler) NewOrder(ctx echo.Context) error {
	logrus.Infof("New order POST request")
	//creating a request body struct piece
	data, body := parseBody(ctx)
	//checking whether we already have that order id in redis db
	//if do, return 400
	if status := h.redis.LookUp(data.OrderID); status {
		return ctx.String(http.StatusBadRequest, "Order with such ID already exists.")
	}

	//logging
	logrus.Infof("added key value pair: %s:\t%s", data.OrderID, data.OrderDESC)
	err := h.redis.NewOrder(data.OrderID, data.OrderDESC)
	if err != nil {
		logrus.Error(err)
		return ctx.String(http.StatusInternalServerError, "Error while writing values to redis")
	}

	//re-send the json body to client side and checking whether we get an error or not
	err = dialer.SendOrder(strings.NewReader(string(body)))
	if err != nil {
		logrus.Errorf("%v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("%s\t%s", data.OrderDESC, data.OrderID))
}

// CloseOrder handler closes the given order
func (h *Handler) CloseOrder(ctx echo.Context) error {
	logrus.Println("Received request to close order --- log from handler")
	data, _ := parseBody(ctx)
	if err := h.redis.CloseOrder(data.OrderID); err != nil {
		logrus.Errorf("%v\n", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.String(http.StatusOK, "order closed")
}

// NewHandler  creates handler instance
func NewHandler(redisStorager storage.Storager) *Handler {
	instance := Handler{
		redis: redisStorager,
	}
	return &instance
}
