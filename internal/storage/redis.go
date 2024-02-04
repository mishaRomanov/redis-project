package storage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisStorage struct {
	Redis *redis.Client
}

// NewInstance creates a redis instance
func NewInstance(port string, password string, db int) Storager {
	storage := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", port),
		Password: password,
		DB:       db})

	//testing whether connection was successful or not
	err := storage.Echo(context.Background(), "~~~Connection to redis established. I guess...~~~")
	if err.Err() != nil {
		logrus.Errorf("%v", err.Err())
		return &RedisStorage{}
	}

	//logging result. it should display the text given to storage.Echo method
	logrus.Infoln(err.Result())
	return &RedisStorage{Redis: storage}
}

// NewOrder creates an order
func (r *RedisStorage) NewOrder(id string, desc string) error {
	err := r.Redis.Set(context.Background(), id, desc, 0)
	if err.Err() != nil {
		logrus.Errorf("error while writing order to redis: %v\n", err.Err())
		return err.Err()
	}
	logrus.Infoln("New order added.")
	return nil
}

// CloseOrder closes the order
func (r *RedisStorage) CloseOrder(id string) error {
	err := r.Redis.Del(context.Background(), id)
	if err.Err() != nil {
		logrus.Errorf("error while deleting values: %v\n", err.Err())
		return err.Err()
	}
	logrus.Infof("Order %s closed.\n", id)
	return nil
}
