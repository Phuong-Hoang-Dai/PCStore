# search_service

Search microservice for the DDStore PC Store. It indexes products into
Elasticsearch and exposes a full-text search API over them.

## Architecture

Layered, dependency-inverted structure (mirrors `catalog_service`):

```
search_service/
├── cmd/
│   └── main.go                     # entrypoint: load config, connect ES, run HTTP server
├── configs/
│   ├── config.go                   # env-based configuration
│   └── example.env                 # copy to .env for local dev
├── db/
│   └── elastic.go                  # Elasticsearch client setup + health check
├── internal/
│   ├── handler/                    # HTTP layer (Gin)
│   │   ├── router.go               # routes + server wiring
│   │   └── search_handler.go       # request/response handling
│   ├── service/                    # business logic
│   │   ├── interface.go            # SearchService + ProductRepos interfaces
│   │   └── search.go               # SearchService implementation
│   ├── repos/                      # persistence (Elasticsearch)
│   │   ├── elasticProductRepos.go  # index / bulk / delete / search
│   │   └── mapping.go              # products index mapping
│   └── model/                      # domain models
│       ├── product.go
│       ├── search.go               # Paging, SearchRequest, SearchResult
│       └── const.go
├── Dockerfile
└── go.mod
```

Flow: `handler → service → repos → Elasticsearch`. Each layer depends on the
interface of the one below, so the service can be tested with a mock repo.

## Configuration

Set via environment variables (or a `configs/.env` file locally):

| Variable       | Default                 | Description                        |
| -------------- | ----------------------- | ---------------------------------- |
| `APP_NAME`     | `PC Store Search`       | Service name                       |
| `PORT`         | `8889`                  | HTTP listen port                   |
| `ES_ADDRESSES` | `http://localhost:9200` | Comma-separated Elasticsearch URLs |
| `ES_USERNAME`  | _(empty)_               | Basic-auth username (optional)     |
| `ES_PASSWORD`  | _(empty)_               | Basic-auth password (optional)     |
| `ES_INDEX`     | `products`              | Index name                         |

## Running locally

```bash
cp configs/example.env configs/.env   # then edit if needed
go mod download
go run ./cmd
```

The service creates the `products` index (with mapping) on startup if it
does not already exist.

## API

Base path: `/search/products`

### Insert one product

```
POST /search/products
Content-Type: application/json

{
  "id": "6650f0c3a1b2c3d4e5f60718",
  "name": "ASUS ROG Strix G16",
  "description": "Gaming laptop with RTX 4070",
  "brand": "ASUS",
  "cate": { "id": 1, "name": "Laptop" },
  "price": 1799.99,
  "type": "standard",
  "images": [{ "url": "https://.../g16.png", "role": "thumbnail" }]
}
```

`id` is optional — if omitted, Elasticsearch generates one and it is returned.

### Insert many products (bulk)

```
POST /search/products/bulk
Content-Type: application/json

[ { ...product... }, { ...product... } ]
```

### Delete a product

```
DELETE /search/products/:id
```

### Search products

```
GET /search/products?q=gaming+laptop&brand=ASUS&type=standard&cate_id=1&min_price=1000&max_price=2000&limit=10&offset=0
```

| Param                 | Description                                       |
| --------------------- | ------------------------------------------------- |
| `q`                   | Full-text query over name, brand, description     |
| `brand`               | Exact brand filter                                |
| `type`                | Product type filter (`standard` / `composite`)    |
| `cate_id`             | Category id filter                                |
| `min_price`/`max_price` | Price range filter                              |
| `limit`/`offset`      | Pagination (default limit 10, max 50)             |

Response:

```json
{
  "success": true,
  "message": "Search completed successfully",
  "data": [ { ...product... } ],
  "pagination": { "offset": 0, "limit": 10, "total": 42 }
}
```

### Health check

```
GET /health  ->  { "status": "ok" }
```
