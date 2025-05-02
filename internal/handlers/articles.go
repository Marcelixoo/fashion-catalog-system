package handlers

import (
	"mini-search-platform/internal/models"
	"mini-search-platform/internal/search"

	"github.com/gin-gonic/gin"
)

type ArticleRepository interface {
	Save(*models.Article) (int, error)
}

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

func AddArticles(repository ArticleRepository, finder AuthorsFinder, engine search.SearchEngine) gin.HandlerFunc {
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

			article := models.NewArticle(input.Title, input.Body, author, input.Tags)

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

		err := engine.IndexArticles(inserted)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to index articles"})
			return
		}

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

func AddArticle(repository ArticleRepository, finder AuthorsFinder, engine search.SearchEngine) gin.HandlerFunc {
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

		article := models.NewArticle(input.Title, input.Body, author, input.Tags)

		lastInsertedId, err := repository.Save(article)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to insert article"})
			return
		}

		article.ID = lastInsertedId

		err = engine.IndexArticles([]*models.Article{article})
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to index article"})
			return
		}

		c.JSON(201, article)
	}
}
