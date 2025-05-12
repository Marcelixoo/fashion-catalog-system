```structurizr
workspace {

    model {
        user = person "API Consumer" {
            description "A developer, user, or system that sends search/indexing requests."
        }

        searchPlatform = softwareSystem "Tenant Search Platform" {
            description "A multi-tenant platform to search and index articles and products using Meilisearch or Elasticsearch."

            api = container "API Service" {
                description "Handles HTTP requests for indexing and searching."
                technology "Go + HTTPS"
            }

            router = container "Search Router" {
                description "Determines which backend to use (Meilisearch or Elasticsearch)."
                technology "Go"
            }

            registry = container "Index Registry" {
                description "Tracks index metadata per tenant."
                technology "Go + SQL"
            }

            redis = container "Redis (Routing Cache)" {
                description "Caches index metadata for backend routing."
                technology "Redis"
            }

            postgres = container "PostgreSQL (Registry DB)" {
                description "Stores tenant/index metadata."
                technology "PostgreSQL"
            }

            variantRedis = container "Redis (Variant Index)" {
                description "Stores variant data (stock, price) for fast access."
                technology "Redis"
            }

            variantPostgres = container "PostgreSQL (Variant Store)" {
                description "Stores persistent variant data for all articles."
                technology "PostgreSQL"
            }

            meili = container "Meilisearch" {
                description "Search backend for small/simple indexes."
                technology "Meilisearch"
            }

            elastic = container "Elasticsearch" {
                description "Search backend for large, complex indexes."
                technology "Elasticsearch"
            }

            user -> api "Sends search/index requests"
            api -> router "Forwards tenant requests"
            router -> redis "Looks up routing info"
            router -> registry "Fetches index metadata"
            registry -> postgres "Reads/writes registry data"

            router -> meili "Search product index (small)"
            router -> elastic "Search product index (large)"

            api -> variantRedis "Fetches variant data (real-time)"
            api -> variantPostgres "Persists and updates variants"
        }
    }

    views {
        systemContext searchPlatform {
            include *
            autoLayout lr
            title "Context Diagram - Multi-Tenant Search Platform with Variant Store"
        }

        container searchPlatform {
            include *
            autoLayout lr
            title "Container Diagram - Multi-Tenant Search Platform with Variant Store"
        }

        styles {
            element "Container" {
                background "#1168bd"
                color "#ffffff"
            }

            element "Person" {
                background "#08427b"
                color "#ffffff"
                shape "person"
            }

            element "Software System" {
                background "#438dd5"
                color "#ffffff"
            }
        }
    }
}
```

**Note:** Image generated using https://structurizr.com/dsl
