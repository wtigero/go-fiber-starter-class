# 🏗️ Go Fiber Architecture Boilerplates

ตัวอย่างการออกแบบ Architecture Patterns ต่างๆ สำหรับ Go Fiber

## 📁 โครงสร้าง

```
07-boilerplate/
├── 01-layered/          # Layered Architecture
├── 02-clean/            # Clean Architecture  
├── 03-hexagonal/        # Hexagonal Architecture
├── 04-onion/            # Onion Architecture
└── README.md            # ไฟล์นี้
```

## 🎯 วัตถุประสงค์

เรียนรู้และเปรียบเทียบ Architecture Patterns ที่นิยมใช้ในการพัฒนา Backend:

### **1. 🏢 Layered Architecture**
- **หลักการ**: แบ่งชั้นตาม function (Presentation → Business → Data)
- **ข้อดี**: เข้าใจง่าย, implement เร็ว
- **ข้อเสีย**: tight coupling, ทดสอบยาก

### **2. 🧹 Clean Architecture**
- **หลักการ**: Dependency inversion, Independent frameworks
- **ข้อดี**: testable, maintainable, framework independent
- **ข้อเสีย**: ซับซ้อน, setup ใช้เวลา

### **3. 🔗 Hexagonal Architecture (Ports & Adapters)**
- **หลักการ**: Core business อยู่กลาง, external เป็น adapters
- **ข้อดี**: isolate business logic, easy to swap adapters
- **ข้อเสีย**: อาจ over-engineering สำหรับ app เล็ก

### **4. 🧅 Onion Architecture**
- **หลักการ**: Dependencies point inward, business rules ไม่ depend external
- **ข้อดี**: highly testable, technology independent
- **ข้อเสีย**: learning curve สูง

## 🚀 วิธีใช้

แต่ละโฟลเดอร์มีตัวอย่าง User CRUD API เหมือนกัน แต่ใช้ pattern ต่างกัน:

```bash
# ลองรัน Layered Architecture
cd 01-layered && go run main.go

# ลองรัน Clean Architecture  
cd 02-clean && go run main.go

# ลองรัน Hexagonal Architecture
cd 03-hexagonal && go run main.go

# ลองรัน Onion Architecture
cd 04-onion && go run main.go
```

## 📚 แนะนำการเรียนรู้

1. **เริ่มจาก Layered** → เข้าใจง่ายที่สุด
2. **ไปที่ Hexagonal** → เรียนรู้ concept Ports & Adapters
3. **ศึกษา Clean** → เข้าใจ Uncle Bob's principles  
4. **จบที่ Onion** → เข้าใจ dependency direction

## 🔍 เปรียบเทียบ

| Pattern | Complexity | Testability | Maintainability | Learning Curve |
|---------|------------|-------------|-----------------|----------------|
| Layered | ⭐ | ⭐⭐ | ⭐⭐ | ⭐ |
| Clean | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Hexagonal | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| Onion | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

**เลือกตาม Project:**
- **Small/Prototype**: Layered
- **Medium/Growing**: Hexagonal  
- **Large/Enterprise**: Clean หรือ Onion

---

💡 **คำแนะนำ**: ศึกษาทุก pattern แล้วเลือกที่เหมาะกับ project และทีมของคุณ! 