package mopidy

import (
	"encoding/json"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Type      string `json:"type"`
		Message   string `json:"message"`
		Traceback string `json:"traceback"`
	} `json:"data"`
}

type Response struct {
	Error  *json.RawMessage `json:"error"`
	Result *json.RawMessage `json:"result"`
}

type BasicResponse struct {
	Error  `json:"error"`
	Result interface{} `json:"result"`
}

type TrackResponse struct {
	Error  `json:"error"`
	Result *Track `json:"result"`
}

type TracksResponse struct {
	Error  `json:"error"`
	Result []Track `json:"result"`
}
