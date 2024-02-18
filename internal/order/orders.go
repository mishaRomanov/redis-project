package order

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	//
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// OrdersMap stores orders in map so client can access them
var OrdersMap = make(map[string]string)

// OrderReceiver struct is used to parse order data from request body
type OrderReceiver struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

// Parsebody.... parses the json body!! Wow!!
func ParseBody(ctx echo.Context) (OrderReceiver, []byte) {
	//creating a request body struct piece
	data := OrderReceiver{}
	//defer the closure of the body
	defer ctx.Request().Body.Close()
	//reading json body
	r, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorf("error while reading the body: %\n", err)
		return OrderReceiver{}, nil
	}
	//parsing json
	err = json.Unmarshal(r, &data)
	if err != nil {
		logrus.Errorf("error while parsing json: %v\n", err)
		return OrderReceiver{}, nil
	}
	return data, r
}

// SendOrder call client api and sends order there
// WARNING ! DEFAULT CLIENT PORT IS 3030
func SendOrder(body io.Reader) error {
	resp, err := http.Post("http://client:3030/place-order", "application/json", body)
	if err != nil {
		logrus.Errorf("Error while sending order to client: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

// CloseOrder calls server api for order closure
// WARNING ! DEFAULT SERVER PORT IS 8080
func CloseOrder(id string) error {
	logrus.Println("Received request to close order --- log from order")
	body := strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id))
	resp, err := http.NewRequest(http.MethodDelete, "http://server:8080/order", body)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
