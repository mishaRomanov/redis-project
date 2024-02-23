package tools

import (
	"encoding/json"
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/entities"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// CreateOrderBody is a helper function that creates a Reader that can be sent to client
func CreateOrderBody(id, description string) io.Reader {
	return strings.NewReader(fmt.Sprintf(`{"id":"%s","description":"%s"}`, id, description))
}

// ParseAuthBody parses the auth request body
func ParseAuthBody(body io.Reader) entities.AuthResponse {
	var token entities.AuthResponse
	r, err := io.ReadAll(body)
	if err != nil {
		logrus.Errorf("error while reading the body: %\n", err)
		return token
	}
	err = json.Unmarshal(r, &token)
	if err != nil {
		logrus.Errorf("error while parsing json: %v\n", err)
		return token
	}
	return token
}

// ParseBody .... parses the json body!!
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
