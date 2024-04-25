# fashion-catalog-service

## Outline

The Fashion Catalog Service is a supercharged PIM service specialised in managing product catalog metadata from the fashion industry.

The service is inspired by a real-world and proprietary implementation of the PIM service that runs for the webshop https://momoxfashion.com (company where I work at the time of writing). Some features were removed and others were added based on my wishes for an open-source version of the product.

All functionatilly is exposed to collaborating systems via a unified REST API. Therefore, no user interface is available within this project.

## Use cases

For further reference, we'll identify a few of the actors involved in the main workflows supporded by the service: 

1. The **Merchandising team**, responsible for defining product types and reviewing generated content.
2. The **Business Intelligence team**, interested in items recently added to the catalog.
3. The **Warehouse team**, responsible for adding new items to the catalog with an initial set of basic attributes, e.g. colour, size, brand, type.

- **Warehouse Service** adds a new batch of products.
  - Product attributes are filled in English.
  - Product attributes are translated into German & French before storage.
  - Rich descriptions are generated for each product.
- **Merchandising Service** creates a product type definition.
- **Merchandising Service** removes a product type definition.
- **Merchandising Service** updates the definition of a product type.
- **Merchandising Service** lists all existing product type definitions.
- **Merchandising Service** lists products with pending approval.
- **Merchandising Service** updates product attributes.
  - Step necessary for reviewing generated content & correcting potential mistakes.
- **Merchandising Service** approves recently added product.
  - Only "approved" products are available on listings.
- **BI Service** receives notifications about recently added products.
  - Notifications take the form of messages from a message queue.
  - Only "approved" products are notified.

**Note:** Products are stored for undetermined period of time for auditing purposes and/or in case of returns.

## Non-functional requirements

1. Durability: service should operate 24/7 with uptime of 99%.
2. Observability: easy visualization of daily updates & aging of data.
3. Agility: strive for streamlined maintenance & isolated testing.
4. Resiliency: service should not stop if dependencies are down (e.g. Enrichment API, Translation API).

 ## Visualize the structure

- High-level diagrams (C4 model with context, component levels)
- Main workflows
- Consider using icePanel (see https://icepanel.io/c4-model)
