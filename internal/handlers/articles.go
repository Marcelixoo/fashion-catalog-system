package handlers

import (
	"context"
	"mini-search-platform/internal/models"
	"mini-search-platform/internal/search"
	"mini-search-platform/pkg/retry"

	"github.com/gin-gonic/gin"
)

type AuthorsFinder interface {
	FindAuthorById(id int) (*models.Author, error)
}

type ArticleInput struct {
	Title    string   `json:"title" binding:"required"`
	Body     string   `json:"body" binding:"required"`
	AuthorID int      `json:"author_id" binding:"required"`
	Author   string   `json:"author"`
	Tags     []string `json:"tags"`
}

type AddArticlesSummary struct {
	TotalInserted int `json:"total_inserted"`
	TotalFailed   int `json:"total_failed"`
}

type AddArticlesResponse struct {
	Summary  AddArticlesSummary        `json:"summary"`
	Inserted []*models.Article         `json:"inserted"`
	Failed   []map[string]ArticleInput `json:"failed"`
}

func AddArticles(repository models.ArticleRepository, finder AuthorsFinder, tagsRepository models.TagsRepository, sync *search.IndexSyncManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputs []ArticleInput
		if err := c.ShouldBindJSON(&inputs); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var inserted []*models.Article
		var failed = []map[string]ArticleInput{}
		for _, input := range inputs {
			author, err := finder.FindAuthorById(input.AuthorID)
			if err != nil {
				failed = append(failed, map[string]ArticleInput{
					"author not found": input,
				})
				continue
			}

			tags, err := tagsRepository.FindByLabels(input.Tags)
			if err != nil {
				failed = append(failed, map[string]ArticleInput{
					"tags not found": input,
				})
				continue
			}

			article := models.NewArticle(input.Title, input.Body, author, tags)

			lastInsertedId, err := repository.Save(article)
			if err != nil {
				failed = append(failed, map[string]ArticleInput{
					err.Error(): input,
				})
				continue
			}

			article.ID = lastInsertedId
			inserted = append(inserted, article)
		}

		resync := func(articlesToSync []*models.Article) error {
			operation := func() error {
				return sync.SyncAfterArticlesChanged(articlesToSync)
			}
			return retry.WithBackoff(context.Background(), operation)
		}
		go resync(inserted)

		c.JSON(201, AddArticlesResponse{
			Summary: AddArticlesSummary{
				TotalInserted: len(inserted),
				TotalFailed:   len(failed),
			},
			Inserted: inserted,
			Failed:   failed,
		})
	}
}

func AddArticle(repository models.ArticleRepository, finder AuthorsFinder, tagsRepository models.TagsRepository, sync *search.IndexSyncManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input ArticleInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		author, err := finder.FindAuthorById(input.AuthorID)
		if err != nil {
			c.JSON(400, gin.H{"error": "Author not found"})
			return
		}

		tags, err := tagsRepository.FindByLabels(input.Tags)
		if err != nil {
			c.JSON(400, gin.H{"error": "Could not find one (or more) tags"})
			return
		}

		article := models.NewArticle(input.Title, input.Body, author, tags)

		lastInsertedId, err := repository.Save(article)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save article"})
			return
		}

		article.ID = lastInsertedId

		resync := func(articlesToSync []*models.Article) error {
			operation := func() error {
				return sync.SyncAfterArticlesChanged(articlesToSync)
			}
			return retry.WithBackoff(context.Background(), operation)
		}
		go resync([]*models.Article{article})

		c.JSON(201, article)
	}
}
