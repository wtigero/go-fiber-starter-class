# 🏗️ Microservices - แยก Services ขนาดเล็ก (90 นาที)

## 📚 จุดประสงค์
เรียนรู้การแยก API เดียวให้เป็นหลาย ๆ services เล็ก ๆ ที่สื่อสารกันได้

## 🎯 สิ่งที่จะเรียนรู้
- แยก User Service และ Todo Service
- Service-to-service communication
- API Gateway pattern
- Service discovery แบบง่าย
- Health checks ระหว่าง services

## 🏢 Architecture ที่จะสร้าง

```
Client --> API Gateway --> User Service
                      --> Todo Service
```

### Services ที่จะมี:
1. **API Gateway** (port 3000) - จุดเข้าหลัก
2. **User Service** (port 3001) - จัดการ users
3. **Todo Service** (port 3002) - จัดการ todos

## 📋 API Endpoints

### API Gateway (localhost:3000)
- `GET /health` - Health check ทุก services
- `POST /users` - สร้าง user (proxy ไป User Service)
- `GET /users/:id` - ดู user (proxy ไป User Service)
- `GET /todos` - ดู todos (proxy ไป Todo Service)
- `POST /todos` - สร้าง todo (proxy ไป Todo Service)

### User Service (localhost:3001)
- `GET /health` - Health check
- `POST /users` - สร้าง user
- `GET /users/:id` - ดู user

### Todo Service (localhost:3002)
- `GET /health` - Health check
- `GET /todos` - ดู todos
- `POST /todos` - สร้าง todo
- `GET /todos/user/:user_id` - ดู todos ของ user

## 🏃‍♂️ วิธีรัน

### เริ่มทั้ง 3 services พร้อมกัน:

#### Terminal 1: User Service
```bash
cd user-service
go mod tidy
go run main.go
# รันที่ port 3001
```

#### Terminal 2: Todo Service
```bash
cd todo-service
go mod tidy
go run main.go
# รันที่ port 3002
```

#### Terminal 3: API Gateway
```bash
cd api-gateway
go mod tidy
go run main.go
# รันที่ port 3000
```

## 🧪 ทดสอบ Microservices

### 1. ตรวจสอบ Health ทุก services
```bash
curl http://localhost:3000/health
```
**ผลลัพธ์:**
```json
{
  "gateway": "healthy",
  "services": {
    "user_service": "healthy",
    "todo_service": "healthy"
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 2. สร้าง User ผ่าน Gateway
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "นาย Microservice",
    "email": "micro@example.com"
  }'
```

### 3. สร้าง Todo ผ่าน Gateway
```bash
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "ทดสอบ Microservices",
    "user_id": 1
  }'
```

### 4. ทดสอบเรียกตรง ๆ User Service
```bash
curl http://localhost:3001/users/1
```

## 🔍 สิ่งสำคัญที่เรียนรู้

### 1. Service Discovery (แบบง่าย)
```go
type ServiceRegistry struct {
    services map[string]string
}

func (sr *ServiceRegistry) Register(name, url string) {
    sr.services[name] = url
}

func (sr *ServiceRegistry) GetService(name string) (string, bool) {
    url, exists := sr.services[name]
    return url, exists
}
```

### 2. HTTP Client สำหรับเรียก services
```go
func callService(serviceURL, endpoint string, data interface{}) (*http.Response, error) {
    jsonData, _ := json.Marshal(data)
    
    resp, err := http.Post(serviceURL+endpoint, "application/json", 
        bytes.NewBuffer(jsonData))
    return resp, err
}
```

### 3. Health Check System
```go
func checkServiceHealth(serviceURL string) bool {
    resp, err := http.Get(serviceURL + "/health")
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    return resp.StatusCode == 200
}
```

## 📁 โครงสร้างโฟลเดอร์
```
04-microservices/
├── starter/
│   ├── api-gateway/
│   │   ├── main.go       # TODO: สร้าง gateway
│   │   └── go.mod
│   ├── user-service/
│   │   ├── main.go       # TODO: สร้าง user service
│   │   └── go.mod
│   └── todo-service/
│       ├── main.go       # TODO: สร้าง todo service
│       └── go.mod
├── complete/
│   ├── api-gateway/      # เฉลยสมบูรณ์
│   ├── user-service/
│   └── todo-service/
└── docker-compose.yml    # รันทั้งหมดใน Docker
```

## 📝 ใน starter/ จะมี:
- [ ] TODO: สร้าง User Service พื้นฐาน
- [ ] TODO: สร้าง Todo Service พื้นฐาน
- [ ] TODO: สร้าง API Gateway ที่ proxy requests
- [ ] TODO: เพิ่ม Health Check system
- [ ] TODO: เพิ่ม Service Discovery

## ✅ ใน complete/ จะมี:
- ✅ 3 services ทำงานอิสระ
- ✅ API Gateway ที่ proxy อย่างฉลาด
- ✅ Health monitoring
- ✅ Error handling ระหว่าง services
- ✅ Docker Compose สำหรับรันง่าย ๆ

## 💡 ข้อดี/ข้อเสีย

### ข้อดี:
- แต่ละ service มีความรับผิดชอบชัด
- Deploy แยกได้อิสระ
- Scale แต่ละ service ตามความต้องการ
- ทีมพัฒนาแยกกันได้

### ข้อเสีย:
- ซับซ้อนกว่า monolithic
- ต้องจัดการ network latency
- การ debug ยากขึ้น
- ต้องมี monitoring ดี ๆ

## ⏭️ ขั้นต่อไป
- เพิ่ม Load Balancer
- ใช้ Message Queue แทน HTTP calls
- เพิ่ม Circuit Breaker pattern
- ใช้ Service Mesh เช่น Istio

---
**เวลาเรียน:** 90 นาที | **ความยาก:** ⭐⭐⭐⭐☆  
**เหมาะสำหรับ:** คนที่ต้องการทำ system ขนาดใหญ่ 