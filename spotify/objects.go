package spotify

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type Album struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Uri    string  `json:"uri"`
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Uri      string   `json:"uri"`
	Artists  []Artist `json:"artists"`
	Album    Album    `json:"album"`
	Duration int      `json:"duration_ms"`
}

type SearchOptions struct {
	Query  string
	Type   string
	Market string
	Limit  int
	Offset int
}

type SearchResult struct {
	Tracks struct {
		Items []Track `json:"items"`
		Total int     `json:"total"`
		Limit int     `json:"limit"`
	}
	Artists struct {
		Items []Artist `json:"items"`
		Total int      `json:"total"`
		Limit int      `json:"limit"`
	}
	Albums struct {
		Items []Album `json:"items"`
		Total int     `json:"total"`
		Limit int     `json:"limit"`
	}
}
