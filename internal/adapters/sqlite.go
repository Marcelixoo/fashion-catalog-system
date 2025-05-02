package adapters

import (
	"database/sql"
	"mini-search-platform/internal/models"
)

type SQLliteAuthorsRepository struct {
	db *sql.DB
}

func NewSQLliteAuthorsRepository(db *sql.DB) *SQLliteAuthorsRepository {
	return &SQLliteAuthorsRepository{db: db}
}
func (r *SQLliteAuthorsRepository) Save(author *models.Author) (int, error) {
	query := `
		INSERT INTO authors (
			id,
			name, 
			created_at
		) VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query,
		author.ID,
		author.Name,
		author.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}
func (r *SQLliteAuthorsRepository) FindAuthorById(id int) (*models.Author, error) {
	query := `
		SELECT id, name, created_at
		FROM authors
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var author models.Author
	err := row.Scan(&author.ID, &author.Name, &author.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

type SQLliteArticleRepository struct {
	db *sql.DB
}

func NewSQLliteArticleRepository(db *sql.DB) *SQLliteArticleRepository {
	return &SQLliteArticleRepository{db: db}
}
func (r *SQLliteArticleRepository) Save(article *models.Article) (int, error) {
	query := `
		INSERT INTO articles (
			title, 
			body, 
			author_id,
			created_at
		) VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		article.Title,
		article.Body,
		article.AuthorID,
		article.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}
