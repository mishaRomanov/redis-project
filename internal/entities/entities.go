package entities

// OrdersMap stores orders in map so client can access them
var (
	OrdersMap = make(map[string]string)
	Token     string
)

// OrderReceiver struct is used to parse order data from request body
type OrderReceiver struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

// AuthBody struct to parse the body to
type AuthBody struct {
	Num             int  `json:"num,omitempty"`
	IsClientService bool `json:"isClientService"`
}

// AuthResponse struct to send token with
type AuthResponse struct {
	Token string `json:"token,omitempty"`
}

// WriteToken writes jwt token to the builder
func WriteToken(token string) string {
	return "Bearer " + token
}
