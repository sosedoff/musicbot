package main

import (
	"os"
)

func main() {
	bot := NewBot(BotConfig{
		MopidyHost: os.Getenv("MOPIDY_HOST"),
		SlackToken: os.Getenv("SLACK_TOKEN"),
		Channel:    os.Getenv("SLACK_CHANNEL"),
	})

	bot.Run()

	// dummy
	chexit := make(chan bool)
	<-chexit
}
