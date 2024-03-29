package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisStorage struct {
	Redis *redis.Client
}

// NewInstance creates a redis instance
func NewInstance(port string, password string, db int) (Storager, error) {
	if port == "" {
		logrus.Errorln("NewInstance func error: port cannot be empty")
		return nil, errors.New("redis connect error: port cannot be empty")
	}
	storage := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis:%s", port),
		Password: password,
		DB:       db,
	})

	//testing whether connection was successful or not
	err := storage.Echo(context.Background(), "~~~Connection to redis established~~~")
	if err.Err() != nil {
		logrus.Errorf("%v", err.Err())
		return &RedisStorage{}, err.Err()
	}

	//logging result. it should display the text given to storage.Echo method
	logrus.Infoln(err.Result())
	return &RedisStorage{Redis: storage}, nil
}

// NewOrder creates an order
func (r *RedisStorage) NewOrder(desc string) (string, error) {
	//creating orderID using uuid package
	orderID := uuid.New().String()[:10]
	//setting order in redis
	err := r.Redis.Set(context.Background(), orderID, desc, 0)
	if err.Err() != nil {
		logrus.Errorf("error while writing order to redis: %v\n", err.Err())
		return "", err.Err()
	}
	logrus.Infoln("New order added. Order id is: %s\n", orderID)
	return orderID, nil
}

// CloseOrder closes the order
func (r *RedisStorage) CloseOrder(id string) error {
	//checking whether we have an order with given id
	if status := r.LookUp(id); !status {
		return errors.New(fmt.Sprintf("can't delete the order %s: such order doesn't exist", id))
	}

	//performing the delete operation in redis
	err := r.Redis.Del(context.Background(), id)
	if err.Err() != nil {
		logrus.Errorf("error while deleting values: %v\n", err.Err())
		return err.Err()
	}
	logrus.Infof("Order %s closed.\n", id)
	return nil
}

// LookUp func is used to check whether redis has given key in it
func (r *RedisStorage) LookUp(id string) bool {
	res := r.Redis.Get(context.Background(), id)
	//if redis error is redis.Nil then we return false i.e. not found
	if res.Err() == redis.Nil {
		return false
	}
	return true
}
