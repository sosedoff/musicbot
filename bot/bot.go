package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/sosedoff/musicbot/mopidy"
	"github.com/sosedoff/musicbot/slack"
)

type BotConfig struct {
	SlackChannel string
	SlackToken   string
	MopidyHost   string
	Debug        bool
}

type Bot struct {
	config      BotConfig
	slack       *slack.Client
	mopidy      *mopidy.Client
	slackEvent  chan slack.Event
	mopidyEvent chan mopidy.Event
	commands    []Command
}

func NewBot(config BotConfig) Bot {
	return Bot{config: config}
}

func (bot *Bot) initSlack() error {
	client := slack.NewClient(bot.config.SlackToken)

	err := client.Connect()
	if err != nil {
		return err
	}

	bot.slack = client
	bot.slackEvent = make(chan slack.Event)

	return nil
}

func (bot *Bot) initMopidy() error {
	client := mopidy.New(bot.config.MopidyHost)

	err := client.Connect()
	if err != nil {
		return err
	}

	bot.mopidy = client
	bot.mopidyEvent = make(chan mopidy.Event)

	return nil
}

func (bot *Bot) handleSlackEvents() {
	for {
		ev := <-bot.slackEvent
		switch ev.Data.(type) {
		case slack.MessageEvent:
			msg := ev.Data.(slack.MessageEvent)
			bot.handleMessage(&msg.Message)
		}
	}
}

func (bot *Bot) handleMopidyEvents() {
	for {
		ev := <-bot.mopidyEvent

		switch ev.Data.(type) {
		case mopidy.TrackPlaybackStarted:
			data := ev.Data.(mopidy.TrackPlaybackStarted)
			track := data.TlTrack.Track
			link := fmt.Sprintf("http://open.spotify.com/track/%s", strings.Split(track.Uri, ":")[2])

			log.Printf("Playing track: %s - %s", track.Uri, track.String())
			bot.Say(fmt.Sprintf(":musical_note: %s - %s", track.String(), link))

		case mopidy.VolumeChangedEvent:
			data := ev.Data.(mopidy.VolumeChangedEvent)
			bot.Say(fmt.Sprintf("Volume is set to %v%s", data.Volume, "%"))
		}
	}
}

func (bot *Bot) handleMessage(msg *slack.Message) {
	if (msg.Channel != nil && msg.Channel.Name != bot.config.SlackChannel) {
		fmt.Println("ignoreing message from other chan")
		return
	}
	for _, cmd := range bot.commands {
		match := cmd.Match(msg.Text)

		if match != nil {
			go cmd.handler(bot, match)
			break
		}
	}
}

func (bot *Bot) addCommand(expr string, handler HandlerFunc) {
	cmd := NewCommand(bot, expr, handler)
	bot.commands = append(bot.commands, cmd)
}

func (bot *Bot) Say(text string) {
	err := bot.slack.SendMessage(bot.config.SlackChannel, text)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (bot *Bot) Run() {
	if err := bot.initMopidy(); err != nil {
		log.Fatalln(err)
	}

	if err := bot.initSlack(); err != nil {
		log.Fatalln(err)
	}

	setupCommands(bot)

	go bot.handleSlackEvents()
	go bot.handleMopidyEvents()

	go bot.mopidy.Run(bot.mopidyEvent)
	go bot.slack.Run(bot.slackEvent)
}
