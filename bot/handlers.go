package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sosedoff/musicbot/spotify"
)

var usage = `
#help# - Print commands
#current# - Show current track
#next|skip# - Play next track
#pause|resume|stop# - Control playback
#play <query># - Play 10 tracks that match query
#tracks|list# - Show tracks in queue
#clear# - Remove all tracks and stop playback
#state# - Get player state
#vol|volume# - Get current volume level
#vol|volume up|down|0-100# - Set volume level
`

func setupCommands(bot *Bot) {
	bot.addCommand("^help$", Help)
	bot.addCommand("^current$", CurrentTrack)
	bot.addCommand("^(next|skip)$", NextTrack)
	bot.addCommand("^pause$", Pause)
	bot.addCommand("^resume$", Resume)
	bot.addCommand("^stop$", Stop)
	bot.addCommand("^play$", Resume)
	bot.addCommand("^play (.*)", Play)
	bot.addCommand("^(tracks|list)$", Tracks)
	bot.addCommand("^clear$", Clear)
	bot.addCommand("^state$", State)
	bot.addCommand("^(vol|volume)$", Volume)
	bot.addCommand("^(vol|volume) (up|down|[0-9]+)$", SetVolume)
}

func Help(bot *Bot, match *Match) {
	bot.Say(strings.TrimSpace(strings.Replace(usage, "#", "`", -1)))
}

func CurrentTrack(bot *Bot, match *Match) {
	track, err := bot.mopidy.CurrentTrack()
	if err != nil {
		bot.Say("Cant get current track")
		return
	}

	if track == nil {
		bot.Say("No current track")
		return
	}

	bot.Say(track.String())
}

func NextTrack(bot *Bot, match *Match) {
	err := bot.mopidy.PlayNextTrack()
	if err != nil {
		bot.Say("Cant play next track")
		return
	}
}

func Pause(bot *Bot, match *Match) {
	err := bot.mopidy.Pause()
	if err != nil {
		bot.Say("Cant pause")
		return
	}
}

func Resume(bot *Bot, match *Match) {
	err := bot.mopidy.Resume()
	if err != nil {
		bot.Say("Cant resume")
		return
	}
}

func Stop(bot *Bot, match *Match) {
	err := bot.mopidy.Stop(true)
	if err != nil {
		bot.Say("Cant stop, wont stop!")
		return
	}
}

func Tracks(bot *Bot, match *Match) {
	current := ""

	track, _ := bot.mopidy.CurrentTrack()
	if track != nil {
		current = track.Uri
	}

	tracks, err := bot.mopidy.Tracks()
	if err != nil {
		bot.Say("Cant get tracks")
		return
	}

	if len(tracks) == 0 {
		bot.Say("Queue is empty")
		return
	}

	lines := make([]string, len(tracks))

	// For now just print 10 tracks. Slack will cut WS if you try to send a lot of data.
	num := len(tracks)
	if num > 10 {
		num = 10
	}

	for i, track := range tracks[0:num] {
		if track.Uri == current {
			lines[i] = fmt.Sprintf("*%d. %s*", i+1, track.String())
		} else {
			lines[i] = fmt.Sprintf("%d. %s", i+1, track.String())
		}
	}

	bot.Say(strings.Join(lines, "\n"))
}

func Clear(bot *Bot, match *Match) {
	err := bot.mopidy.ClearTracklist()
	if err != nil {
		bot.Say("Cant clear queue")
		return
	}

	bot.Say("Queue is cleared")
}

func State(bot *Bot, match *Match) {
	state, err := bot.mopidy.State()
	if err != nil {
		bot.Say("Cant get state")
		return
	}

	bot.Say("Player state: " + state)
}

func Play(bot *Bot, match *Match) {
	query := match.Values[0]

	opts := spotify.SearchOptions{
		Query: query,
		Type:  "track",
		Limit: 10,
	}

	result, err := spotify.Search(opts)
	if err != nil {
		fmt.Println(err)
		bot.Say("_Spotify search failed_")
		return
	}

	if len(result.Tracks.Items) == 0 {
		bot.Say("Nothing found for: " + query)
		return
	}

	// If player is stopped we should clear old track list so that playback will s
	// start with only new tracks. This is needed to keep the track list small.
	state, _ := bot.mopidy.State()
	if state == "stopped" {
		bot.mopidy.ClearTracklist()
	}

	err = bot.mopidy.AddSpotifyTracks(result.Tracks.Items)
	if err != nil {
		bot.Say("Cant add tracks to the queue")
		return
	}

	// Start playback only if player is stopped.
	state, _ = bot.mopidy.State()
	if state == "stopped" {
		bot.mopidy.Play()
	}

	// Build a string that only includes 10 tracks. Its a dirty hack to make sure
	// that amount of data sent to slack stays low, otherwise slack will terminate
	// websocket connection. TODO: need a better way of handing this.
	lines := make([]string, len(result.Tracks.Items))
	for i, track := range result.Tracks.Items {
		lines[i] = fmt.Sprintf("%v. %s - %s", i+1, track.Name, track.Album.Name)
	}

	bot.Say("Added tracks:\n" + strings.Join(lines, "\n"))
}

func Volume(bot *Bot, match *Match) {
	vol, err := bot.mopidy.Volume()
	if err != nil {
		bot.Say("Cant get volume")
		return
	}

	bot.Say(fmt.Sprintf("Current volume: %v%s", vol, "%"))
}

func SetVolume(bot *Bot, match *Match) {
	vol, err := bot.mopidy.Volume()
	if err != nil {
		bot.Say("Cant get volume")
		return
	}

	newvol := match.Values[1]

	switch newvol {
	case "up":
		vol += 10
	case "down":
		vol -= 10
	default:
		vol, err = strconv.Atoi(newvol)
		if err != nil {
			bot.Say("Invalid volume value")
			return
		}
	}

	if vol > 110 || vol < 0 {
		bot.Say("Volume range is 0-110")
		return
	}

	err = bot.mopidy.SetVolume(vol)
	if err != nil {
		bot.Say("Cant change volume")
		return
	}
}
