package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func connectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
func setDataToRedis(rdb *redis.Client, value []byte) {
	err := rdb.Set(ctx, "key", value, 0).Err()
	if err != nil {
		panic(err)
	}

}
func getDataFromRedis(rdb *redis.Client) string {
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	return val

}
