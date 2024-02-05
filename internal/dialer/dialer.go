package dialer

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func SendOrder(body io.Reader) error {
	resp, err := http.Post("http://127.0.0.1:3030/place-order", "application/json", body)
	if err != nil {
		logrus.Errorf("Error while sending order to client: %v\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		return errors.New("error: status bad request found. check everything")
	}
	return nil
}
