package main

import (
	"fmt"
	"mini-search-platform/internal/adapters"
	"mini-search-platform/internal/handlers"
	"mini-search-platform/internal/migrations"
	"mini-search-platform/pkg/sqlite"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("running api-server...")

	engine := adapters.Init()

	db, err := sqlite.Init()
	if err != nil {
		panic(err)
	}
	defer sqlite.Close(db)

	_ = migrations.Migrate(db)

	r := gin.Default()
	r.POST("/articles", handlers.AddArticle(adapters.NewSQLliteArticleRepository(db), adapters.NewSQLliteAuthorsRepository(db), engine))
	r.POST("/articles/batch", handlers.AddArticles(adapters.NewSQLliteArticleRepository(db), adapters.NewSQLliteAuthorsRepository(db), engine))
	r.POST("/authors", handlers.AddAuthor(adapters.NewSQLliteAuthorsRepository(db)))
	r.POST("/authors/batch", handlers.AddAuthors(adapters.NewSQLliteAuthorsRepository(db)))
	r.GET("/search", handlers.SearchArticles(engine))

	r.Run(":8080")
}
