package dyatl

import (
	"strings"
	"testing"
)

func TestPreview(t *testing.T) {
	c := NewClient()
	for link, testCase := range map[string]Preview{
		"https://www.youtube.com/watch?v=eLAHSRmFFzE": {
			ThumbnailUrl: "https://i.ytimg.com/vi/eLAHSRmFFzE/hqdefault.jpg",
			Title:        "Noize MC",
			YoutubeId:    "eLAHSRmFFzE",
		},
		"https://www.youtube.com/watch?v=Ud3p8p9hqHc": {
			ThumbnailUrl: "https://i.ytimg.com/vi/Ud3p8p9hqHc/hqdefault.jpg",
			Title:        "Коронавирус в Китае:",
			YoutubeId:    "", // empty for private video because can't be embedded
		},
		"http://example.com.": {
			Title: "Example Domain",
		},
		"https://vk.com/shantitrash": {
			// skip: has different servers for thumbnails
			// ThumbnailUrl: "http://sun1-89.userapi.com/impf/c631625/v631625378/401/bXrjmBHZwP8.jpg?size=200x0&quality=96&crop=59,23,527,528&sign=b6e64806e0947daa745c9ed26df5c427&ava=1",
			Title:        "Шанти-трэш",
		},
		"https://habr.com/ru/post/354472/": {
			ThumbnailUrl: "https://habr.com/share/publication/354472/590e252e0c6748a2c367a8c72a67c422/?v=1",
			Title:        "Бэкап переписки в telegram",
		},
		"https://yandex.ru/maps/213/moscow/?ll=37.656245%2C55.759324&mode=search&oid=1111668074&ol=biz&sctx=ZAAAAAgCEAAaKAoSCYnQCDauz0JAEabtX1lp4EtAEhIJA137Anqh8j8RmPxP%2Fu4d7D8iBQABAgQFKAAwATi47Z3olunF%20cQBQNUBSAFVzczMPlgAYiRtaWRkbGVfYXNrX2RpcmVjdF9xdWVyeV90eXBlcz1ydWJyaWNiKG1pZGRsZV9pbmZsYXRlX2RpcmVjdF9maWx0ZXJfd2luZG93PTUwMDBiEnJlbGV2X2RydWdfYm9vc3Q9MWJEbWlkZGxlX2RpcmVjdF9zbmlwcGV0cz1waG90b3MvMi54LGJ1c2luZXNzcmF0aW5nLzIueCxtYXNzdHJhbnNpdC8xLnhiNW1pZGRsZV93aXpleHRyYT10cmF2ZWxfY2xhc3NpZmllcl92YWx1ZT0wLjAxOTcyODM2MjU2YidtaWRkbGVfd2l6ZXh0cmE9YXBwbHlfZmVhdHVyZV9maWx0ZXJzPTFiKG1pZGRsZV93aXpleHRyYT1vcmdtbl93YW5kX3RocmVzaG9sZD0wLjliKW1pZGRsZV93aXpleHRyYT1yZXF1ZXN0X3NvZnRfdGltZW91dD0wLjA1YiNtaWRkbGVfd2l6ZXh0cmE9dHJhbnNpdF9hbGxvd19nZW89MWI9bWlkZGxlX3dpemV4dHJhPXRyYXZlbF9jbGFzc2lmaWVyX29yZ21hbnlfdmFsdWU9MC4wMTY0ODYxNDkyOGIebWlkZGxlX2Fza19kaXJlY3RfcGVybWFsaW5rcz0xYiptaWRkbGVfaW5mbGF0ZV9kaXJlY3RfcmVxdWVzdF93aW5kb3c9MTAwMDBiHXJlbGV2X2ZpbHRlcl9nd2tpbmRzPTAuMywwLjQ1YilyZWFycj1zY2hlbWVfTG9jYWwvR2VvL0FsbG93VHJhdmVsQm9vc3Q9MWIxcmVhcnI9c2NoZW1lX0xvY2FsL0dlb3VwcGVyL2ZlYXR1cmVzRnJvbU9iamVjdHM9MWIvcmVhcnI9c2NoZW1lX0xvY2FsL0dlby9Qb3N0ZmlsdGVyL0Fic1RocmVzaD0wLjJiKXJlYXJyPXNjaGVtZV9Mb2NhbC9HZW8vQ3V0QWZpc2hhU25pcHBldD0xYjVyZWFycj1zY2hlbWVfTG9jYWwvR2VvL0hvdGVsQm9vc3Q9bWVhbl9jb252ZXJzaW9uXzEyd2IpcmVhcnI9c2NoZW1lX0xvY2FsL0dlby9Vc2VHZW9UcmF2ZWxSdWxlPTFqAnJ1cAGVAQAAAACdAc3MTD6gAQGoAQC9ASpEmwfCAQvq6oqSBLSz54m4Bg%3D%3D&sll=37.656245%2C55.759324&source=wizbiz_new_map_single&text=copyprint%20%D0%BA%D1%83%D1%80%D1%81%D0%BA%D0%B0%D1%8F&z=14": {
			ThumbnailUrl: "https://avatars.mds.yandex.net/get-altay/910613/2a000001622d8e17052f70e2a85a00c18bfe/L",
			Title:        "Copyprint, копировальный центр, Яковоапостольский пер., 17 — Яндекс.Карты",
		},
		"http://hidemyname.link/EN8f3g": {
			ThumbnailUrl: "https://hidemy.name/media/img/news/evropejskim-chinovnikam-rekomendujut-udalit-whatsapp-i-facebook-messenger.jpg",
			Title:        "Европейским чиновникам рекомендуют удалить WhatsApp и Facebook Messenger",
		},
		"https://cl.ly/391E2W3E111s": {
			ThumbnailUrl: "https://f.v1.n0.cdn.getcloudapp.com/items/3P0J3Z3y1o380L2l3X36/Image%202016-09-15%20at%2012.02.06%20PM.png",
			Title:        "Image 2016-09-15 at 12.02.06 PM.png",
		},
		"https://meduza.io/shapito/2020/04/22/vsegda-priyatno-posmotret-na-usatyh-loshadey-a-kogda-tak-malo-horoshih-novostey-tem-bolee": {
			ThumbnailUrl: "https://meduza.io/imgly/share/1587572704/shapito/2020/04/22/vsegda-priyatno-posmotret-na-usatyh-loshadey-a-kogda-tak-malo-horoshih-novostey-tem-bolee",
			Title:        "Всегда приятно посмотреть на усатых лошадей",
		},
		"https://zoom.us/j/97953527417?pwd=ZU9YaVFTNjBhb0t3dDh4bFJoZlVKQT09": {
			Title: "Join our Cloud HD Video Meeting",
		},
		"https://m.facebook.com/id77777/posts/10154846872521076": {
			//Title: "Telegram собрался на ICO. В последнее время я несколько раз слышал от участников рынка слухи о том, что Telegram план...",
			Title: "Для просмотра нужно войти или зарегистрироваться",
		},
		"https://facebook.com/permalink.php?story_fbid=1159419377602287&id=100006027061492": {
			// skip: has different servers for thumbnails
			// ThumbnailUrl: "https://external-arn2-1.xx.fbcdn.net/safe_image.php?d=AQHT-p_OOBeRCf7f&w=282&h=282&url=https%3A%2F%2Fmedia.giphy.com%2Fmedia%2FeIMB5D8d3SOR2%2Fgiphy.gif&cfs=1&_nc_cb=1&_nc_hash=AQE4cPaYplc9-Cit",
			Title:        "Мойрэрамос Гамос",
		},
		"https://github.com/meetecho/janus-gateway": {
			ThumbnailUrl: "https://repository-images.githubusercontent.com/16734696/cb265c80-651c-11ea-9c93-c54a2d2284e7",
			Title:        "meetecho/janus-gateway",
			//Title: "GitHub - meetecho/janus-gateway: Janus WebRTC Server",
		},
		// FIXME: stopped after 10 redirects
		//"https://developer.android.com/guide/topics/connectivity/telecom": {
		//	Title: "Telecom framework overview |",
		//},
	} {
		t.Run(link, func(t *testing.T) {
			result, err := c.Preview(link)
			if err != nil {
				t.Fatal(err)
			}
			if testCase.ThumbnailUrl != "" && result.ThumbnailUrl != testCase.ThumbnailUrl {
				t.Error("invalid ThumbnailUrl:\nwant:", testCase.ThumbnailUrl, "\ngot:", result.ThumbnailUrl)
			}
			if !strings.HasPrefix(result.Title, testCase.Title) {
				t.Error("invalid Title: want:", testCase.Title, " got:", result.Title)
			}
			if result.YoutubeId != testCase.YoutubeId {
				t.Error("invalid YoutubeId: want:", testCase.YoutubeId, "got:", result.YoutubeId)
			}
		})
	}
}
