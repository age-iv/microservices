package storage

import (
    "database/sql"
    "fmt"
    "news-service/internal/model"
    "os"
    _ "github.com/lib/pq"
)

type Storage struct {
    db *sql.DB
}

func New() (*Storage, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    s := &Storage{db: db}
    if err = s.createTable(); err != nil {
        return nil, err
    }
    return s, nil
}

func (s *Storage) createTable() error {
    q := `CREATE TABLE IF NOT EXISTS news (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        pub_time BIGINT NOT NULL,
        link TEXT UNIQUE,
        source TEXT
    )`
    _, err := s.db.Exec(q)
    return err
}

func (s *Storage) SaveNews(news model.News) error {
    _, err := s.db.Exec(`INSERT INTO news (title, content, pub_time, link, source)
        VALUES ($1, $2, $3, $4, $5) ON CONFLICT (link) DO NOTHING`,
        news.Title, news.Content, news.PubTime, news.Link, news.Source)
    return err
}

func (s *Storage) GetNews(page, perPage int, search string) ([]model.News, int, error) {
    offset := (page - 1) * perPage
    var total int
    countQuery := "SELECT COUNT(*) FROM news"
    dataQuery := "SELECT id, title, content, pub_time, link, source FROM news"
    var args []interface{}
    if search != "" {
        filter := " WHERE title ILIKE $1"
        countQuery += filter
        dataQuery += filter
        args = append(args, "%"+search+"%")
    }
    err := s.db.QueryRow(countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    dataQuery += " ORDER BY pub_time DESC LIMIT $2 OFFSET $3"
    args = append(args, perPage, offset)
    rows, err := s.db.Query(dataQuery, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    var news []model.News
    for rows.Next() {
        var n model.News
        err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.PubTime, &n.Link, &n.Source)
        if err != nil {
            return nil, 0, err
        }
        news = append(news, n)
    }
    return news, total, nil
}

func (s *Storage) GetNewsByID(id int) (model.News, error) {
    var n model.News
    err := s.db.QueryRow("SELECT id, title, content, pub_time, link, source FROM news WHERE id = $1", id).
        Scan(&n.ID, &n.Title, &n.Content, &n.PubTime, &n.Link, &n.Source)
    return n, err
}
