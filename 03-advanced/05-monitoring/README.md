# 📊 Monitoring & Metrics - ติดตาม App ผลงาน (60 นาที)

## 📚 จุดประสงค์
เรียนรู้การติดตาม performance และ health ของ API ให้รู้ว่าทำงานเป็นอย่างไร

## 🎯 สิ่งที่จะเรียนรู้
- เก็บ Metrics แบบง่าย ๆ (ไม่ใช้ Prometheus ซับซ้อน)
- สร้าง Health Check ที่มีประโยชน์
- ทำ Request Logging อย่างมีโครงสร้าง
- แสดง Performance Dashboard แบบ REST API

## 📈 Metrics ที่จะเก็บ
1. **Request Metrics** - จำนวน requests, response time
2. **Error Metrics** - จำนวน errors, error types
3. **System Metrics** - memory usage, goroutines
4. **Business Metrics** - users created, todos completed

## 📋 API Endpoints

### Application APIs (เหมือนเดิม)
- `GET /todos` - ดู todos
- `POST /todos` - สร้าง todo
- `PUT /todos/:id` - แก้ไข todo

### Monitoring APIs (ใหม่)
- `GET /health` - Health check ละเอียด
- `GET /metrics` - Metrics ทั้งหมด
- `GET /dashboard` - Dashboard HTML
- `GET /logs` - Recent logs

## 🏃‍♂️ วิธีรัน

### รัน Application:
```bash
cd starter
go mod tidy
go run main.go
```

### ดู Monitoring:
```bash
# ดู Health
curl http://localhost:3000/health

# ดู Metrics
curl http://localhost:3000/metrics

# เปิด Dashboard ใน Browser
open http://localhost:3000/dashboard
```

## 🧪 ทดสอบ Monitoring

### 1. สร้าง Load เพื่อดู Metrics
```bash
# สร้าง requests เยอะ ๆ
for i in {1..50}; do
  curl -X POST http://localhost:3000/todos \
    -H "Content-Type: application/json" \
    -d '{"title":"Todo '${i}'","done":false}'
  
  curl http://localhost:3000/todos > /dev/null
  
  # สร้าง error บ้าง
  curl http://localhost:3000/todos/999 > /dev/null 2>&1
done
```

### 2. ดู Health Check
```bash
curl http://localhost:3000/health
```
**ผลลัพธ์:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "uptime": "2h30m15s",
  "version": "1.0.0",
  "database": {
    "status": "connected",
    "response_time": "2ms"
  },
  "system": {
    "memory_usage": "45MB",
    "goroutines": 12,
    "cpu_usage": "15%"
  }
}
```

### 3. ดู Metrics
```bash
curl http://localhost:3000/metrics
```
**ผลลัพธ์:**
```json
{
  "requests": {
    "total": 152,
    "success": 142,
    "errors": 10,
    "avg_response_time": "25ms"
  },
  "endpoints": {
    "/todos": {
      "GET": {"count": 50, "avg_time": "15ms"},
      "POST": {"count": 50, "avg_time": "35ms"}
    }
  },
  "errors": {
    "404_not_found": 8,
    "500_internal_error": 2
  },
  "business": {
    "todos_created": 50,
    "todos_completed": 12
  }
}
```

## 🔍 สิ่งสำคัญที่เรียนรู้

### 1. Metrics Collection
```go
type Metrics struct {
    RequestCount    int64         `json:"request_count"`
    ErrorCount      int64         `json:"error_count"`
    ResponseTimes   []time.Duration
    EndpointStats   map[string]EndpointMetric
    mutex           sync.RWMutex
}

func (m *Metrics) RecordRequest(endpoint string, duration time.Duration, statusCode int) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.RequestCount++
    m.ResponseTimes = append(m.ResponseTimes, duration)
    
    if statusCode >= 400 {
        m.ErrorCount++
    }
}
```

### 2. Health Check
```go
type HealthStatus struct {
    Status     string                 `json:"status"`
    Timestamp  time.Time             `json:"timestamp"`
    Uptime     string                `json:"uptime"`
    System     SystemHealth          `json:"system"`
    Database   DatabaseHealth        `json:"database"`
}

func checkHealth() HealthStatus {
    return HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Uptime:    time.Since(startTime).String(),
        System:    getSystemHealth(),
        Database:  getDatabaseHealth(),
    }
}
```

### 3. Middleware สำหรับ Monitoring
```go
func monitoringMiddleware(metrics *Metrics) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // ดำเนินการ request
        err := c.Next()
        
        // เก็บ metrics
        duration := time.Since(start)
        endpoint := c.Route().Path
        statusCode := c.Response().StatusCode()
        
        metrics.RecordRequest(endpoint, duration, statusCode)
        
        return err
    }
}
```

## 📊 Dashboard HTML (แบบง่าย)
```html
<!DOCTYPE html>
<html>
<head>
    <title>API Monitoring Dashboard</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>
    <h1>📊 API Monitoring</h1>
    
    <div class="stats">
        <div class="stat-card">
            <h3>Total Requests</h3>
            <span id="total-requests">Loading...</span>
        </div>
        
        <div class="stat-card">
            <h3>Error Rate</h3>
            <span id="error-rate">Loading...</span>
        </div>
        
        <div class="stat-card">
            <h3>Avg Response Time</h3>
            <span id="avg-response">Loading...</span>
        </div>
    </div>
    
    <canvas id="responseTimeChart"></canvas>
    
    <script>
        // อัปเดตข้อมูลทุก 5 วินาที
        setInterval(updateDashboard, 5000);
    </script>
</body>
</html>
```

## 📝 ใน starter/ จะมี:
- [ ] TODO: สร้าง Metrics struct
- [ ] TODO: สร้าง Health Check
- [ ] TODO: เพิ่ม Monitoring middleware
- [ ] TODO: สร้าง Dashboard endpoint
- [ ] TODO: เพิ่ม Structured logging

## ✅ ใน complete/ จะมี:
- ✅ Metrics collection สมบูรณ์
- ✅ Health check ละเอียด
- ✅ Dashboard แบบ real-time
- ✅ Structured logging
- ✅ System monitoring (memory, CPU)

## 💡 ประโยชน์ในชีวิตจริง

### การ Debug:
- รู้ว่า endpoint ไหนช้า
- รู้ว่า error เกิดที่ไหนบ่อย
- ติดตาม memory leaks

### การ Scale:
- รู้ว่าต้อง scale เมื่อไหร่
- รู้ว่า bottleneck อยู่ตรงไหน
- วางแผน capacity

### Business Intelligence:
- รู้ว่า feature ไหนถูกใช้บ่อย
- รู้ pattern การใช้งาน

## ⏭️ ขั้นต่อไป
หากต้องการ Production-grade:
- ใช้ Prometheus + Grafana
- เพิ่ม Alerting rules
- ทำ Distributed tracing
- เก็บ logs ใน ELK Stack

---
**เวลาเรียน:** 60 นาที | **ความยาก:** ⭐⭐⭐☆☆  
**เหมาะสำหรับ:** ทุกคนที่ทำ API ใช้งานจริง 