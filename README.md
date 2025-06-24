# 🚀 Go Fiber Starter Class - Complete Learning Path

## 📚 Course Overview

เรียนรู้ **Go Fiber Framework** จากระดับเริ่มต้นถึงระดับ **Production-Ready** และ **Performance Optimization**!

## 🗓️ Course Structure

### 01-beginner/ (30 นาที) - 🌱 เริ่มต้น
- สร้าง REST API พื้นฐาน
- HTTP Methods (GET, POST, PUT, DELETE)
- JSON handling
- Basic routing
- **Goal**: เข้าใจ Fiber basics

### 02-intermediate/ (45 นาที) - 🔧 ขั้นกลาง  
- MongoDB integration
- Docker setup
- Environment configuration
- Error handling
- **Goal**: สร้างระบบ CRUD ที่สมบูรณ์

### 03-advanced/ (60+ นาที) - 🎯 ขั้นสูง
ระดับ **Production-Ready** สำหรับการใช้งานจริง

1. **01-jwt-auth** (60 min) - 🔐 JWT Authentication
2. **02-rate-limit-cache** (45 min) - ⚡ Rate Limiting & Caching  
3. **03-database-advanced** (45 min) - 🗄️ Database Advanced
4. **04-microservices** (90 min) - 🏗️ Microservices Architecture
5. **05-monitoring** (60 min) - 📊 Monitoring & Metrics
6. **06-message-queue** (75 min) - 📨 Message Queue
7. **07-boilerplate** - 🏛️ Architecture Patterns
8. **08-dependency-injection** - 🔧 Dependency Injection

### 04-nightmare/ (Expert Level) - 💀 Performance & Optimization
ระดับ **Expert** สำหรับ High-Performance Applications

1. **01-zero-allocation** - 🚀 Zero Memory Allocation
2. **02-memory-pool** - 🏊 Object Pool & Memory Reuse
3. **03-goroutine-pool** - 🔄 Worker Pool Pattern  
4. **04-lock-free** - 🔓 Lock-Free Data Structures
5. **05-cpu-optimization** - ⚡ CPU Cache & SIMD
6. **06-gc-tuning** - 🗑️ Garbage Collector Tuning
7. **07-profiling** - 🔍 Advanced Profiling
8. **08-benchmarking** - 📊 Micro-benchmarks
9. **09-assembly** - 🔧 Assembly Integration
10. **10-real-world** - 🌍 Real Production Cases

## 🎯 Learning Progression

```
🌱 Beginner    →  🔧 Intermediate  →  🎯 Advanced     →  💀 Nightmare
(30 min)          (45 min)           (5+ hours)       (Expert Level)
                                                      
Basic API         MongoDB +          Production       Performance
                  Docker             Features         Optimization
```

## 📊 Skills Matrix

| Level | Duration | Skills Gained | Use Cases |
|-------|----------|---------------|-----------|
| 🌱 Beginner | 30 min | Basic REST API | Learning, Prototypes |
| 🔧 Intermediate | 45 min | Full CRUD + DB | Small Applications |  
| 🎯 Advanced | 5+ hours | Production Ready | Real Applications |
| 💀 Nightmare | Expert | High Performance | High-Scale Systems |

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker (for intermediate+)
- MongoDB (for intermediate+)

### เริ่มต้นเรียน
```bash
# 1. Clone repository
git clone <repository-url>
cd go-fiber-starter-class

# 2. เริ่มจาก beginner
cd 01-beginner/starter
go mod tidy
go run main.go

# 3. ดู complete example
cd ../complete
go run main.go
```

## 📁 Directory Structure

```
go-fiber-starter-class/
├── 01-beginner/           # 🌱 Basic REST API
│   ├── starter/           # โค้ดเริ่มต้น
│   ├── complete/          # โค้ดเสร็จสมบูรณ์
│   └── README.md          # คำอธิบายละเอียด
│
├── 02-intermediate/       # 🔧 MongoDB + Docker
│   ├── starter/
│   ├── complete/
│   ├── docker-compose.yml
│   └── README.md
│
├── 03-advanced/          # 🎯 Production Features
│   ├── 01-jwt-auth/      # Authentication
│   ├── 02-rate-limit-cache/  # Performance
│   ├── 03-database-advanced/ # Database
│   ├── 04-microservices/     # Architecture
│   ├── 05-monitoring/        # Observability
│   ├── 06-message-queue/     # Events
│   ├── 07-boilerplate/       # Design Patterns
│   └── 08-dependency-injection/ # DI Patterns
│
└── 04-nightmare/         # 💀 Performance Optimization
    ├── 01-zero-allocation/   # Memory Optimization
    ├── 02-memory-pool/       # Object Pooling
    ├── 03-goroutine-pool/    # Concurrency
    ├── 04-lock-free/         # Lock-Free Programming
    └── 10-real-world/        # Production Cases
```

## 🎓 Learning Path Recommendations

### 🎯 **For Beginners**
1. Complete `01-beginner` 
2. Move to `02-intermediate`
3. Try 1-2 topics from `03-advanced`

### 🚀 **For Experienced Developers**  
1. Quick review of `01-beginner`
2. Focus on `03-advanced` topics
3. Choose relevant patterns from `07-boilerplate`

### 💀 **For Performance Engineers**
1. Master `03-advanced` first
2. Deep dive into `04-nightmare`
3. Apply learnings to real projects

## 🔧 Development Setup

### Local Development
```bash
# Install dependencies
go mod download

# Run with hot reload (optional)
go install github.com/cosmtrek/air@latest
air
```

### Docker Development
```bash
# For intermediate+ levels
docker-compose up -d
```

## 📚 Additional Resources

### Architecture Patterns (07-boilerplate)
- **Layered Architecture** - ตรงไปตรงมา เหมาะกับโปรเจคเล็ก
- **Clean Architecture** - แยกส่วนชัดเจน เหมาะกับโปรเจคกลาง
- **Hexagonal Architecture** - ยืดหยุ่นสูง เหมาะกับระบบซับซ้อน  
- **Onion Architecture** - Enterprise-grade แยก domain ชัดเจน

### Dependency Injection (08-dependency-injection)
- **Constructor Injection** - Simple, fast, no external dependencies
- **Uber Dig** - Container-based, reflection, flexible
- **Google Wire** - Code generation, compile-time safety, fastest

### Performance Optimization (04-nightmare)
- **Zero Allocation** - ลด GC pressure 90%+
- **Memory Pool** - Reuse objects, reduce allocations
- **Lock-Free** - เพิ่ม concurrency performance 5-10x
- **Real-World Cases** - เรียนรู้จากระบบจริง

## 🏆 Success Metrics

### After Completing This Course:
✅ Build production-ready Go APIs  
✅ Handle 100K+ requests per second  
✅ Design scalable architectures  
✅ Optimize for performance  
✅ Apply best practices  

## 🤝 Contributing

เรายินดีรับ contributions! กรุณา:
1. Fork repository
2. สร้าง feature branch
3. Submit Pull Request

## 📄 License

MIT License - ใช้ได้อย่างอิสระในโปรเจคของคุณ

---

> **🎯 Remember**: "Start simple, scale smart" - เริ่มจากพื้นฐาน แล้วค่อยขยายไปสู่ระบบที่ซับซ้อนขึ้น!

**Happy Coding! 🚀** 