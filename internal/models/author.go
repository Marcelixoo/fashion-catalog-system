package models

import "time"

type Author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func NewAuthor(AuthorId int, Name string) *Author {
	return &Author{
		ID:        AuthorId,
		Name:      Name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}
