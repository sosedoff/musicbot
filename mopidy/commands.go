package mopidy

import (
	"github.com/sosedoff/musicbot/spotify"
)

func (m *Client) ClearTracklist() error {
	_, err := m.Call("core.tracklist.clear", nil)
	return err
}

func (m *Client) Tracks() ([]Track, error) {
	data, err := m.Call("core.tracklist.get_tracks", nil)
	if err != nil {
		return []Track{}, err
	}

	resp := data.(*TracksResponse)
	return resp.Result, nil
}

func (m *Client) AddTracks(tracks []Track) error {
	params := map[string][]Track{"tracks": tracks}
	_, err := m.Call("core.tracklist.add", params)
	return err
}

func (m *Client) AddSpotifyTracks(tracks []spotify.Track) error {
	newTracks := make([]Track, 0)
	for _, track := range tracks {
		newTracks = append(newTracks, spotifyTrack(track))
	}
	return m.AddTracks(newTracks)
}

func (m *Client) State() (string, error) {
	data, err := m.Call("core.playback.get_state", nil)
	if err != nil {
		return "", err
	}

	resp := data.(*BasicResponse)
	return resp.Result.(string), nil
}

func (mopidy *Client) CurrentTrack() (*Track, error) {
	data, err := mopidy.Call("core.playback.get_current_track", nil)
	if err != nil || data == nil {
		return nil, err
	}

	resp := data.(*TrackResponse)
	return resp.Result, nil
}

func (m *Client) PlayNextTrack() error {
	_, err := m.Call("core.playback.next", nil)
	return err
}

func (m *Client) Pause() error {
	_, err := m.Call("core.playback.pause", nil)
	return err
}

func (m *Client) Resume() error {
	_, err := m.Call("core.playback.resume", nil)
	return err
}

func (m *Client) Play() error {
	_, err := m.Call("core.playback.play", nil)
	return err
}

func (m *Client) Stop(clearTrack bool) error {
	_, err := m.Call("core.playback.stop", nil)
	return err
}

func (m *Client) Volume() (int, error) {
	data, err := m.Call("core.playback.get_volume", nil)
	if err != nil {
		return 0, err
	}

	resp := data.(*BasicResponse)
	return int(resp.Result.(float64)), nil
}

func (m *Client) SetVolume(percent int) error {
	params := map[string]int{"volume": percent}
	_, err := m.Call("core.playback.set_volume", params)
	return err
}

func (m *Client) Search(query string) {
	params := map[string][]string{
		"any":  []string{query},
		"uris": []string{"spotify:"},
	}

	m.Call("core.library.search", params)
}
