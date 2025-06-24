# 🧹 Clean Architecture

> **Uncle Bob's Clean Architecture** - Business rules ไม่ depend on frameworks, UI, database

## 🎯 หลักการ

Clean Architecture มุ่งเน้นการแยก business logic ออกจาก external concerns (framework, database, UI) โดยใช้ **Dependency Inversion Principle**

## 📋 โครงสร้าง

```
02-clean/
├── domain/           # Enterprise Business Rules
│   ├── entities/     # Core business entities
│   └── repositories/ # Repository interfaces
├── usecases/         # Application Business Rules
│   ├── interfaces/   # Use case interfaces
│   └── user/         # User use cases
├── infrastructure/   # Frameworks & Drivers
│   ├── database/     # Database implementations
│   └── web/          # HTTP handlers (Fiber)
├── main.go           # Main composition root
└── README.md         # ไฟล์นี้
```

## 🎯 Clean Architecture Layers

```
┌─────────────────────────────────────────────────┐
│                 Frameworks                      │ ← Web, DB, UI
├─────────────────────────────────────────────────┤
│            Interface Adapters                   │ ← Controllers, Presenters
├─────────────────────────────────────────────────┤
│              Use Cases                          │ ← Application Business Rules
├─────────────────────────────────────────────────┤
│               Entities                          │ ← Enterprise Business Rules
└─────────────────────────────────────────────────┘
```

## 🔄 Flow การทำงาน

```
HTTP Request → Controller → Use Case ←→ Repository Interface
                    ↓                        ↑
HTTP Response ← Controller ← Entity    Repository Implementation
```

## 📊 Layer ละเอียด

### **1. 🏛️ Entities (Domain)**
- **Core business objects**
- Independent of everything
- เก็บ business rules ที่สำคัญที่สุด
- ไม่ depend on anything

### **2. 🎯 Use Cases (Application Business Rules)**
- **Application-specific business rules**
- Orchestrate flow of data to/from entities
- ใช้ repository interfaces (not implementations)
- Independent of frameworks และ UI

### **3. 🔌 Interface Adapters**
- **Convert data formats**
- Controllers, Presenters, Gateways
- แปลงข้อมูลระหว่าง use cases และ external world

### **4. 🏗️ Frameworks & Drivers**
- **External tools**
- Database, Web frameworks, UI
- อยู่ชั้นนอกสุด สามารถเปลี่ยนได้ง่าย

## ✅ ข้อดี

- **Testable**: Business logic แยกออกจาก frameworks
- **Framework Independent**: สามารถเปลี่ยน framework ได้
- **Database Independent**: สามารถเปลี่ยน database ได้
- **UI Independent**: สามารถมีหลาย UI
- **Maintainable**: Dependencies ไปทิศทางเดียว (inward)

## ❌ ข้อเสีย

- **Complex**: มี interfaces และ layers เยอะ
- **Over-engineering**: อาจเยอะเกินไปสำหรับ simple apps
- **Learning Curve**: ต้องเข้าใจ dependency inversion
- **Initial Setup**: ใช้เวลาในการ setup เยอะ

## 🔑 หลักการสำคัญ

### **Dependency Rule**
> Dependencies point **INWARD** เท่านั้น
- Outer layers depend on inner layers
- Inner layers ไม่รู้จัก outer layers

### **Entities**
- ไม่ depend on anything
- Pure business logic

### **Use Cases**
- depend on Entities และ Repository interfaces
- ไม่ depend on implementations

### **Interface Adapters**
- depend on Use Cases
- implement Repository interfaces

## 🚀 การรัน

```bash
cd 00-boilerplate/02-clean
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
// Test Use Case without any external dependencies
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    useCase := user.NewCreateUserUseCase(mockRepo)
    
    result, err := useCase.Execute(CreateUserRequest{
        Name: "John", Email: "john@test.com",
    })
    
    assert.NoError(t, err)
    assert.Equal(t, "John", result.Name)
}
```

## 🎓 เมื่อไหร่ควรใช้

✅ **เหมาะกับ:**
- Large applications
- Complex business rules
- Long-term projects
- Team ที่ต้องการ maintainability สูง
- Applications ที่ต้องการ high testability

❌ **ไม่เหมาะกับ:**
- Simple CRUD applications
- Prototypes หรือ MVP
- Team เล็กที่ต้องการ delivery เร็ว
- Projects ที่ business rules ไม่ซับซ้อน

---

💡 **คำแนะนำ**: Clean Architecture เหมาะกับ enterprise applications ที่ต้องการ maintainability และ testability สูง แต่อาจเป็น overkill สำหรับ simple applications 