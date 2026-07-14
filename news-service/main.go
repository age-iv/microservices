package main

import (
    "log"
    "net/http"
    "os"
    "news-service/internal/handler"
    "news-service/internal/middleware"
    "news-service/internal/rss"
    "news-service/internal/storage"
)

func main() {
    store, err := storage.New()
    if err != nil {
        log.Fatal(err)
    }
    rss.Start(store)

    h := handler.New(store)
    mux := http.NewServeMux()
    mux.HandleFunc("/news", h.GetNews)
    mux.HandleFunc("/news/", h.GetNewsDetail)

    wrapped := middleware.RequestIDMiddleware(middleware.LoggingMiddleware(mux))

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8081"
    }
    log.Printf("News service started on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, wrapped))
}
