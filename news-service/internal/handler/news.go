package handler

import (
    "encoding/json"
    "net/http"
    "news-service/internal/model"
    "news-service/internal/storage"
    "strconv"
    "strings"
)

const perPage = 15

type Handler struct {
    store *storage.Storage
}

func New(store *storage.Storage) *Handler {
    return &Handler{store: store}
}

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    page := 1
    if pageStr != "" {
        p, err := strconv.Atoi(pageStr)
        if err == nil && p > 0 {
            page = p
        }
    }
    search := r.URL.Query().Get("s")

    news, total, err := h.store.GetNews(page, perPage, search)
    if err != nil {
        http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
        return
    }
    totalPages := total / perPage
    if total%perPage != 0 {
        totalPages++
    }
    resp := model.NewsResponse{
        News: news,
        Pagination: model.Pagination{
            Page:       page,
            TotalPages: totalPages,
            PerPage:    perPage,
        },
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetNewsDetail(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 3 {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    id, err := strconv.Atoi(parts[2])
    if err != nil {
        http.Error(w, "bad id", http.StatusBadRequest)
        return
    }
    news, err := h.store.GetNewsByID(id)
    if err != nil {
        http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(news)
}
