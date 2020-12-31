package dyatl

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func YoutubeIdFromUrl(u *url.URL) string {
	switch u.Host {
	case "www.youtube.com", "m.youtube.com", "youtube.com":
		return u.Query().Get("v")
	case "youtu.be":
		return strings.Trim(u.Path, "/")
	}
	return ""
}

type YoutubeData struct {
	Type    string `json:"type"`
	Version string `json:"version"`

	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
	AuthorUrl  string `json:"author_url"`

	Html   string `json:"html"`
	Width  int    `json:"width"`
	Height int    `json:"height"`

	ProviderName string `json:"provider_name"`
	ProviderUrl  string `json:"provider_url"`

	ThumbnailUrl    string `json:"thumbnail_url"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
}

var (
	Unauthorized = errors.New("Unauthorized")
	BadRequest   = errors.New("Bad Request")
)

func (c Client) GetYoutubeData(id string) (YoutubeData, error) {
	var data YoutubeData

	dataUrl := "https://www.youtube.com/oembed?url=http%3A//youtube.com/watch%3Fv%3D" + id
	resp, err := c.Get(dataUrl)
	if err != nil {
		return data, err
	}

	dataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, errors.Wrapf(err, "read body fail on `%s`", dataUrl)
	}
	defer resp.Body.Close()

	switch string(dataBytes) {
	case "Bad Request":
		return data, BadRequest
	case "Unauthorized":
		return data, Unauthorized
	default:
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return data, errors.Wrapf(err, "unmarshal fail on `%s`", string(dataBytes))
		}
	}

	return data, err
}
