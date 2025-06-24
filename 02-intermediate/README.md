# ⚙️ Intermediate Level - Go Fiber + Database (45 นาที)

## 📚 จุดประสงค์
เรียนรู้การสร้าง REST API ที่สมบูรณ์ด้วย Go Fiber + MongoDB พร้อม Docker containerization

## 🎯 สิ่งที่จะเรียนรู้
- CRUD operations สมบูรณ์ (PUT, DELETE)
- MongoDB integration ด้วย mongo-go-driver
- Docker และ Docker Compose
- Custom middleware (API Key validation)
- Advanced error handling และ validation
- Environment variables

## 📋 API Endpoints ที่จะสร้าง
- `GET /health` - Health check endpoint
- `GET /todos` - ดูรายการ Todo ทั้งหมด (พร้อม pagination)
- `GET /todos/:id` - ดูรายการ Todo เดียว
- `POST /todos` - เพิ่มรายการ Todo ใหม่
- `PUT /todos/:id` - แก้ไขรายการ Todo
- `DELETE /todos/:id` - ลบรายการ Todo

## 📊 Data Structure
```go
type Todo struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Title     string             `json:"title" bson:"title" validate:"required,min=1"`
    Done      bool               `json:"done" bson:"done"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
```

## 🛠️ เทคโนโลยีที่ใช้
- **Go Fiber v2** - Web framework
- **MongoDB** - NoSQL database
- **mongo-go-driver** - Official MongoDB driver
- **Docker & Docker Compose** - Containerization
- **go-playground/validator** - Data validation

## 🏃‍♂️ วิธีรัน

### ขั้นตอนที่ 1: เตรียม Docker Environment
```bash
# สร้าง MongoDB container
docker-compose up -d

# ตรวจสอบ MongoDB ทำงาน
docker-compose ps
```

### ขั้นตอนที่ 2: รัน Application

#### สำหรับนักเรียน (starter):
```bash
cd starter
go mod tidy
go run main.go
```

#### สำหรับผู้สอน (complete):
```bash
cd complete
go mod tidy
go run main.go
```

Server จะรันที่ `http://localhost:3000`
MongoDB จะรันที่ `mongodb://localhost:27017`

## 🧪 ทดสอบ API

### 1. Health Check
```bash
curl http://localhost:3000/health
```
**ผลลัพธ์:**
```json
{
  "status": "OK",
  "database": "connected",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 2. ดูรายการ todos ทั้งหมด (พร้อม pagination)
```bash
curl "http://localhost:3000/todos?page=1&limit=10" \
  -H "X-API-Key: your-secret-key"
```

### 3. เพิ่ม todo ใหม่
```bash
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-key" \
  -d '{
    "title": "เรียน MongoDB",
    "done": false
  }'
```

### 4. แก้ไข todo
```bash
curl -X PUT http://localhost:3000/todos/65a1b2c3d4e5f6789abcdef0 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-key" \
  -d '{
    "title": "เรียน MongoDB (อัพเดต)",
    "done": true
  }'
```

### 5. ลบ todo
```bash
curl -X DELETE http://localhost:3000/todos/65a1b2c3d4e5f6789abcdef0 \
  -H "X-API-Key: your-secret-key"
```

## 🐳 Docker Configuration

### docker-compose.yml
```yaml
version: '3.8'
services:
  mongodb:
    image: mongo:7
    container_name: todo-mongodb
    restart: always
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: password123
      MONGO_INITDB_DATABASE: todoapp
    volumes:
      - mongodb_data:/data/db
      - ./init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js:ro

volumes:
  mongodb_data:
```

## 🔐 Environment Variables
```bash
# .env file
MONGODB_URI=mongodb://admin:password123@localhost:27017/todoapp?authSource=admin
API_SECRET_KEY=your-secret-key
PORT=3000
ENVIRONMENT=development
```

## 📝 สิ่งที่ควรสังเกต

### ใน starter/:
- มี Docker setup พร้อมใช้
- มี basic structure สำหรับ MongoDB connection
- มี TODO comments ให้เติม CRUD operations

### ใน complete/:
- มี CRUD operations สมบูรณ์
- มี custom middleware สำหรับ API Key validation
- มี proper error handling
- มี pagination และ sorting
- มี data validation ด้วย struct tags

## 🔍 Key Learning Points

1. **MongoDB Connection:**
   ```go
   client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
   ```

2. **Custom Middleware:**
   ```go
   func APIKeyMiddleware() fiber.Handler {
       return func(c *fiber.Ctx) error {
           apiKey := c.Get("X-API-Key")
           if apiKey != expectedKey {
               return c.Status(401).JSON(fiber.Map{"error": "Invalid API Key"})
           }
           return c.Next()
       }
   }
   ```

3. **MongoDB Operations:**
   ```go
   // Insert
   result, err := collection.InsertOne(ctx, todo)
   
   // Update
   filter := bson.M{"_id": objectId}
   update := bson.M{"$set": updatedTodo}
   
   // Delete
   collection.DeleteOne(ctx, bson.M{"_id": objectId})
   ```

4. **Data Validation:**
   ```go
   validate := validator.New()
   if err := validate.Struct(todo); err != nil {
       return c.Status(400).JSON(fiber.Map{"error": err.Error()})
   }
   ```

## 🐛 Troubleshooting

### MongoDB Connection Issues:
```bash
# ตรวจสอบ MongoDB container
docker logs todo-mongodb

# เชื่อมต่อ MongoDB shell
docker exec -it todo-mongodb mongosh --username admin --password password123
```

### API Testing:
```bash
# ทดสอบโดยไม่ใส่ API Key (ควรได้ 401)
curl http://localhost:3000/todos

# ทดสอบด้วย API Key ผิด (ควรได้ 401) 
curl -H "X-API-Key: wrong-key" http://localhost:3000/todos
```

## ⏭️ ถัดไปคือ...
หลังจากเสร็จระดับ Intermediate แล้ว ให้ไปต่อที่:
```bash
cd ../advanced
cat README.md
```

---
**เวลาเรียน:** 45 นาที | **ความยาก:** ⭐⭐⭐☆☆ 