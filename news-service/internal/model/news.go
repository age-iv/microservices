package model

type News struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    PubTime int64  `json:"pub_time"`
    Link    string `json:"link"`
    Source  string `json:"source"`
}

type Pagination struct {
    Page       int `json:"page"`
    TotalPages int `json:"total_pages"`
    PerPage    int `json:"per_page"`
}

type NewsResponse struct {
    News       []News     `json:"news"`
    Pagination Pagination `json:"pagination"`
}
