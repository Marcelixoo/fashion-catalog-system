# Split variants and product indexes for near real-time indexing

## Context and Problem Statement

Before we start building the functionality, we decided to agree on a set of non-functional requirements that would guide our decisions on which architecture style to choose.

The most important architectural characteristics of our system are:

- **Durability:** fault tolerance & archivability of historical data.
- **Observability:** easy visualization of key metrics, e.g. number of daily updates, leading time for items to go from "pending approval" to "approved".
- **Agility:** strive for streamlined maintenance & isolated testing.
- **Resiliency:** service should not stop if dependencies are down or slow to respond (e.g. Cloud Translation API).

Our decision was weighed on the [Achitecture Styles Worksheet](https://www.developertoarchitect.com/downloads/architecture-styles-worksheet.pdf) and gut feelings coming from the experience of our engineers.

## Considered Options

- Option 1: Keep Both Indexes in Elasticsearch (Simple)
  • 🔧 Easier to implement — shared query language, unified stack
  • 🔁 Update both indices when variant stock/availability changes

- Option 2: Split Stores: Elasticsearch + Fast Lookup Store (Optimized)
  • Keep Product Index in Elasticsearch for full-text + facet search
  • Keep Variant Index in a fast key-value store like:
  • Redis (e.g. HGETALL variant:{id})
  • PostgreSQL with JSONB
  • DynamoDB for serverless setups

This is useful because:
• Variant data is small but changes often (price, availability, stock)
• K/V stores or small SQL tables are faster and cheaper for that

## Decision Outcome

Option 2, split stores with a Search Engine (NoSQL), e.g. Elasticsearch, Meilisearch, Solr, etc for product indexing and a Relational Database for variations lookups/filters.

## 1. Product Index (in Elasticsearch)

- Indexed document structure:

  - article_id
  - title
  - brand
  - category
  - facet_data:
    - available_sizes
    - available_colors
    - is_in_stock

- Used for:
  - Full-text search
  - Relevance scoring

## 2. Variant Index (in Redis / SQL / ES)

- Lookup per variant_id
- Schema:

  - article_id
  - variant_id
  - size
  - color
  - price
  - availability (boolean)
  - updated_at

- Used for:
  - Real-time price/stock display
  - Variant selection UI
  - Analytics + tracking
  - Filtering

## 3. Search Flow

1. Client sends query (e.g. `q=nike&size=42`)
2. Query hits Product Index in Elasticsearch
3. Elasticsearch returns top-N matching products with:

   - Precomputed facet data
   - article_ids

4. For each product:
   - App fetches relevant variants from Variant Index (Redis, etc.)
   - Filters/sorts client-side if needed
   - Displays available options

## 4. Update Flow

- When a variant is updated (e.g. stock or price):

  - ✅ Update Variant Index (real-time)
  - ❌ Do **not** reindex full product unless facet data changes

- When a product is updated (e.g. title or category):
  - ✅ Reindex in Product Index
  - 🔁 Optionally update variants if their metadata is affected
