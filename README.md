# mini-search-platform

## Outline

The Mini Search Platform is a news platform specialised in making articles  searchable.

The system is inspired by a real-world and proprietary implementation of a Search Platform that runs for the webshop https://momoxfashion.com (company where I work at the time of writing). Some features were removed and others were added based on my wishes for an open-source version of the product.

All functionatilly is exposed to collaborating systems via a unified REST API. Therefore, no user interface is available within this project.

## Use cases

📝 Articles
`POST /articles`
Create a single article with metadata (title, body, author, tags).
Automatically syncs the article to the search index.

`POST /articles/batch`
Create multiple articles in one request.
Each article includes an author and a list of tags.
All articles are synced to the search engine after insert.

👤 Authors
`POST /authors`
Create a new author with a unique ID and name.

`POST /authors/batch`
Batch insert multiple authors.
Useful during initial data ingestion or import operations.

🏷️ Tags
`POST /tags`
Add a new tag by label.
If the tag already exists, it can be updated or rejected depending on backend logic.

`PATCH /tags/:label`
Update the label of an existing tag.
Triggers a background resync of related articles in the search index to reflect the updated tag.

`POST /tags/batch`
Batch insert multiple tags.
Returns a summary of how many were inserted vs. failed.

`GET /tags`
List all tags stored in the database.
Supports use in filtering UIs or autocomplete features.

`GET /tags/:label`
Retrieve a single tag by its label.
Useful for checking if a tag exists before assigning it to an article.

`GET /tags/:label/articles`
Fetch all articles that are associated with a tag matching the provided label.
Returns full articles, each with a list of their tags (not just the matching one).

🔎 Search
`GET /search`
Perform a full-text search across articles via the search engine.
Supports keyword queries and may include filters (e.g., by tag or author) depending on implementation.


## Non-functional requirements

1. Durability: fault tolerance & archivability of historical data.
3. Agility: strive for streamlined maintenance & isolated testing.
4. Resiliency: service should not stop if dependencies are down or slow to respond (e.g.Cloud Translation API).

## Local setup

### 1. Start Meilisearch container

docker run -it --rm -p 7700:7700 getmeili/meilisearch

### 2. Start the application
go run cmd/api-server/main.go

## Sample requests

Sample requests are available at [examples/mini-search-platform.postman_collection.json](examples/mini-search-platform.postman_collection.json)

🗂️ Entity-Relationship Model (ER Model)
This system models a publishing platform with articles, authors, and tags. It supports a many-to-many relationship between articles and tags.

📊 Tables Overview
### 👤 `authors`

| Column     | Type      | Constraints                          |
|------------|-----------|--------------------------------------|
| `id`       | INTEGER   | Primary key, Auto-increment          |
| `name`     | TEXT      | Not null, Unique                     |
| `created_at` | TIMESTAMP | Defaults to `CURRENT_TIMESTAMP`     |

---

### 📝 `articles`

| Column      | Type      | Constraints                              |
|-------------|-----------|------------------------------------------|
| `id`        | INTEGER   | Primary key, Auto-increment              |
| `title`     | TEXT      | Not null                                 |
| `body`      | TEXT      | Not null                                 |
| `author_id` | INTEGER   | Foreign key → `authors(id)`, Not null    |
| `created_at`| TIMESTAMP | Defaults to `CURRENT_TIMESTAMP`          |

---

### 🏷️ `tags`

| Column       | Type      | Constraints                          |
|--------------|-----------|--------------------------------------|
| `id`         | INTEGER   | Primary key, Auto-increment          |
| `label`      | TEXT      | Not null, Unique                     |
| `created_at` | TIMESTAMP | Defaults to `CURRENT_TIMESTAMP`      |
| `updated_at` | TIMESTAMP | Nullable                             |

---

### 🔗 `article_tags`

| Column      | Type    | Constraints                                                  |
|-------------|---------|--------------------------------------------------------------|
| `article_id`| INTEGER | Primary key (with `tag_id`), Foreign key → `articles(id)`    |
| `tag_id`    | INTEGER | Primary key (with `article_id`), Foreign key → `tags(id)`    |



🔁 Relationships
- **1 Author → many Articles**
Each article is written by a single author (articles.author_id → authors.id).
- **Many Articles ↔ Many Tags**
Represented via the join table article_tags.

📘 Diagram (Text-based)
```markdown
authors
 └───< articles
           └───< article_tags >───┐
                                  │
                                tags
```

## Further information

- [Architectural Overview](ARCHITECTURE.md)
- [Architectural Decision Records (ADRs)](decisions/README.md)
- [Threat Model Analysis](THREAT_MODEL_ANALYSIS.md)
