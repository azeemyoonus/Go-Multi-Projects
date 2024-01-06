package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Event")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	errr := godotenv.Load()
	if errr != nil {
		log.Fatal("Error loading .env file")
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	go printCommandEvents(bot.CommandEvents())
	bot.Command("ping", &slacker.CommandDefinition{
		Description: "Ping!",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("Pong!")
		},
	})

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "YOB calculator",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				response.Reply("Invalid year")
				return
			}
			age := 2024 - yob
			r := fmt.Sprintf("Your age is %d", age)
			response.Reply(r)
		},
	})

	bot.Command("my name is <name>", &slacker.CommandDefinition{
		Description: "Name",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			name := request.Param("name")
			r := fmt.Sprintf("Your name is %s", name)
			response.Reply(r)
		},
	})
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

}
