package main

import (
	"twitch_chat_analysis/cmd/helper"
	"twitch_chat_analysis/cmd/model"
)

var (
	rdb  = helper.ConnectToRedis()
	body model.Message
)

func main() {
	helper.SetDataToRedis(rdb)
	helper.GetDataFromRedis(rdb)

}
