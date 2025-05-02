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

func (r *SQLliteArticleRepository) FindByTags(tags []*models.Tag) ([]*models.Article, error) {
	return []*models.Article{}, nil
}

type SQLliteTagsRepository struct {
	db *sql.DB
}

func NewSQLliteTagsRepository(db *sql.DB) *SQLliteTagsRepository {
	return &SQLliteTagsRepository{db: db}
}

func (r *SQLliteTagsRepository) Save(tag *models.Tag) (int, error) {
	query := `
		INSERT INTO tags (label, updated_at, created_at)
		VALUES (?, ?, ?)
		ON CONFLICT(label) DO UPDATE SET
			label = ?,
			updated_at = ?;
	`

	result, err := r.db.Exec(query,
		tag.Label,
		tag.UpdatedAt,
		tag.CreatedAt,
		tag.Label,
		tag.UpdatedAt,
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

func (r *SQLliteTagsRepository) FindByLabel(label string) (*models.Tag, error) {
	query := `
		SELECT id, label, created_at, updated_at
		FROM tags
		WHERE label = ?
	`
	row := r.db.QueryRow(query, label)

	var tag models.Tag
	err := row.Scan(&tag.ID, &tag.Label, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *SQLliteTagsRepository) FindById(id int) (*models.Tag, error) {
	query := `
		SELECT id, label, created_at, updated_at
		FROM tags
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var tag models.Tag
	err := row.Scan(&tag.ID, &tag.Label, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *SQLliteTagsRepository) FindAll() ([]*models.Tag, error) {
	query := `
		SELECT id, label, created_at, updated_at
		FROM tags
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Label, &tag.CreatedAt, &tag.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}
