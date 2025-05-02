# mini-search-lpatform

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

## Further information

- [Architectural Overview](ARCHITECTURE.md)
- [Architectural Decision Records (ADRs)](decisions/README.md)
- [Threat Model Analysis](THREAT_MODEL_ANALYSIS.md)
