# 12. File Upload & Storage

## เวลาที่ใช้: 45 นาที

## สิ่งที่จะได้เรียนรู้

1. **Single File Upload** - อัพโหลดไฟล์ทีละไฟล์
2. **Multiple File Upload** - อัพโหลดหลายไฟล์พร้อมกัน
3. **File Validation** - ตรวจสอบประเภทและขนาดไฟล์
4. **Image Processing** - Resize, compress รูปภาพ
5. **Storage Options** - Local, S3, Cloud Storage

## Use Cases

- Profile picture upload
- Document management
- Image galleries
- File sharing systems

## โครงสร้างโปรเจค

```
12-file-upload/
├── README.md
├── starter/
│   ├── go.mod
│   └── main.go
└── complete/
    ├── go.mod
    ├── main.go
    ├── handlers/
    │   └── upload.go
    ├── storage/
    │   ├── local.go      # Local file storage
    │   └── interface.go  # Storage interface
    ├── utils/
    │   └── image.go      # Image processing
    ├── uploads/          # Upload directory
    └── static/
        └── index.html    # Upload UI
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | / | Upload UI |
| POST | /upload | Single file upload |
| POST | /upload/multiple | Multiple files upload |
| POST | /upload/image | Image upload with resize |
| GET | /files | List uploaded files |
| GET | /files/:filename | Download file |
| DELETE | /files/:filename | Delete file |

## File Validation

```go
// Allowed file types
var allowedTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/gif":  true,
    "image/webp": true,
    "application/pdf": true,
}

// Max file size: 10MB
const maxFileSize = 10 * 1024 * 1024
```

## Request Examples

### Single Upload
```bash
curl -X POST http://localhost:3000/upload \
  -F "file=@photo.jpg"
```

### Multiple Upload
```bash
curl -X POST http://localhost:3000/upload/multiple \
  -F "files=@photo1.jpg" \
  -F "files=@photo2.jpg"
```

### Image with Resize
```bash
curl -X POST http://localhost:3000/upload/image \
  -F "image=@photo.jpg" \
  -F "width=800" \
  -F "height=600"
```

## Response Format

```json
{
  "success": true,
  "file": {
    "original_name": "photo.jpg",
    "saved_name": "abc123.jpg",
    "size": 102400,
    "type": "image/jpeg",
    "url": "/files/abc123.jpg"
  }
}
```

## Security Best Practices

1. **ตรวจสอบ MIME type** - ไม่เชื่อ file extension
2. **จำกัดขนาดไฟล์** - ป้องกัน DoS
3. **สร้างชื่อไฟล์ใหม่** - ป้องกัน path traversal
4. **เก็บนอก webroot** - ป้องกัน direct execution
5. **Scan malware** - สำหรับ production
