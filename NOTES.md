## Scalability & Uptime Guarantee Checklist

1. Elastic Scaling
   • 📦 “Each customer can create their own index, and we scale horizontally across nodes using Kubernetes.”
   • 🧩 “We use index-per-tenant design, which simplifies isolation and scaling.”

2. Sharding
   • 🧱 “Large indexes are split into shards so we can parallelize indexing and querying.”
   • 🎯 “Shard count is optimized per index size — small tenants get 1 shard, big ones may get 5+.”

3. Replication
   • 🔁 “Each shard has replicas for high availability. If one node fails, the replica takes over automatically.”
   • 🔒 “We use 2 replicas for read redundancy and durability.”

4. Distributed Infrastructure
   • 🌍 “Indexes and replicas are spread across multiple zones or nodes to reduce single points of failure.”
   • 🛰 “We use cloud-native scheduling (e.g. Kubernetes, StatefulSets) to balance load.”

5. Fault Tolerance
   • ⚙️ “Our API layer uses health checks and timeouts. If a node is slow or down, we fail over or retry automatically.”
   • 🧯 “We maintain daily snapshots of all indexes to recover from critical failures.”

6. Fast Recovery
   • 💾 “We can restore an index from backups in minutes if needed.”
   • 🔄 “Read queries are routed to replicas when primaries are rebalancing or restarting.”

7. Performance for All Customers
   • 🚀 “Frequent queries and product lookups are cached at the edge or in Redis.”
   • 📊 “We monitor per-customer indexing and query load to adjust shard allocation dynamically.”

⸻

## Sharding euristics

🔧 What is a Shard?

A shard is a low-level unit of data and computation:
• Primary shard: holds the original data.
• Replica shard: copy of a primary shard used for failover and load balancing.

⸻

⚖️ Choosing Shard Count: Small vs. Large Indexes

✅ Small Index (up to ~1M documents / <2 GB total size)
• Recommended shards:
🧱 1 primary, 1–2 replicas
• Why:
Over-sharding small indexes adds unnecessary overhead (too many file handles, memory usage).
• Example:
Boutique shop with 20k products → 1 primary, 1 replica.
• Elasticsearch config:

"number_of_shards": 1,
"number_of_replicas": 1

⸻

⚖️ Medium Index (~1M–10M documents / 2–20 GB)
• Recommended shards:
🧱 2–3 primaries, 1–2 replicas
• Why:
Better parallelism for indexing and querying. One shard could become a bottleneck.
• Example:
Mid-sized marketplace with 5 million SKUs.

⸻

🚀 Large Index (>10M+ docs / >20 GB total size)
• Recommended shards:
🧱 5–10 primaries, depending on CPU, RAM, and query patterns
➕ 1–2 replicas
• Why:
Large documents must be split to:
• Distribute load across nodes
• Parallelize indexing/querying
• Fit in RAM/heap comfortably per shard
• Example:
Enterprise catalog with 50M variants across thousands of products.

⸻

📊 Rule of Thumb for Elasticsearch
• ⏱ Target shard size: ~10–30 GB per shard
• 🧠 Target shard count per node: ~20–30 shards per GB of heap (depends on usage)
• ✅ Shards too small → overhead (CPU/memory); too large → slower recovery and possible heap issues

⸻

🧩 Meilisearch Notes
• Meilisearch uses a single-shard design per index by default (no built-in sharding yet).
• To scale:
• Use index-per-tenant model (you already do this)
• Partition data at app layer (e.g. per category, region, brand)
• Horizontally scale Meilisearch nodes with load balancers or proxies (e.g. Meili Pro or your own sharding proxy)

⸻

🧠 Final Tip: Auto-Adjust Shards Based on Tenant Type

Tenant Type Document Count Recommended Shards Example Config
Hobby shop < 500k 1 primary, 1 replica Simple, fast access
Medium store 1M–5M 2–3 primaries Better query parallelism
Enterprise brand > 10M 5–10 primaries Load balancing + scalability

⸻

## Sharding by vendor type

📐 Sizing Calculator (Rules of Thumb)

Input:
• Number of documents
• Average document size (JSON payload size)

⸻

📊 1. Elasticsearch

Metric Estimate
Index size ~1.5–3× avg doc size × doc count
Shard count 1 shard per ~10–30 GB
Heap need ~1 GB heap per 20–30 shards
Field-heavy docs Consider 3× inflation vs raw JSON
Replication overhead Multiply index size by (1 + replica count)

📌 Example
• 5M products × 3 KB → ~15 GB raw
• ES index size: ~30–45 GB
• Shards: 2–4 primaries
• Heap: 2–4 GB minimum (with ~1 GB for each 25 shards)

⸻

⚡ 2. Meilisearch

Metric Estimate
Index size ~2–4× raw JSON (due to full-text + facets)
Max index size ~10–20M docs per node (recommended)
Memory usage ~1–2× index size in RAM (everything is mmap’d)
Sharding support ❌ No native sharding — you shard at app level

📌 Example
• 2M products × 3 KB → 6 GB raw → 12–15 GB Meili index
• RAM usage: 15–25 GB
• Shard manually via:
• Separate nodes per index (e.g. DE vs FR)
• Proxy and merge at app layer

⸻

☀️ 3. Apache Solr

Metric Estimate
Index size ~1.5–2.5× raw JSON
Shard count Similar to ES (~10–30 GB/shard)
Heap usage ~1–1.5 GB per 10M docs for basic queries
Replicas Use SolrCloud for HA

📌 Example
• 5M × 3 KB = 15 GB raw → ~30 GB Solr index
• Heap: 2–4 GB for indexing + search (adjust JVM -Xmx)
• Shards: 2–3 primaries, 1–2 replicas

⸻

🧮 Rough Sizing Calculator (Pseudocode)

func EstimateShardSize(docCount int, avgDocKB int) (esIndexSizeGB, meiliIndexSizeGB, solrIndexSizeGB float64) {
rawSizeGB := float64(docCount*avgDocKB) / 1024.0 / 1024.0
esIndexSizeGB = rawSizeGB * 2.5
meiliIndexSizeGB = rawSizeGB _ 3.0
solrIndexSizeGB = rawSizeGB _ 2.0
return
}

⸻

✅ Summary Table

Engine Inflated Index Size Shardable? Memory Footprint When to Choose
Elasticsearch ~2–3× ✅ Native ~1 GB/25 shards Large, filtered/faceted search
Meilisearch ~2–4× ❌ Manual RAM ≈ Index Size Lightweight, instant search
Solr ~1.5–2.5× ✅ Native JVM heap ~1.5 GB/10M Powerful but more ops-heavy
