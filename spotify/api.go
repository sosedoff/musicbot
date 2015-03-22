package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func makeUrl(path string) string {
	return "https://api.spotify.com/v1" + path
}

func makeSearchRequest(options SearchOptions) (*http.Request, error) {
	req, err := http.NewRequest("GET", makeUrl("/search"), nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("query", options.Query)
	params.Add("type", options.Type)

	if options.Market != "" {
		params.Add("market", options.Market)
	}

	if options.Offset > 0 {
		params.Add("offset", strconv.Itoa(options.Offset))
	}

	if options.Limit > 0 {
		params.Add("limit", strconv.Itoa(options.Limit))
	}

	req.URL.RawQuery = params.Encode()

	return req, nil
}

func GetAlbum(id string) (*Album, error) {
	url := fmt.Sprintf("%s/%s", makeUrl("/albums"), id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var album Album

	err = json.Unmarshal(body, &album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}

func Search(options SearchOptions) (*SearchResult, error) {
	req, err := makeSearchRequest(options)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", body)

	result := SearchResult{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
