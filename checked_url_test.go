package dyatl

import (
	"testing"
)

func TestHostAvailable(t *testing.T) {
	for url, ok := range map[string]bool{
		"https://t.me/joinchat/23948092384098230948230": true,
		"https://example.com:1443/sentry/":              true,
		"https://tewyy24rywe08hwqyef.me":                false,
	} {
		if NewCheckedURL(url).HostAvailable() != ok {
			t.Error("fail: available:", url, "!=", ok)
		}
	}
}

func TestLooksCorrect(t *testing.T) {
	for url, ok := range map[string]bool{
		"https://t.me/joinchat/siodfusoiudfoisaudfoiu": true,
		"https://example.com:1443/sentry/":             true,
		"https://tewyy24rywe08hwqyef.me":               true,
		"http://123.ru":                                true,
		"http://123.ru.":                               true,
		"123.ru":                                       false,
		"http://sdfsdf":                                false,
		"http://sdfsdf.sdjfajsdklfjalk":                false,
		"http://sdfsdf.sdjfajsdklfjalk.aero":           true,
		"HTTPS://ya.ru":                                true,
		".com":                                         false,
		"http://.com":                                  false,
	} {
		if NewCheckedURL(url).LooksCorrect() != ok {
			t.Error("fail: looks correct:", url, "!=", ok)
		}
	}
}
