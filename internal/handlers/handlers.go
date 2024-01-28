package handlers

import (
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	redis *redis.Client
}
