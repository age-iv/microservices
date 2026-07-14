package rss

import (
    "encoding/xml"
    "io"
    "net/http"
    "time"
    "news-service/internal/model"
    "news-service/internal/storage"
)

type rssFeed struct {
    Channel struct {
        Item []struct {
            Title       string `xml:"title"`
            Description string `xml:"description"`
            Link        string `xml:"link"`
            PubDate     string `xml:"pubDate"`
        } `xml:"item"`
    } `xml:"channel"`
}

var feeds = []string{
    "https://lenta.ru/rss",
    "https://www.interfax.ru/rss.asp",
}

func Start(st *storage.Storage) {
    go func() {
        for {
            for _, url := range feeds {
                fetchAndStore(url, st)
            }
            time.Sleep(10 * time.Minute)
        }
    }()
}

func fetchAndStore(url string, st *storage.Storage) {
    resp, err := http.Get(url)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    var feed rssFeed
    xml.Unmarshal(body, &feed)
    for _, item := range feed.Channel.Item {
        pubTime, _ := time.Parse(time.RFC1123Z, item.PubDate)
        n := model.News{
            Title:   item.Title,
            Content: item.Description,
            PubTime: pubTime.Unix(),
            Link:    item.Link,
            Source:  url,
        }
        st.SaveNews(n)
    }
}
