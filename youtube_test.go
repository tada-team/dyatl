package dyatl

import (
	"net/url"
	"strings"
	"testing"
)

func TestYoutubeIdFromUrl(t *testing.T) {
	for link, id := range map[string]string{
		"https://www.youtube.com/watch?v=iJHV-96Z6ZY": "iJHV-96Z6ZY",
		"https://youtu.be/GCYZMjFczmM?t=7":            "GCYZMjFczmM",
		"https://xxxxx.com":                           "",
	} {
		t.Run(link, func(t *testing.T) {
			u, err := url.Parse(link)
			if err != nil {
				t.Fatal(err)
			}
			v := YoutubeIdFromUrl(u)
			if v != id {
				t.Error("invalid id, want:", id, "got:", v)
			}
		})
	}
}

func TestGetYoutubeData(t *testing.T) {
	c := NewClient()
	t.Run("private video", func(t *testing.T) {
		u, err := url.Parse("https://www.youtube.com/watch?v=Ud3p8p9hqHc")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.GetYoutubeData(YoutubeIdFromUrl(u)); err != Unauthorized {
			t.Fatal("must be Unauthorized, got:", err)
		}
	})

	t.Run("unknown video", func(t *testing.T) {
		u, err := url.Parse("https://www.youtube.com/watch?v=jfoiwufoisudofi")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.GetYoutubeData(YoutubeIdFromUrl(u)); err != BadRequest {
			t.Fatal("must be BadRequest, got:", err)
		}
	})

	t.Run("valid video", func(t *testing.T) {
		u, err := url.Parse("https://www.youtube.com/watch?v=zFLp-lJfA64")
		if err != nil {
			t.Fatal(err)
		}
		data, err := c.GetYoutubeData(YoutubeIdFromUrl(u))
		if err != nil{
			t.Fatal(err)
		}
		if !strings.HasPrefix(data.Title, "А.Лаэртский") {
			t.Fatal("invalid title:", data.Title)
		}
	})
}
