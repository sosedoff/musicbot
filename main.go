package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/jessevdk/go-flags"
	"github.com/sosedoff/musicbot/bot"
)

var options struct {
	MopidyHost   string `long:"mopidy" description:"Mopidy server host:port" env:"MOPIDY_HOST"`
	SlackToken   string `long:"slack-token" description:"Slack integration token" env:"SLACK_TOKEN"`
	SlackChannel string `long:"slack-channel" description:"Slack channel name" default:"general" env:"SLACK_CHANNEL"`
	Debug        bool   `short:"d" long:"debug" description:"Enable debugging mode" default:"false"`
}

func init() {
	_, err := flags.ParseArgs(&options, os.Args)

	if err != nil {
		os.Exit(1)
	}

	if options.MopidyHost == "" {
		fmt.Println("Error: Mopidy host is not provided")
		os.Exit(1)
	}

	if options.SlackToken == "" {
		fmt.Println("Error: Slack token is not provided")
		os.Exit(1)
	}
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func main() {
	bot := bot.NewBot(bot.BotConfig{
		MopidyHost: options.MopidyHost,
		SlackToken: options.SlackToken,
		Channel:    options.SlackChannel,
	})

	bot.Run()
	handleSignals()
}
