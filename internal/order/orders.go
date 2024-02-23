package order

import (
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/entities"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// SendOrder call client api and sends order there
// WARNING ! DEFAULT CLIENT PORT IS 3030
func SendOrder(body io.Reader) error {
	resp, err := http.Post("http://client:3030/add", "application/json", body)
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
	host := &http.Client{}
	logrus.Printf("Received request to close order %s", id)
	url := fmt.Sprintf("http://server:8080/client/order/%s", id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return err
	}
	req.Header.Add("Authorization", entities.Token)
	resp, err := host.Do(req)
	if err != nil {
		logrus.Errorf("%v\n", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
