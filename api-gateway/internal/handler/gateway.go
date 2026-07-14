package handler

import (
    "encoding/json"
    "net/http"
    "os"
    "strings"
    "api-gateway/internal/client"
)

var (
    newsServiceURL       = os.Getenv("NEWS_SERVICE_URL")
    commentsServiceURL   = os.Getenv("COMMENTS_SERVICE_URL")
    censorshipServiceURL = os.Getenv("CENSORSHIP_SERVICE_URL")
)

func GetNews(w http.ResponseWriter, r *http.Request) {
    query := r.URL.RawQuery
    url := newsServiceURL + "/news?" + query
    var data interface{}
    if err := client.GetJSON(url, &data); err != nil {
        http.Error(w, `{"error":"news service unavailable"}`, http.StatusServiceUnavailable)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func GetNewsDetail(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 3 {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    id := parts[2]
    newsURL := newsServiceURL + "/news/" + id
    commentsURL := commentsServiceURL + "/comments?news_id=" + id

    type result struct {
        data interface{}
        err  error
    }
    ch := make(chan result, 2)
    go func() {
        var news interface{}
        err := client.GetJSON(newsURL, &news)
        ch <- result{news, err}
    }()
    go func() {
        var comments interface{}
        err := client.GetJSON(commentsURL, &comments)
        ch <- result{comments, err}
    }()

    var newsData, commentsData interface{}
    var errs []string
    for i := 0; i < 2; i++ {
        res := <-ch
        if res.err != nil {
            errs = append(errs, res.err.Error())
            continue
        }
        if i == 0 {
            newsData = res.data
        } else {
            commentsData = res.data
        }
    }
    if len(errs) > 0 {
        http.Error(w, `{"error":"failed to fetch data"}`, http.StatusServiceUnavailable)
        return
    }
    resp := map[string]interface{}{
        "news":     newsData,
        "comments": commentsData,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
        return
    }
    text, ok := req["text"].(string)
    if !ok || text == "" {
        http.Error(w, `{"error":"text required"}`, http.StatusBadRequest)
        return
    }

    censorReq := map[string]string{"text": text}
    var censorResp map[string]string
    err := client.PostJSON(censorshipServiceURL+"/check", censorReq, &censorResp)
    if err != nil || censorResp["status"] != "approved" {
        http.Error(w, `{"error":"comment rejected by censorship"}`, http.StatusBadRequest)
        return
    }

    var createResp map[string]int
    err = client.PostJSON(commentsServiceURL+"/comments", req, &createResp)
    if err != nil {
        http.Error(w, `{"error":"comment creation failed"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createResp)
}
