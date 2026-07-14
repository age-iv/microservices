package main

import (
    "log"
    "net/http"
    "os"
    "censorship-service/internal/handler"
    "censorship-service/internal/middleware"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/check", handler.Check)

    wrapped := middleware.RequestIDMiddleware(middleware.LoggingMiddleware(mux))

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8083"
    }
    log.Printf("Censorship service started on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, wrapped))
}
