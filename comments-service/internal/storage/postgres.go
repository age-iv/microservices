package storage

import (
    "database/sql"
    "fmt"
    "os"
    "time"
    "comments-service/internal/model"
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
    q := `CREATE TABLE IF NOT EXISTS comments (
        id SERIAL PRIMARY KEY,
        news_id INT NOT NULL,
        parent_id INT,
        text TEXT NOT NULL,
        created_at BIGINT NOT NULL
    )`
    _, err := s.db.Exec(q)
    return err
}

func (s *Storage) CreateComment(req model.CreateCommentRequest) (int, error) {
    var id int
    err := s.db.QueryRow(`INSERT INTO comments (news_id, parent_id, text, created_at)
        VALUES ($1, $2, $3, $4) RETURNING id`,
        req.NewsID, req.ParentID, req.Text, time.Now().Unix()).Scan(&id)
    return id, err
}

func (s *Storage) GetCommentsByNews(newsID int) ([]model.Comment, error) {
    rows, err := s.db.Query(`SELECT id, news_id, parent_id, text, created_at FROM comments WHERE news_id = $1 ORDER BY created_at ASC`, newsID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var comments []model.Comment
    for rows.Next() {
        var c model.Comment
        if err := rows.Scan(&c.ID, &c.NewsID, &c.ParentID, &c.Text, &c.Created); err != nil {
            return nil, err
        }
        comments = append(comments, c)
    }
    return comments, nil
}
