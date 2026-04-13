package redis

import (
	"context"
	"fmt"
	"jvalleyverse/pkg/config"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var IsConnected bool = false

func ConnectRedis() {
    cfg := config.AppConfig
    Client = redis.NewClient(&redis.Options{
        Addr:     cfg.RedisHost,
        Password: cfg.RedisPass,
        DB:       cfg.RedisDB,
    })

    // test connection
    if err := Client.Ping(context.Background()).Err(); err != nil {
        log.Printf("⚠️  Redis connection failed (caching disabled): %v\n", err)
        IsConnected = false
        Client = nil
        return
    }
    
    IsConnected = true
    fmt.Println("✅ Redis connected successfully")
}

// IsAvailable checks if Redis is available
func IsAvailable() bool {
    return IsConnected && Client != nil
}