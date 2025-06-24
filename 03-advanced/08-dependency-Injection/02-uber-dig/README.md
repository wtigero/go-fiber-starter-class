# 🏗️ Uber Dig (Container-based DI)

## 🎯 วัตถุประสงค์
เรียนรู้ **Uber Dig** - Container-based Dependency Injection ที่จัดการ dependencies อัตโนมัติ

## 💡 หลักการ Uber Dig

```go
// ✅ GOOD: ใช้ Container จัดการ
container := dig.New()

// ลงทะเบียน constructors
container.Provide(repository.NewUserRepository)
container.Provide(service.NewUserService)
container.Provide(controller.NewUserController)

// ให้ Dig หา dependencies และสร้างให้
container.Invoke(func(uc *controller.UserController) {
    // ใช้งาน UserController ที่สร้างเสร็จแล้ว
})
```

## 🏗️ โครงสร้าง

```
02-uber-dig/
├── models/
│   └── user.go           # Entity (เหมือนเดิม)
├── repository/
│   └── user_repository.go # Data layer
├── service/
│   └── user_service.go    # Business logic  
├── controller/
│   └── user_controller.go # HTTP handlers
├── container/
│   └── container.go       # 🏗️ Dig Container setup
├── main.go               # 🔌 Container-based DI
└── go.mod                # + dig dependency
```

## 🔄 Data Flow (Uber Dig)

```
main.go → Container → Provide all constructors → Invoke → Auto-wire dependencies
         ↓
    Dig Container ทำ:
    1. วิเคราะห์ dependencies
    2. เรียก constructors ตามลำดับ
    3. ส่ง dependencies ให้อัตโนมัติ
```

## 🏗️ Container Setup

```go
// container/container.go
func BuildContainer() *dig.Container {
    container := dig.New()
    
    // ลงทะเบียนทุก constructor
    container.Provide(repository.NewUserRepository)
    container.Provide(service.NewUserService)       // รับ UserRepository
    container.Provide(controller.NewUserController) // รับ UserService
    
    return container
}
```

## ✅ ข้อดี Uber Dig

- **🤖 Auto-wiring**: ไม่ต้องจำลำดับการสร้าง dependencies
- **🔄 Flexible**: เปลี่ยน constructor ได้ง่าย
- **🎯 Lifecycle**: จัดการ singleton pattern ได้
- **🧪 Testable**: Mock dependencies ได้ง่าย
- **📈 Scalable**: เหมาะกับโปรเจคใหญ่

## ❌ ข้อเสีย

- **🐌 Slower**: ใช้ reflection (ช้ากว่า constructor)
- **🤔 Learning Curve**: ต้องเรียนรู้ Dig concepts
- **⚠️ Runtime Errors**: Error ออกมาตอน runtime
- **📦 External Dependency**: ต้องติดตั้ง uber/dig

## 🧪 การทดสอบ

```bash
# ติดตั้ง dependencies
go mod tidy

# รันแอป
go run main.go

# ทดสอบ API
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'

curl http://localhost:3000/users
```

## 📊 API Endpoints (เหมือน Constructor)

| Method | Endpoint     | Description |
|--------|-------------|-------------|
| GET    | /users      | ดึงผู้ใช้ทั้งหมด |
| POST   | /users      | สร้างผู้ใช้ใหม่ |
| GET    | /users/:id  | ดึงผู้ใช้ตาม ID |
| PUT    | /users/:id  | แก้ไขผู้ใช้ |
| DELETE | /users/:id  | ลบผู้ใช้ |

## 🔍 Advanced Features

### 1. **Named Dependencies**
```go
// หลาย implementation
container.Provide(mysql.NewUserRepository, dig.Name("mysql"))
container.Provide(mongo.NewUserRepository, dig.Name("mongo"))

// เลือกใช้
func NewUserService(repo UserRepository `name:"mysql"`) UserService {
    return &userService{repo: repo}
}
```

### 2. **Lifecycle Management**
```go
// Singleton pattern
container.Provide(NewDatabase, dig.As(new(Database)))

// Cleanup
defer container.Invoke(func(db Database) {
    db.Close()
})
```

## 🎯 เมื่อไหร่ควรใช้?

✅ **ใช้เมื่อ:**
- โปรเจคกลาง-ใหญ่ (10-50 services)
- มี complex dependencies
- ต้องการ flexibility
- ทีมงานคุ้นชิน DI patterns

❌ **ไม่ควรใช้เมื่อ:**
- โปรเจคเล็กๆ (< 5 services)
- ต้องการ performance สูงสุด
- ทีมใหม่ที่ไม่คุ้นชิน reflection

## 🔄 เปรียบเทียบกับ Constructor

| Feature | Constructor | Uber Dig |
|---------|-------------|----------|
| **Setup** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Auto-wiring** | ❌ | ✅ |
| **Learning** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Flexibility** | ⭐⭐ | ⭐⭐⭐⭐⭐ | 