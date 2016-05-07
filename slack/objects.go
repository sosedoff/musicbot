package slack

type Team struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RtmResponse struct {
	Ok       bool       `json:"ok"`
	Url      string     `json:"url"`
	Team     Team       `json:"team"`
	Users    []*User    `json:"users"`
	Channels []*Channel `json:"channels"`
	Groups   []*Channel `json:"groups"`
}

type Message struct {
	User      *User    `json:"-"`
	Channel   *Channel `json:"-"`
	ChannelId string   `json:"channel"`
	Text      string   `json:"text"`
	Timestamp string   `json:"ts"`
}
