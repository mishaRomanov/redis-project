package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mishaRomanov/redis-project/internal/entities"
	"github.com/mishaRomanov/redis-project/internal/order"
	"github.com/mishaRomanov/redis-project/internal/storage"
	"github.com/mishaRomanov/redis-project/internal/tools"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	//
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	permissionDeniedError = errors.New("permission denied")
	secretKey             = []byte("key_for_client")
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

// NewOrder creates a new order
func (h *Handler) NewOrder(ctx echo.Context) error {
	logrus.Infof("New order POST request")
	//creating a request body struct piece
	data, _ := tools.ParseBody(ctx)

	// checking whether body is empty or not
	if data.Description == "" {
		return ctx.String(http.StatusBadRequest, "Empty description")
	}
	//writing the order to the database and creating an orderID
	orderID, err := h.redis.NewOrder(data.Description)
	if err != nil {
		logrus.Error(err)
		return ctx.String(http.StatusInternalServerError, "Error while writing values to redis")
	}

	//re-send the json body to client side and checking whether we get an error or not
	err = order.SendOrder(tools.CreateOrderBody(orderID, data.Description))
	if err != nil {
		logrus.Errorf("%v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("Order %s created.", orderID))
}

// CloseOrder closes the order in redis
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

// ClientHandlerAdd shares new orders with the client
func ClientHandlerAdd(ctx echo.Context) error {
	data, body := tools.ParseBody(ctx)
	if body == nil {
		return ctx.String(http.StatusBadRequest, "Error while reading json body")
	}
	logrus.Infof("New order received: %s: %s\n", data.ID, data.Description)
	//adding new order to the orders map
	entities.OrdersMap[data.ID] = data.Description
	logrus.Println(entities.OrdersMap)
	return ctx.String(http.StatusOK, "OK")
}

// ClientHandlerDelete deletes order through client interface
func ClientHandlerDelete(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if err := order.CloseOrder(orderID); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, "OK")
}

// RegisterClientSession makes authorization for client and returns a JWT token
func (h *Handler) RegisterClientSession(ctx echo.Context) error {
	logrus.Infoln("New request to create a token.")
	logrus.Infoln(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	data, err := checkAuthRequest(ctx)
	if err != nil {
		if errors.Is(err, permissionDeniedError) {
			return ctx.String(http.StatusForbidden, err.Error())
		}
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	// generating the payload for the token
	payload := jwt.MapClaims{
		"sub": data.Num,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	//generating the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	//signing the token
	signed, err := token.SignedString(secretKey)
	if err != nil {
		logrus.Errorf("error while signing the token: %v\n", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	//returning the token with json
	return ctx.JSON(http.StatusOK, entities.AuthResponse{Token: signed})
}

// SendAuthRequest sends request to server to authorize client
func SendAuthRequest(num int) (string, error) {
	rdr := strings.NewReader(fmt.Sprintf(`{"num":%d,"isClientService":true}`, num))
	resp, err := http.Post("http://server:8080/auth", "application/json", rdr)
	if err != nil {
		logrus.Errorf("error while sending the request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	//parsing the response body
	token := tools.ParseAuthBody(resp.Body)
	return token.Token, nil
}

// helper function to check if the request is authorized
func checkAuthRequest(ctx echo.Context) (entities.AuthBody, error) {
	//creating a struct to parse the json body to
	data := entities.AuthBody{}
	defer ctx.Request().Body.Close()
	// reading json body
	r, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Errorf("error while reading the body: %\n", err)
		return entities.AuthBody{}, err
	}
	//parsing json
	err = json.Unmarshal(r, &data)
	if err != nil {
		logrus.Errorf("error while parsing json: %v\n", err)
		return entities.AuthBody{}, err
	}
	//this is important because we here check whether the request was sent from actual client service
	//only client service knows about this json field
	if data.IsClientService == false {
		return entities.AuthBody{}, permissionDeniedError
	}
	return data, nil
}

// Helper function to send authorization request to server
func SendRequestToAuthAndWriteToken() {
	time.Sleep(time.Second * 2)
	token, err := SendAuthRequest(rand.Int())
	if err != nil {
		logrus.Errorf("error while sending the token creation request: %v\n", err)
	}

	//here we write the token to the importable variable
	entities.Token = entities.WriteToken(token)
}
