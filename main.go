package main

import (
	"fmt"
	"os"

	"github.com/sosedoff/musicbot/bot"
)

func main() {
	if os.Getenv("MOPIDY_HOST") == "" {
		fmt.Println("MOPIDY_HOST is not provided")
		return
	}

	if os.Getenv("SLACK_TOKEN") == "" {
		fmt.Println("SLACK_TOKEN is not provided")
		return
	}

	if os.Getenv("SLACK_CHANNEL") == "" {
		fmt.Println("SLACK_CHANNEL is not provided")
		return
	}

	bot := bot.NewBot(bot.BotConfig{
		MopidyHost: os.Getenv("MOPIDY_HOST"),
		SlackToken: os.Getenv("SLACK_TOKEN"),
		Channel:    os.Getenv("SLACK_CHANNEL"),
	})

	bot.Run()

	// dummy
	chexit := make(chan bool)
	<-chexit
}
