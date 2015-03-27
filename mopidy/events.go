package mopidy

type Event struct {
	Type string `json:"event"`
	Data interface{}
}

type VolumeChangedEvent struct {
	Volume int `json:"volume"`
}

type TrackPlaybackStarted struct {
	TlTrack `json:"tl_track"`
}
