package handler

import (
    "encoding/json"
    "net/http"
    "comments-service/internal/model"
    "comments-service/internal/storage"
    "strconv"
)

type Handler struct {
    store *storage.Storage
}

func New(store *storage.Storage) *Handler {
    return &Handler{store: store}
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
    var req model.CreateCommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
        return
    }
    if req.Text == "" || req.NewsID == 0 {
        http.Error(w, `{"error":"missing fields"}`, http.StatusBadRequest)
        return
    }
    id, err := h.store.CreateComment(req)
    if err != nil {
        http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
    newsIDStr := r.URL.Query().Get("news_id")
    newsID, err := strconv.Atoi(newsIDStr)
    if err != nil {
        http.Error(w, `{"error":"invalid news_id"}`, http.StatusBadRequest)
        return
    }
    flat, err := h.store.GetCommentsByNews(newsID)
    if err != nil {
        http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
        return
    }
    tree := buildTree(flat)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tree)
}

func buildTree(comments []model.Comment) []*model.Comment {
    nodeMap := make(map[int]*model.Comment)
    var roots []*model.Comment
    for i, c := range comments {
        cm := &c
        cm.Children = []*model.Comment{}
        nodeMap[c.ID] = cm
        comments[i] = c
    }
    for _, c := range comments {
        if c.ParentID != nil {
            if parent, ok := nodeMap[*c.ParentID]; ok {
                parent.Children = append(parent.Children, nodeMap[c.ID])
            } else {
                roots = append(roots, nodeMap[c.ID])
            }
        } else {
            roots = append(roots, nodeMap[c.ID])
        }
    }
    return roots
}
