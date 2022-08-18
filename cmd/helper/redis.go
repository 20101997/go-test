package helper

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"sort"
	"time"
)

var ctx = context.Background()

func ConnectToRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
func SetDataToRedis(rdb *redis.Client, message map[string]string) {
	convertedMessage, _ := json.Marshal(message)
	err := rdb.Set(ctx, message["sender"]+message["receiver"]+message["date"], convertedMessage, 0).Err()
	if err != nil {
		panic(err)
	}

}
func GetDataFromRedis(rdb *redis.Client, key string) map[string]string {
	message, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		panic(message)
	}
	var v map[string]string
	json.Unmarshal(message, &v)
	return v

}

func GetSortedMessagesListFromRedis(rdb *redis.Client, sender string, receiver string) []map[string]string {
	var cursor uint64
	var messagesList []map[string]string

	keys, cursor, err := rdb.Scan(ctx, cursor, sender+receiver+"*", 0).Result()
	keys2, cursor, err := rdb.Scan(ctx, cursor, receiver+sender+"*", 0).Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		messagesList = append(messagesList, GetDataFromRedis(rdb, key))
	}
	for _, key := range keys2 {
		messagesList = append(messagesList, GetDataFromRedis(rdb, key))
	}

	sort.Slice(messagesList, func(i, j int) bool {
		date1, _ := time.Parse("2006-01-02 15:04:05", messagesList[i]["date"][:19])
		date2, _ := time.Parse("2006-01-02 15:04:05", messagesList[j]["date"][:19])
		return date1.After(date2)
	})
	return messagesList

}
