package dyatl

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
)

type Client struct {
	http.Client
	UserAgent      string
	AcceptLanguage string
}

func NewClient() Client {
	c := Client{
		UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
		AcceptLanguage: "ru-RU, ru;q=0.9, en-US;q=0.8, en;q=0.7",
	}
	c.Timeout = 5 * time.Second
	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSHandshakeTimeout: c.Timeout,
	}
	return c
}

var (
	htmlRegex     = regexp.MustCompile(`(?i)\s*<!DOCTYPE`)
	redirectRegex = regexp.MustCompile(`<script[^>]*>\s*window.location.href = "([^"]+)"\s*</script>`)
)

func (c Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Accept-Language", c.AcceptLanguage)
	return c.Do(req)
}

var (
	NotHttpUrl      = errors.New("not http url")
	PreviewNotFound = errors.New("preview not found")
)

type Preview struct {
	Title        string
	ThumbnailUrl string
	YoutubeId    string
}

func (c Client) Preview(link string) (Preview, error) {
	u, err := url.Parse(link)
	if err != nil {
		return Preview{}, err
	}

	if !(u.Scheme == "http" || u.Scheme == "https") {
		return Preview{}, NotHttpUrl
	}

	if id := YoutubeIdFromUrl(u); id != "" {
		data, _ := c.GetYoutubeData(id)
		if data.ThumbnailUrl != "" {
			p := Preview{
				Title:        data.Title,
				ThumbnailUrl: data.ThumbnailUrl,
				YoutubeId:    id,
			}
			return p, nil
		}
	}

	return c.previewFromContent(u)
}

func (c Client) previewFromContent(u *url.URL) (Preview, error) {
	resp, err := c.Get(u.String())
	if err != nil {
		return Preview{}, err
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	switch {
	case ct == "":
		buf := make([]byte, 20*1024)
		if _, err := resp.Body.Read(buf); err != nil && err != io.EOF {
			return Preview{}, err
		}
		if htmlRegex.Match(buf) {
			return c.fromGoquery(u, bytes.NewReader(buf))
		}
	case strings.HasPrefix(ct, "text/html"):
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return Preview{}, nil
		}
		bits := strings.Split(strings.ToLower(ct), "charset=")
		if len(bits) == 2 {
			switch bits[1] { // FIXME: all codepages
			case "windows-1251":
				dec := charmap.Windows1251.NewDecoder()
				body, err = dec.Bytes(body)
				if err != nil {
					return Preview{}, errors.Wrap(err, "decode cp1251 fail")
				}
			}
		}
		if m := redirectRegex.FindSubmatch(body); len(m) > 0 {
			if newLink := string(m[1]); newLink != u.String() {
				return c.Preview(newLink)
			}
		}
		return c.fromGoquery(u, bytes.NewReader(body))
	case strings.HasPrefix(ct, "image/"):
		return Preview{ThumbnailUrl: u.String()}, nil
	}
	return Preview{}, PreviewNotFound
}

func (c Client) fromGoquery(base *url.URL, r io.Reader) (Preview, error) {
	q, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return Preview{}, errors.Wrap(err, "NewDocumentFromReader fail")
	}

	title := strings.Join(strings.Fields(q.Find("title").Text()), " ")
	p := Preview{Title: title}

	q.Find("meta").Each(func(i int, selection *goquery.Selection) {
		switch selection.AttrOr("property", "") {
		case "og:title":
			if v := selection.AttrOr("content", ""); v != "" {
				p.Title = v
			}
		case "og:image":
			if v := selection.AttrOr("content", ""); v != "" {
				switch {
				case strings.HasPrefix(v, "//"):
					p.ThumbnailUrl = "http:" + v
				case strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://"):
					p.ThumbnailUrl = v
				case strings.HasPrefix(v, "/"):
					p.ThumbnailUrl = base.Scheme + "://" + base.Host + v
				}
			}
		}
	})
	return p, nil
}

func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}
	if strings.Contains(err.Error(), "Client.Timeout exceeded") {
		return true
	}
	return false
}
