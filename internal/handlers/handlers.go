package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type handler struct {
	redis *redis.Client
}

// поля делаем экспортируемые чтобы json.Unmarshall смог распарсить тело запроса
type requestBody struct {
	TaskID      string `json:"task_id"`
	Description string `json:"description"`
}

// Info handles /about request
func (h *handler) Info(ctx echo.Context) error {
	logrus.Infoln("New request ")
	return ctx.String(http.StatusOK, "Hello world! Handler works.")
}

// InsertValue handles /add POST request
func (h *handler) InsertValue(ctx echo.Context) error {
	logrus.Infof("New /add request")
	//creating a requestbody struct piece
	data := requestBody{}

	//defer the closure of the body
	defer ctx.Request().Body.Close()

	//reading json body
	r, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorf("error while reading the body: %\n", err)
		return ctx.String(http.StatusUnprocessableEntity, "")
	}

	//parsing json
	err = json.Unmarshal(r, &data)
	if err != nil {
		logrus.Errorf("error while parsing json: %v\n", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	//logging
	logrus.Infof("%s\t%s", data.Description, data.TaskID)

	return ctx.String(http.StatusOK, fmt.Sprintf("%s\t%s", data.Description, data.TaskID))
}

// NewHandler  creates handler instance
func NewHandler(client *redis.Client) *handler {
	instance := handler{
		redis: client,
	}
	return &instance
}
