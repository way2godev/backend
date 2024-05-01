package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitCache() {
	host:= os.Getenv("REDIS_HOST")
	port:= os.Getenv("REDIS_PORT")
	user:= os.Getenv("REDIS_USER")
	password:= os.Getenv("REDIS_PASSWORD")
	if user == "" || password == "" || host == "" || port == "" {
		log.Fatalf("Failed to establish cache connection: Credentials not provided")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Username: user,
		Password: password,
		DB:       0,
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatalf("Failed to establish cache connection: %v", err)
	}

	log.Print("Cache connection established")
}

func GetRedisClient() *redis.Client {
	return redisClient
}