package main

import (
	"log"
	"main/pkg/bot"
	"main/pkg/session"
	"os"
)

func main() {
	token, err := os.ReadFile("token.txt")
	if err != nil {
		log.Panic(err)
	}

	myBot, err := bot.NewBot(string(token))
	if err != nil {
		log.Panic(err)
	}

	userSessions := make(map[int64]*session.UserSession)
	myBot.StartListening(userSessions)
}
