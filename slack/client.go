package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	token      string
	accepted   bool
	ws         *websocket.Conn
	users      map[string]*User
	channels   map[string]*Channel
	channelIds map[string]string
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		accepted:   false,
		users:      map[string]*User{},
		channels:   map[string]*Channel{},
		channelIds: map[string]string{},
	}
}

func (s *Client) Connect() error {
	log.Println("Connecting to Slack API...")
	url, err := s.getSocketUrl()
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Println("Connecting to Slack Websocket")
	ws, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	s.ws = ws

	return nil
}

func (s *Client) FindUser(name string) *User {
	return s.users[name]
}

func (s *Client) FindChannel(name string) *Channel {
	return s.channels[name]
}

func (s *Client) Close() {
	if s.ws != nil {
		s.ws.Close()
	}
}

func (s *Client) SendMessage(channel string, text string) error {
	msg := map[string]string{
		"type":    "message",
		"channel": s.channelIds[channel],
		"text":    text,
	}

	return s.ws.WriteJSON(msg)
}

func (s *Client) Run(receiver chan Event) {
	for {
		event := json.RawMessage{}

		err := s.ws.ReadJSON(&event)
		if err != nil {
			fmt.Println("Websocket error:", err)
			time.Sleep(time.Second * 3)
			if s.Connect() != nil {
				return
			}
		}

		s.handleEvent(receiver, event)
	}
}

func parseMessage(data json.RawMessage) (Message, error) {
	msg := Message{}

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (s *Client) handleMessageEvent(receiver chan Event, data json.RawMessage) {
	msg, err := parseMessage(data)
	if err != nil {
		fmt.Println("Err:", err)
		return
	}

	if msg.Text == "" {
		return
	}

	event := Event{
		Data: MessageEvent{
			Message: msg,
		},
	}

	receiver <- event
}

func (s *Client) handleEvent(receiver chan Event, data json.RawMessage) {
	event := Event{}

	err := json.Unmarshal(data, &event)
	if err != nil {
		fmt.Println("JSON err:", err)
		return
	}

	switch event.Type {
	case "hello":
		receiver <- Event{Data: HelloEvent{}}
	case "message":
		s.handleMessageEvent(receiver, data)
	}
}

func (s *Client) getSocketUrl() (string, error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", s.token)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var rtm RtmResponse
	err = json.Unmarshal(buff, &rtm)
	if err != nil {
		return "", err
	}

	for _, user := range rtm.Users {
		s.users[user.Id] = user
	}

	for _, channel := range rtm.Channels {
		s.channels[channel.Id] = channel
		s.channelIds[channel.Name] = channel.Id
	}

	return rtm.Url, nil
}
