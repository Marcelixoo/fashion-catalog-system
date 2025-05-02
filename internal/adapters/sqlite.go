package adapters

import (
	"database/sql"
	"fmt"
	"mini-search-platform/internal/models"
	"strings"
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
		) VALUES (?, ?, ?, ?);
	`

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(query,
		article.Title,
		article.Body,
		article.AuthorID,
		article.CreatedAt,
	)
	if err != nil {
		return 0, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	query = `
		INSERT INTO article_tags (article_id, tag_id)
		VALUES (?, ?)
	`

	for _, tag := range article.Tags {
		_, err := tx.Exec(query, lastInsertedId, tag.ID)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(lastInsertedId), err
}

func (r *SQLliteArticleRepository) FindByTag(tag *models.Tag) ([]*models.Article, error) {
	query := `
		SELECT
			a.id,
			a.title, 
			a.body, 
			a.author_id,
			au.name,
			a.created_at,
			t.id,
			t.label,
			t.created_at,
			t.updated_at
		FROM articles a
		JOIN authors au ON a.author_id = au.id
		JOIN tags t ON at.tag_id = t.id
		JOIN article_tags at ON a.id = at.article_id
		WHERE a.id IN (
			SELECT at.article_id
			FROM article_tags at
			WHERE at.tag_id = ?
		)
	`

	rows, err := r.db.Query(query, tag.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articleMap := make(map[int]*models.Article)

	for rows.Next() {
		var (
			articleID                            int
			title, body                          string
			authorID                             int
			authorName, createdAt                string
			tagID                                int
			tagLabel, tagCreatedAt, tagUpdatedAt string
		)

		err := rows.Scan(
			&articleID, &title, &body,
			&authorID, &authorName, &createdAt,
			&tagID, &tagLabel, &tagCreatedAt, &tagUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		article, exists := articleMap[articleID]
		if !exists {
			article = &models.Article{
				ID:        articleID,
				Title:     title,
				Body:      body,
				AuthorID:  authorID,
				Author:    authorName,
				CreatedAt: createdAt,
				Tags:      []*models.Tag{},
			}
			articleMap[articleID] = article
		}

		article.Tags = append(article.Tags, &models.Tag{
			ID:        tagID,
			Label:     tagLabel,
			CreatedAt: tagCreatedAt,
			UpdatedAt: tagUpdatedAt,
		})
	}

	// Convert map to slice
	var articles []*models.Article
	for _, a := range articleMap {
		articles = append(articles, a)
	}

	return articles, nil

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

func (r *SQLliteTagsRepository) FindByLabels(labels []string) ([]*models.Tag, error) {
	if len(labels) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(labels)-1) + "?"

	query := fmt.Sprintf(`
		SELECT id, label, created_at, updated_at
		FROM tags
		WHERE label IN (%s)
	`, placeholders)

	args := make([]interface{}, len(labels))
	for i, label := range labels {
		args[i] = label
	}

	rows, err := r.db.Query(query, args...)
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
