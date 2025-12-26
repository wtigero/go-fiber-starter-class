# 11. WebSocket & Real-time Communication

## เวลาที่ใช้: 60 นาที

## สิ่งที่จะได้เรียนรู้

1. **WebSocket Basics** - การทำงานของ WebSocket กับ Fiber
2. **Real-time Chat** - สร้างระบบ chat แบบ real-time
3. **Broadcasting** - ส่งข้อความไปยังทุก clients
4. **Rooms/Channels** - แยก channel สำหรับกลุ่มต่างๆ
5. **Connection Management** - จัดการ connections

## WebSocket vs HTTP

| HTTP | WebSocket |
|------|-----------|
| Request-Response | Full-duplex |
| Client เริ่มเสมอ | ทั้งสองฝ่ายส่งได้ |
| Stateless | Stateful (connection) |
| Overhead ทุก request | Overhead ครั้งแรก |

## Use Cases

- Chat applications
- Live notifications
- Real-time dashboards
- Collaborative editing
- Live sports scores
- Stock tickers

## โครงสร้างโปรเจค

```
11-websocket/
├── README.md
├── starter/
│   ├── go.mod
│   └── main.go
└── complete/
    ├── go.mod
    ├── main.go
    ├── hub/
    │   └── hub.go       # Connection manager
    ├── handlers/
    │   └── ws.go        # WebSocket handlers
    └── static/
        └── index.html   # Chat UI
```

## การรัน

```bash
cd complete && go run main.go
# เปิด browser ไปที่ http://localhost:3000
```

## WebSocket Events

### Client → Server
```json
{"type": "join", "room": "general", "username": "john"}
{"type": "message", "content": "Hello!", "room": "general"}
{"type": "leave", "room": "general"}
{"type": "typing", "room": "general"}
```

### Server → Client
```json
{"type": "welcome", "message": "Connected to server"}
{"type": "message", "from": "john", "content": "Hello!", "timestamp": "..."}
{"type": "user_joined", "username": "jane", "room": "general"}
{"type": "user_left", "username": "john", "room": "general"}
{"type": "typing", "username": "jane"}
{"type": "online_users", "users": ["john", "jane"], "count": 2}
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | / | Chat UI |
| WS | /ws | WebSocket endpoint |
| GET | /rooms | List active rooms |
| GET | /stats | Connection statistics |

## Performance Tips

1. **Heartbeat/Ping-Pong** - ตรวจสอบ connection ยังอยู่
2. **Message Compression** - ใช้ permessage-deflate
3. **Connection Pooling** - จำกัดจำนวน connections
4. **Binary Messages** - ใช้ binary แทน JSON ถ้าเป็นไปได้
