package search

import (
	"mini-search-platform/internal/models"
)

var (
	ARTICLES_INDEX_NAME = "articles"
)

type SearchEngine interface {
	Search(q string, options SearchOptions) (SearchResponse, error)
	IndexArticles(articles []*models.Article) error
}

type SearchOptions struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Sort   []string `json:"sort"`
	Filter string   `json:"filter"`
	Facets string   `json:"facets"`
}

type SearchHit struct {
	ID     int          `json:"id"`
	Title  string       `json:"title"`
	Author string       `json:"author"`
	Body   string       `json:"body"`
	Tags   []models.Tag `json:"tags"`
}

type SearchHits struct {
	Hits []SearchHit `json:"hits"`
}

type SearchResponse struct {
	Query  string      `json:"query"`
	Hits   []SearchHit `json:"hits"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Total  int         `json:"total"`
}

type IndexSyncManager struct {
	Engine             SearchEngine
	ArticlesRepository models.ArticleRepository
	TagsRepository     models.TagsRepository
}

func NewIndexSyncManager(engine SearchEngine, articlesRepository models.ArticleRepository, tagsRepository models.TagsRepository) *IndexSyncManager {
	return &IndexSyncManager{
		Engine:             engine,
		ArticlesRepository: articlesRepository,
		TagsRepository:     tagsRepository,
	}
}

func (m *IndexSyncManager) SyncAfterTagsChanged(tagToSync *models.Tag) error {
	articles, err := m.ArticlesRepository.FindByTag(tagToSync)
	if err != nil {
		return err
	}

	return m.Engine.IndexArticles(articles)
}

func (m *IndexSyncManager) SyncAfterArticlesChanged(articlesToSync []*models.Article) error {
	err := m.Engine.IndexArticles(articlesToSync)
	if err != nil {
		return err
	}

	return nil
}
