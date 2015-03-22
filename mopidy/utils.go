package mopidy

import (
	"github.com/sosedoff/musicbot/spotify"
)

func spotifyTrack(track spotify.Track) Track {
	artists := make([]Artist, 0)
	for _, artist := range track.Artists {
		artists = append(artists, Artist{
			Model: "Artist",
			Name:  artist.Name,
			Uri:   artist.Uri,
		})
	}

	album := Album{
		Model:   "Album",
		Name:    track.Album.Name,
		Uri:     track.Album.Uri,
		Artists: artists,
	}

	return Track{
		Model:   "Track",
		Name:    track.Name,
		Uri:     track.Uri,
		Length:  track.Duration,
		Artists: artists,
		Album:   album,
	}
}
