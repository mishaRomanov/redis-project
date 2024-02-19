package entities

import (
	"fmt"
	"io"
	"strings"
)

// OrdersMap stores orders in map so client can access them
var OrdersMap = make(map[string]string)

// OrderReceiver struct is used to parse order data from request body
type OrderReceiver struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateOrderBody is a helper function that creates a Reader that can be sent to client
func CreateOrderBody(id, description string) io.Reader {
	return strings.NewReader(fmt.Sprintf(`{"id":"%s","description":"%s"}`, id, description))
}
