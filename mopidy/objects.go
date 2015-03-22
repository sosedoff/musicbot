package mopidy

import (
	"fmt"
)

type Artist struct {
	Model string `json:"__model__"`
	Name  string `json:"name"`
	Uri   string `json:"uri"`
}

type Album struct {
	Model   string   `json:"__model__"`
	Name    string   `json:"name"`
	Uri     string   `json:"uri"`
	Artists []Artist `json:"artists"`
}

type Track struct {
	Model   string   `json:"__model__"`
	Name    string   `json:"name"`
	Uri     string   `json:"uri"`
	Length  int      `json:"length"`
	TrackNo int      `json:"track_no"`
	Album   Album    `json:"album"`
	Artists []Artist `json:"artists"`
}

type TlTrack struct {
	Model string `json:"__model__"`
	TlId  int    `json:"tlid"`
	Track Track  `json:"track"`
}

type SearchResult struct {
	Artists []Artist `json:"artists"`
	Albums  []Artist `json:"albums"`
	Tracks  []Track  `json:"tracks"`
}

func (t Track) String() string {
	return fmt.Sprintf("%s - %s - %s", t.Name, t.Album.Name, t.Artists[0].Name)
}
