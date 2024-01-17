package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-6461516691651-6447130371815-DXwgQflBQx1bBMu1UGTw90ZV")
	os.Setenv("SLACK_APP_TOKEN", "xapp-1-A06DY89DND7-6464103379220-4b2db0b58757f65e8e82d14bfc68220c3cf9da28e116116ec6f16b461b0269b3")
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "yob calculator",
		Examples:    []string{"my yob is 2020"},
		Handler: func(botCtx slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			year := r.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println("Couldn't parse request:", err)
			}

			currYear := time.Now().Year()
			age := currYear - yob

			result := fmt.Sprintf("age is %v", age)
			w.Reply(result)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
}
