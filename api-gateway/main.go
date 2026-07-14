package main

import (
    "log"
    "net/http"
    "os"
    "api-gateway/internal/handler"
    "api-gateway/internal/middleware"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/news", handler.GetNews)
    mux.HandleFunc("/news/", handler.GetNewsDetail)
    mux.HandleFunc("/comments", handler.CreateComment)

    wrapped := middleware.RequestIDMiddleware(middleware.LoggingMiddleware(mux))

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("API Gateway started on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, wrapped))
}
