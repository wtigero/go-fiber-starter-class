# 🔗 Hexagonal Architecture (Ports & Adapters)

> **Ports & Adapters Pattern** - Business logic อยู่กลาง, external เชื่อมผ่าน ports & adapters

## 🎯 หลักการ

Hexagonal Architecture แยก **core business logic** ออกจาก **external concerns** โดยใช้ **Ports** (interfaces) และ **Adapters** (implementations) เพื่อให้ business logic ไม่ขึ้นกับ external systems

## 📋 โครงสร้าง

```
03-hexagonal/
├── core/              # Business Logic (Hexagon Center)
│   ├── domain/        # Entities & Business Rules
│   ├── ports/         # Interfaces (Primary & Secondary)
│   └── services/      # Application Services
├── adapters/          # External World Connections
│   ├── primary/       # Driving Adapters (HTTP, CLI, Tests)
│   └── secondary/     # Driven Adapters (Database, Email, etc)
├── main.go            # Application composition
└── README.md          # ไฟล์นี้
```

## 🔍 Hexagonal Pattern

```
                    Primary Adapters
                   (Driving/Input)
                  ┌─────────────────┐
                  │   HTTP Handler  │
                  │   CLI Commands  │
                  │   Test Cases    │
                  └─────────┬───────┘
                           │
                  ┌────────▼────────┐
                  │  Primary Port   │ ← Interface
                  │ (Use Case API)  │
                  └────────┬────────┘
          ┌─────────────────▼─────────────────┐
          │                                   │
          │         HEXAGON CORE              │
          │      (Business Logic)             │
          │                                   │
          │  ┌─────────────────────────────┐  │
          │  │     Domain Entities         │  │
          │  │     Business Rules          │  │
          │  │     Application Services    │  │
          │  └─────────────────────────────┘  │
          │                                   │
          └─────────────────┬─────────────────┘
                  ┌────────▼────────┐
                  │ Secondary Port  │ ← Interface
                  │(Repository API) │
                  └────────┬────────┘
                           │
                  ┌────────▼────────┐
                  │Secondary Adapter│
                  │   Database      │
                  │   File System   │
                  │   External API  │
                  └─────────────────┘
                   (Driven/Output)
```

## 🔄 Flow การทำงาน

```
HTTP Request → Primary Adapter → Primary Port → Core Business → Secondary Port → Secondary Adapter → Database
                     ↓                                                                  ↑
HTTP Response ← Primary Adapter ← Primary Port ← Core Business ← Secondary Port ← Secondary Adapter ← Database
```

## 📊 Components ละเอียด

### **1. 🎯 Core (Hexagon Center)**
- **Domain Entities**: Core business objects
- **Application Services**: Business logic orchestration
- **ไม่รู้จัก external systems**
- **ไม่มี dependencies ไปข้างนอก**

### **2. 🚪 Primary Ports (Driving Ports)**
- **Interfaces สำหรับ input**
- Use case interfaces
- API contracts ที่ core expose ให้ outside world
- เป็น "entry points" เข้าสู่ business logic

### **3. 🔌 Primary Adapters (Driving Adapters)**
- **Implement primary ports**
- HTTP handlers, CLI commands, Test cases
- แปลง external requests เป็น core calls

### **4. 📡 Secondary Ports (Driven Ports)**
- **Interfaces สำหรับ output**
- Repository interfaces, Email interfaces
- Contracts ที่ core ต้องการจาก external systems

### **5. 🔧 Secondary Adapters (Driven Adapters)**
- **Implement secondary ports**
- Database repositories, Email services, File systems
- แปลง core requests เป็น external system calls

## ✅ ข้อดี

- **Isolation**: Business logic แยกจาก external systems
- **Testable**: Mock adapters ได้ง่าย
- **Flexible**: เปลี่ยน adapters ได้โดยไม่กระทบ core
- **Multiple Interfaces**: สามารถมี HTTP, CLI, gRPC พร้อมกัน
- **Clear Boundaries**: Port/Adapter boundary ชัดเจน

## ❌ ข้อเสีย

- **More Complex**: มี ports และ adapters เยอะ
- **Initial Overhead**: Setup ใช้เวลา
- **Over-engineering**: อาจเยอะเกินไปสำหรับ simple apps
- **Learning Curve**: ต้องเข้าใจ concept ports & adapters

## 🔑 หลักการสำคัญ

### **Dependency Direction**
- Core ไม่ depend on adapters
- Adapters depend on ports (interfaces)
- Primary adapters call into core
- Core calls out through secondary ports

### **Ports = Interfaces**
- Primary ports: Input interfaces
- Secondary ports: Output interfaces
- Core only knows interfaces, not implementations

### **Adapters = Implementations**
- Primary adapters: External → Core
- Secondary adapters: Core → External

## 🚀 การรัน

```bash
cd 03-advanced/07-boilerplate/03-hexagonal
go run main.go
```

## 🔧 API Endpoints

```bash
# สร้าง user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# ดู users ทั้งหมด
curl http://localhost:3000/users

# ดู user ตาม ID
curl http://localhost:3000/users/1

# อัปเดต user
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'

# ลบ user
curl -X DELETE http://localhost:3000/users/1
```

## 🧪 Testing Benefits

```go
// Test core business logic with mock adapters
func TestUserService(t *testing.T) {
    mockRepo := &MockUserRepository{}
    userService := services.NewUserService(mockRepo)
    
    user, err := userService.CreateUser("John", "john@test.com")
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

## 🎓 เมื่อไหร่ควรใช้

✅ **เหมาะกับ:**
- Medium to large applications
- Multiple input interfaces (HTTP, CLI, gRPC)
- Complex integration requirements
- ต้องการ high testability
- Applications ที่ต้องเปลี่ยน external systems บ่อย

❌ **ไม่เหมาะกับ:**
- Simple CRUD applications
- Single interface applications
- Prototypes หรือ MVP
- Team ที่ไม่คุ้นเคยกับ pattern

## 🆚 เปรียบเทียบกับ Clean Architecture

| Aspect | Hexagonal | Clean |
|--------|-----------|-------|
| Focus | Ports & Adapters | Dependency Direction |
| Complexity | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| Learning | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| Flexibility | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

💡 **คำแนะนำ**: Hexagonal เหมาะกับ applications ที่ต้องการ flexibility ในการเชื่อมต่อ external systems และ multiple interfaces 