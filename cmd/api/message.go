package main

type message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}
