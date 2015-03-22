package mopidy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	rpc "github.com/gorilla/rpc/v2/json2"
	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	url  string
	rpc  string
	conn *websocket.Conn
}

func New(host string) *Client {
	return &Client{
		url: fmt.Sprintf("ws://%s/mopidy/ws", host),
		rpc: fmt.Sprintf("http://%s/mopidy/rpc", host),
	}
}

func (mopidy *Client) Connect() error {
	ws, _, err := dialer.Dial(mopidy.url, nil)
	if err != nil {
		return err
	}

	mopidy.conn = ws
	return nil
}

func reponseObjectForCommand(command string) interface{} {
	mapping := map[string]interface{}{
		"core.playback.get_current_track": &TrackResponse{},
		"core.tracklist.get_tracks":       &TracksResponse{},
	}

	obj := mapping[command]
	if obj == nil {
		obj = &BasicResponse{}
	}

	return obj
}

func (m *Client) Call(command string, params interface{}) (interface{}, error) {
	if params == nil {
		params = map[string]string{}
	}

	buff, err := rpc.EncodeClientRequest(command, params)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buff)
	resp, err := http.Post(m.rpc, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("***\n%s\n***\n", body)

	var r Response
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&r); err != nil {
		return nil, err
	}

	if r.Error != nil {
		jsonErr := Error{}
		err = json.Unmarshal(*r.Error, &jsonErr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(jsonErr.Message)
	}

	if r.Result == nil {
		return nil, nil
	}

	obj := reponseObjectForCommand(command)
	err = json.Unmarshal(body, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (mopidy *Client) handleEvent(receiver chan Event, data json.RawMessage) {
	event := Event{}

	err := json.Unmarshal(data, &event)
	if err != nil {
		fmt.Println("JSON err:", err)
		return
	}

	switch event.Type {
	case "track_playback_started":
		ev := TrackPlaybackStarted{}
		json.Unmarshal(data, &ev)
		receiver <- Event{Data: ev}

	case "volume_changed":
		msg := VolumeChangedEvent{}
		json.Unmarshal(data, &msg)
		receiver <- Event{Data: msg}
	}
}

func (mopidy *Client) Run(receiver chan Event) {
	for {
		event := json.RawMessage{}

		err := mopidy.conn.ReadJSON(&event)
		if err != nil {
			fmt.Println("WS Error:", err)
			return
		}

		fmt.Printf("%s\n", event)
		mopidy.handleEvent(receiver, event)
	}
}
