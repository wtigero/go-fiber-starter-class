# 🚦 Rate Limiting & Caching - แบบง่าย ๆ (45 นาที)

## 📚 จุดประสงค์
เรียนรู้การควบคุมการใช้ API และการเก็บ cache เพื่อเพิ่มประสิทธิภาพ

## 🎯 สิ่งที่จะเรียนรู้
- จำกัดจำนวนครั้งที่เรียก API (Rate Limiting)
- เก็บข้อมูลใน Memory Cache
- แสดงข้อมูลสถิติการใช้งาน
- เข้าใจวิธีป้องกัน API abuse

## 📋 API Endpoints
- `GET /todos` - ดู todos (มี rate limit 10 ครั้ง/นาที)
- `GET /stats` - ดูสถิติการใช้งาน
- `GET /cache/clear` - ล้าง cache

## 🔧 การทำงาน
1. **Rate Limiting**: จำกัด IP ละ 10 requests ต่อนาที
2. **Memory Cache**: เก็บผลลัพธ์ API ไว้ 30 วินาที
3. **Stats**: นับจำนวนการเรียกใช้งาน

## 🏃‍♂️ วิธีรัน

### สำหรับนักเรียน:
```bash
cd starter
go mod tidy
go run main.go
```

### สำหรับดูเฉลย:
```bash
cd complete
go mod tidy
go run main.go
```

## 🧪 ทดสอบ Rate Limiting

### 1. เรียก API ปกติ (ครั้งที่ 1-10)
```bash
curl http://localhost:3000/todos
```
**ผลลัพธ์:** ปกติ

### 2. เรียก API เกิน limit (ครั้งที่ 11+)
```bash
# เรียกซ้ำ ๆ เร็ว ๆ
for i in {1..15}; do curl http://localhost:3000/todos; echo; done
```
**ผลลัพธ์ครั้งที่ 11:**
```json
{
  "error": "Rate limit exceeded. Try again later.",
  "retry_after": "45s"
}
```

### 3. ดูสถิติ
```bash
curl http://localhost:3000/stats
```
**ผลลัพธ์:**
```json
{
  "total_requests": 15,
  "cache_hits": 8,
  "cache_misses": 2,
  "rate_limited": 5
}
```

## 🔍 สิ่งสำคัญที่เรียนรู้

### 1. Rate Limiting (ใน Memory)
```go
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
}

func (r *RateLimiter) Allow(ip string) bool {
    // ตรวจสอบจำนวน requests ใน 1 นาทีที่ผ่านมา
}
```

### 2. Memory Cache
```go
type Cache struct {
    data   map[string]CacheItem
    mutex  sync.RWMutex
}

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}
```

### 3. Middleware การใช้งาน
```go
app.Use("/todos", rateLimitMiddleware())
app.Use("/todos", cacheMiddleware())
```

## 📝 ใน starter/ จะมี:
- [ ] TODO: สร้าง RateLimiter struct
- [ ] TODO: สร้าง Cache struct  
- [ ] TODO: สร้าง middleware functions
- [ ] TODO: เพิ่ม statistics tracking

## ✅ ใน complete/ จะมี:
- ✅ Rate limiter แบบ in-memory
- ✅ Cache system พื้นฐาน
- ✅ Stats tracking
- ✅ Clean error messages
- ✅ Thread-safe operations

## 💡 ข้อดี/ข้อจำกัด

### ข้อดี:
- เรียนรู้ง่าย ไม่ต้องติดตั้ง Redis
- เหมาะสำหรับ app เล็ก ๆ
- ทำความเข้าใจหลักการ

### ข้อจำกัด:
- หาย cache เมื่อ restart app
- ไม่เหมาะกับ multiple servers
- ใช้ memory บน server

## ⏭️ ขั้นต่อไป
หากต้องการใช้งานจริง ควรอัปเกรดเป็น:
- Redis สำหรับ cache และ rate limiting
- Database สำหรับเก็บ statistics
- Distributed rate limiting

---
**เวลาเรียน:** 45 นาที | **ความยาก:** ⭐⭐⭐☆☆  
**เหมาะสำหรับ:** เข้าใจหลักการก่อนใช้ Redis 