# Which architecture style to adopt?

## Context and Problem Statement

Before we start building the functionality, we decided to agree on a set of non-functional requirements that would guide our decisions on which architecture style to choose.

The most important architectural characteristics of our system are:
- **Durability:** fault tolerance & archivability of historical data.
- **Observability:** easy visualization of key metrics, e.g. number of daily updates, leading time for items to go from "pending approval" to "approved".
- **Agility:** strive for streamlined maintenance & isolated testing.
- **Resiliency:** service should not stop if dependencies are down or slow to respond (e.g. Cloud Translation API).

Our decision was weighed on the [Achitecture Styles Worksheet](https://www.developertoarchitect.com/downloads/architecture-styles-worksheet.pdf) and gut feelings coming from the experience of our engineers.

## Considered Options

* Ports & Adapters
* Event-Driven
* Microservices

![Spectrum of different architecture style patterns](https://github.com/Marcelixoo/momox-catalog-service/assets/29285152/5c71d44e-d091-4042-9f8e-a7803ea6d578)

## Decision Outcome

Chosen option: "Microservices", because of specific steps of our workflow requiring more processing than others. During the enrichment of product information two external services will be required– Generative AI and Cloud Translating API–which add additional latency yet should not block further processing of incoming requests.

Therefore, having separate deployable units seems to be the best way to go. On top of that, we see the benefit of easer evolvability & testability of the services.
