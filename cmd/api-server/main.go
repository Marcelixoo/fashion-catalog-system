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

	_ = migrations.Migrate(db)

	articles := adapters.NewSQLliteArticleRepository(db)
	authors := adapters.NewSQLliteAuthorsRepository(db)
	// tags := adapters.NewSQLliteTagsRepository(db)

	engine := adapters.Init()

	sync := search.NewIndexSyncManager(engine, articles)

	r := gin.Default()
	// resource: articles
	r.POST("/articles", handlers.AddArticle(articles, authors, engine))
	r.POST("/articles/batch", handlers.AddArticles(articles, authors, engine, sync))

	// resource: authors
	r.POST("/authors", handlers.AddAuthor(authors))
	r.POST("/authors/batch", handlers.AddAuthors(authors))

	// resource: search
	r.GET("/search", handlers.SearchArticles(engine))

	r.Run(":8080")
}
