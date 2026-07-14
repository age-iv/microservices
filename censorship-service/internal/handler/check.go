package handler

import (
    "encoding/json"
    "net/http"
    "strings"
)

var forbiddenWords = []string{"qwerty", "йцукен", "zxvbnm"}

type CheckRequest struct {
    Text string `json:"text"`
}

type CheckResponse struct {
    Status string `json:"status"`
}

func Check(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
        return
    }
    var req CheckRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
        return
    }
    text := strings.ToLower(req.Text)
    for _, word := range forbiddenWords {
        if strings.Contains(text, word) {
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(CheckResponse{Status: "rejected"})
            return
        }
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(CheckResponse{Status: "approved"})
}
