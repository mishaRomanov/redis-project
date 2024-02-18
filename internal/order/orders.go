package dialer

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// OrderReceiver struct is used to parse order data from request body
type OrderReceiver struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
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
	logrus.Println("Received request to close order --- log from dialer")
	body := strings.NewReader(fmt.Sprintf(`{"order-id":"%s"}`, id))
	resp, err := http.Post("http://server:8080/delete-order", "application/json", body)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
