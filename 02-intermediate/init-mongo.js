// MongoDB initialization script
// สร้าง database และ collection พร้อม sample data

db = db.getSiblingDB('todoapp');

// สร้าง user สำหรับ application
db.createUser({
  user: 'todouser',
  pwd: 'todopass',
  roles: [
    {
      role: 'readWrite',
      db: 'todoapp'
    }
  ]
});

// สร้าง collection และใส่ sample data
db.todos.insertMany([
  {
    title: "เรียน Go Programming",
    done: false,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    title: "ทำโปรเจกต์ Todo API",
    done: true,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    title: "เรียน MongoDB",
    done: false,
    created_at: new Date(),
    updated_at: new Date()
  }
]);

// สร้าง index สำหรับ performance
db.todos.createIndex({ "created_at": -1 });
db.todos.createIndex({ "title": "text" });

print('Database initialization completed!'); 