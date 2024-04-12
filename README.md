# momox-catalog-service

## Outline

The Momox Catalog Service is responsible for centralizing product catalog metadata (attributes, categories and translations).
Catalog metadata is made available to other teams via a unified API that's constantly updated to catch-up with the pace of the warehouse workers.

## Functional requirements

1. Product attributes are stored in the German language.
2. Product catalog is enriched with additional visual attributes.
3. Product attributes are translated into English and French (EU and FR marketplaces).
4. Product categories can be assigned/updated at any time by the Merchandising team.
5. Product catalog updates are recorded.
6. Product catalog updates can be fetched by the BI team.

## Non-functional requirements

1. Durability: service should operate 24/7 with uptime of 99%.
2. Observability: easy visualization of daily updates & aging of data.
3. Agility: strive for streamlined maintenance & isolated testing.
4. Resiliency: service should not stop if dependencies are down (e.g. Enrichment API, Translation API).

 ## Visualize the structure

- High-level diagrams (C4 model with context, component levels)
- Main workflows
- Consider using icePanel (see https://icepanel.io/c4-model)
