package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/entities"
	"github.com/mishaRomanov/redis-project/internal/order"
	"github.com/mishaRomanov/redis-project/internal/storage"
	"io"
	"net/http"
	//
	"github.com/labstack/echo/v4"
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
	data, _ := ParseBody(ctx)

	//writing the order to the database and creating an orderID

	orderID, err := h.redis.NewOrder(data.Description)
	if err != nil {
		logrus.Error(err)
		return ctx.String(http.StatusInternalServerError, "Error while writing values to redis")
	}

	//re-send the json body to client side and checking whether we get an error or not
	err = order.SendOrder(entities.CreateOrderBody(orderID, data.Description))
	if err != nil {
		logrus.Errorf("%v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Order %s created.", orderID))
}

// Deletes the order
// /////////////////
// Тут доступ должен быть только со стороны клиента
// мб сделать авторизацию через токен? хз короче надо подумать. но короче этот эндпоинт не должен быть доступен извне
// он должен вызываться только через хендлер клиента
func (h *Handler) CloseOrder(ctx echo.Context) error {
	logrus.Infof("New request to delete order")
	orderID := ctx.Param("id")
	if err := h.redis.CloseOrder(orderID); err != nil {
		logrus.Errorf("%v\n", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Order %s closed.", orderID))
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

// client add
func ClientHandlerAdd(ctx echo.Context) error {
	data, body := ParseBody(ctx)
	if body == nil {
		return ctx.String(http.StatusBadRequest, "Error while reading json body")
	}
	logrus.Infof("New order received: %s: %s\n", data.ID, data.Description)
	//adding new order to the orders map
	entities.OrdersMap[data.ID] = data.Description
	logrus.Println(entities.OrdersMap)
	return ctx.String(http.StatusOK, "OK")
}

// client delete
func ClientHandlerDelete(ctx echo.Context) error {
	data, body := ParseBody(ctx)
	if body == nil {
		return ctx.String(http.StatusBadRequest, "Error while reading json body")
	}
	if err := order.CloseOrder(data.ID); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, "OK")
}

// Parsebody.... parses the json body!! Wow!!
func ParseBody(ctx echo.Context) (entities.OrderReceiver, []byte) {
	//creating a request body struct piece
	data := entities.OrderReceiver{}
	//defer the closure of the body
	defer ctx.Request().Body.Close()
	//reading json body
	r, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorf("error while reading the body: %\n", err)
		return entities.OrderReceiver{}, nil
	}
	//parsing json
	err = json.Unmarshal(r, &data)
	if err != nil {
		logrus.Errorf("error while parsing json: %v\n", err)
		return entities.OrderReceiver{}, nil
	}
	return data, r
}
