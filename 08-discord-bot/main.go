package main

import (
	"log"

	"github.com/amitamrutiya2210/08-discord-bot/bot"
	"github.com/joho/godotenv"
)

func main() {
	erro := godotenv.Load()

	if erro != nil {
		log.Fatal("Error loading .env file")
	}

	bot.Start()

	<-make(chan struct{})
	return

}
