# 10. Redis Cache & Session Management

## เวลาที่ใช้: 60 นาที

## สิ่งที่จะได้เรียนรู้

1. **Redis Basics** - การเชื่อมต่อและใช้งาน Redis กับ Go
2. **Caching Strategies** - Cache-aside, Write-through, Write-behind
3. **Rate Limiting with Redis** - Distributed rate limiting
4. **Session Management** - การจัดการ session ด้วย Redis
5. **Pub/Sub** - Real-time messaging ด้วย Redis

## ทำไมต้องใช้ Redis?

| In-Memory (Module 02) | Redis |
|----------------------|-------|
| ข้อมูลหายเมื่อ restart | ข้อมูลคงอยู่ |
| ใช้ได้แค่ 1 instance | หลาย instances แชร์ได้ |
| ไม่มี TTL ที่ดี | TTL ในตัว |
| ไม่มี Pub/Sub | มี Pub/Sub |

## โครงสร้างโปรเจค

```
10-redis-cache/
├── README.md
├── docker-compose.yml
├── starter/
│   ├── go.mod
│   └── main.go          # โครงสร้างพื้นฐาน
└── complete/
    ├── go.mod
    ├── main.go          # ระบบเต็ม
    ├── cache/
    │   └── redis.go     # Redis client wrapper
    ├── middleware/
    │   ├── cache.go     # Cache middleware
    │   ├── ratelimit.go # Rate limiter
    │   └── session.go   # Session middleware
    └── handlers/
        └── product.go   # Example handlers
```

## Caching Strategies

### 1. Cache-Aside (Lazy Loading)
```
Read: App → Check Cache → Miss → Read DB → Write Cache → Return
Write: App → Write DB → Invalidate Cache
```

### 2. Write-Through
```
Write: App → Write Cache → Write DB (sync)
Read: App → Read Cache (always hit)
```

### 3. Write-Behind (Write-Back)
```
Write: App → Write Cache → Return → Write DB (async)
Read: App → Read Cache
```

## การรัน

```bash
# Start Redis
docker-compose up -d

# Run starter
cd starter && go run main.go

# Run complete
cd complete && go run main.go
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /products | List products (cached) |
| GET | /products/:id | Get product (cached) |
| POST | /products | Create product (invalidate cache) |
| PUT | /products/:id | Update product (invalidate cache) |
| DELETE | /products/:id | Delete product (invalidate cache) |
| GET | /stats | Cache statistics |
| POST | /cache/clear | Clear all cache |

## Rate Limiting

```
# Sliding window algorithm
# 100 requests per minute per IP
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1703520000
```

## Session Example

```bash
# Login และได้ session
curl -X POST /login -d '{"email":"test@test.com","password":"123456"}'
# Response: Set-Cookie: session_id=xxx

# ใช้ session เข้าถึง protected route
curl -H "Cookie: session_id=xxx" /profile
```

## Performance Comparison

| Scenario | Without Redis | With Redis |
|----------|--------------|------------|
| Product list (1000 items) | 50ms | 2ms |
| Rate limit check | N/A (in-memory) | 0.5ms |
| Session lookup | DB query 10ms | 0.3ms |
| Cache hit rate | 0% | 95%+ |
