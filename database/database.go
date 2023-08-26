package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var redisClient *redis.Client

func GetRedis() *redis.Client {
	if redisClient == nil {
		fmt.Println("[Databases] Connecting to RedisDB")
		opt, err := redis.ParseURL(os.Getenv("REDIS_URI"))

		if err != nil {
			panic(err)
		}
		redisClient = redis.NewClient(opt)

		// Ping the Redis server to check the connection
		_, err = redisClient.Ping(context.Background()).Result()
		if err != nil {
			panic(err)
		}

		fmt.Println("[Databases] Connected to RedisDB")
	}

	return redisClient
}

func GetMongo() *mongo.Client {
	if mongoClient == nil {
		fmt.Println("[Databases] Connecting to MongoDB")
		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
		clientOptions.SetMaxConnIdleTime(time.Minute * 10)
		clientOptions.SetMaxPoolSize(uint64(100))
		// Connect to MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			panic(err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			panic(err)
		}

		mongoClient = client
		fmt.Println("[Databases] Connected to MongoDB")
	}

	return mongoClient
}
