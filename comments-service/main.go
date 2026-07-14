package main

import (
    "log"
    "net/http"
    "os"
    "comments-service/internal/handler"
    "comments-service/internal/middleware"
    "comments-service/internal/storage"
)

func main() {
    store, err := storage.New()
    if err != nil {
        log.Fatal(err)
    }
    h := handler.New(store)

    mux := http.NewServeMux()
    mux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.GetComments(w, r)
        case http.MethodPost:
            h.CreateComment(w, r)
        default:
            http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
        }
    })

    wrapped := middleware.RequestIDMiddleware(middleware.LoggingMiddleware(mux))

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8082"
    }
    log.Printf("Comments service started on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, wrapped))
}
