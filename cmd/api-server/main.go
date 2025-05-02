package main

import (
	"fmt"
	"mini-search-platform/internal/adapters"
	"mini-search-platform/internal/handlers"
	"mini-search-platform/internal/migrations"
	"mini-search-platform/pkg/sqlite"

	"mini-search-platform/internal/search"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("running api-server...")

	db, err := sqlite.Init()
	if err != nil {
		panic(err)
	}
	defer sqlite.Close(db)

	err = migrations.Migrate(db)
	if err != nil {
		panic(err)
	}

	articles := adapters.NewSQLliteArticleRepository(db)
	authors := adapters.NewSQLliteAuthorsRepository(db)
	tags := adapters.NewSQLliteTagsRepository(db)

	engine := adapters.Init()

	sync := search.NewIndexSyncManager(engine, articles, tags)

	r := gin.Default()
	// resource: articles
	r.POST("/articles", handlers.AddArticle(articles, authors, tags, sync))
	r.POST("/articles/batch", handlers.AddArticles(articles, authors, tags, sync))

	// resource: authors
	r.POST("/authors", handlers.AddAuthor(authors))
	r.POST("/authors/batch", handlers.AddAuthors(authors))

	// resource: tags
	r.POST("/tags", handlers.AddTag(tags))
	r.PATCH("/tags/:label", handlers.UpdateTagWithLabel(tags, sync))
	r.POST("/tags/batch", handlers.AddTagsInBatch(tags))
	r.GET("/tags", handlers.ListAllTags(tags))
	r.GET("/tags/:label", handlers.GetTagByLabel(tags))
	r.GET("/tags/:label/articles", handlers.FindArticlesByLabels(articles, tags))

	// resource: search
	r.GET("/search", handlers.SearchArticles(engine))

	r.Run(":8080")
}
