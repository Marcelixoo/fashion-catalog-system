# fashion-catalog-system

## Outline

The Fashion Catalog System is a supercharged PIM-like software specialised in managing product catalog metadata from the fashion industry.

The system is inspired by a real-world and proprietary implementation of the PIM service that runs for the webshop https://momoxfashion.com (company where I work at the time of writing). Some features were removed and others were added based on my wishes for an open-source version of the product.

All functionatilly is exposed to collaborating systems via a unified REST API. Therefore, no user interface is available within this project.

## Use cases

For further reference, we'll identify a few of the actors involved in the main workflows supporded by the system: 

1. The **Downstream teams**, interested in items recently added to the catalog, e.g. Business Inteligence, Pricing, Inventory.
2. The **Merchandising team**, responsible for defining product types and reviewing generated content.
3. The **Warehouse team**, responsible for adding new items to the catalog with an initial set of basic attributes, e.g. colour, size, brand, type.

- **Warehouse Employee** adds a new batch of products.
  - Product attributes are filled in English.
  - Product attributes are translated into German & French before storage.
  - Rich descriptions are generated for each product.
- **Merchandising Employee** creates a product type definition.
- **Merchandising Employee** removes a product type definition.
- **Merchandising Employee** updates the definition of a product type.
- **Merchandising Employee** lists all existing product type definitions.
- **Merchandising Employee** lists products with pending approval.
- **Merchandising Employee** updates product attributes.
  - Step necessary for reviewing generated content & correcting potential mistakes.
- **Merchandising Employee** approves recently added product.
  - Only "approved" products are available on listings.
- Recently approved products are broadcasted to **Downstream teams**.
  - Notifications take the form of messages from a message queue.
  - Only "approved" products are notified.
  - **Downstream teams** have the ability to subscribe to "recently-added" topic to consume updates.

**Note:** Products are stored for undetermined period of time for auditing purposes and/or in case of returns.

## Non-functional requirements

1. Durability: fault tolerance & archivability of historical data.
2. Observability: easy visualization of key metrics, e.g. number of daily updates, leading time for items to go from "pending approval" to "approved".
3. Agility: strive for streamlined maintenance & isolated testing.
4. Resiliency: service should not stop if dependencies are down or slow to respond (e.g.Cloud Translation API).

 ## Visualize the structure

- High-level diagrams (C4 model with context, component levels)
- Main workflows
- Consider using icePanel (see https://icepanel.io/c4-model)
