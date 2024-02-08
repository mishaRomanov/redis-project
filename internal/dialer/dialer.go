package dialer

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// SendOrder call client api and sends order there
// WARNING ! DEFAULT CLIENT PORT IS 3030
func SendOrder(body io.Reader) error {
	resp, err := http.Post("http://127.0.0.1:3030/place-order", "application/json", body)
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
	resp, err := http.Post("http://127.0.0.1:8080/delete-order", "application/json", body)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
