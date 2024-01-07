package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Krognol/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
)

var wolframClient *wolfram.Client

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
	client := witai.NewClient(os.Getenv("WIT_SERVER_ACCESS_TOKEN"))

	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}
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

	bot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			q := request.Param("message")
			r := fmt.Sprintf("You asked %s", q)
			response.Reply(r)
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: q,
			})
			data, _ := json.MarshalIndent(msg, "", "  ")
			value := gjson.Get(string(data), "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			answer := value.String()
			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 100)
			if err != nil {
				fmt.Println(err)
			}
			response.Reply(res)
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
