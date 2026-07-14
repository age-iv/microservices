package model

type Comment struct {
    ID       int        `json:"id"`
    NewsID   int        `json:"news_id"`
    ParentID *int       `json:"parent_id,omitempty"`
    Text     string     `json:"text"`
    Created  int64      `json:"created_at"`
    Children []*Comment `json:"children,omitempty"`
}

type CreateCommentRequest struct {
    NewsID   int    `json:"news_id"`
    ParentID *int   `json:"parent_id,omitempty"`
    Text     string `json:"text"`
}
