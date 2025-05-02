package models

import "time"

type Tag struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewTag(label string) *Tag {
	return &Tag{
		Label:     label,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

func (t *Tag) Update(label string) {
	t.Label = label
	t.UpdatedAt = time.Now().Format(time.RFC3339)
}

type TagsRepository interface {
	Save(*Tag) (int, error)
	FindById(id int) (*Tag, error)
	FindByLabel(label string) (*Tag, error)
	FindByLabels(labels []string) ([]*Tag, error)
	FindAll() ([]*Tag, error)
}
