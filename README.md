# Links checker

Takes title and og:image from any page:
```go
c := dyatl.Client()
result, err := c.Preview("https://meduza.io/shapito/2020/12/29/kazhduyu-zimu-v-rossii-edyat-zamerzshiy-doshirak-ved-chem-esche-zanyatsya-v-yakutske-ili-novosibirske")
if err != nil {
    panic(err)
}

fmt.Println("Title:", c.Title)
fmt.Println("ThumbnailUrl:", c.ThumbnailUrl)
```

Special cases like YouTube:
```go
c := dyatl.Client()
result, err := c.Preview("https://www.youtube.com/watch?v=eLAHSRmFFzE")
if err != nil {
    panic(err)
}

fmt.Println("Title:", c.Title)
fmt.Println("ThumbnailUrl:", c.ThumbnailUrl)
fmt.Println("YoutubeId:", c.YoutubeId)
```
