package models

import "time"

type Article struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Author    string   `json:"author"`
	AuthorID  int      `json:"author_id"`
	CreatedAt string   `json:"created_at"`
	Tags      []string `json:"tags"`
}

func NewArticle(title, body string, author *Author, tags []string) *Article {
	return &Article{
		Title:     title,
		Body:      body,
		Author:    author.Name,
		AuthorID:  author.ID,
		Tags:      tags,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}
