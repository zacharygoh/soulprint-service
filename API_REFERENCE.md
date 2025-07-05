# Soulprint Backend API Reference

## 🌐 **Base URL**: `http://localhost:8080`

## 🔗 **Endpoints**

### 🏥 **Health & Info**
```http
GET /health
```
**Response**: `{"status": "healthy", "timestamp": "2025-07-05T12:00:00Z"}`

---

### 📝 **Journal Entries**

#### Create Entry
```http
POST /api/v1/entries
Content-Type: application/json

{
  "title": "My Day",
  "content": "Today was amazing...",
  "mood": "happy",
  "tags": ["work", "achievement"]
}
```
**Response**: `{"success": true, "data": {...}}`

#### Get All Entries
```http
GET /api/v1/entries
```
**Response**: `{"success": true, "data": [...]}`

#### Get Single Entry
```http
GET /api/v1/entries/{id}
```
**Response**: `{"success": true, "data": {...}}`

#### Update Entry
```http
PUT /api/v1/entries/{id}
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content...",
  "mood": "excited",
  "tags": ["updated"]
}
```

#### Delete Entry
```http
DELETE /api/v1/entries/{id}
```
**Response**: `{"success": true, "message": "Entry deleted successfully"}`

---

### 🧠 **AI Reflections**

#### Generate Reflection
```http
POST /api/v1/reflect
Content-Type: application/json

{
  "entry_id": "6868b0a2065bc4e96b88a833",
  "type": "insight"
}
```
**Types**: `insight`, `summary`, `analysis`  
**Response**: `{"success": true, "data": {...}}`

#### Get All Reflections
```http
GET /api/v1/reflections
```
**Response**: `{"success": true, "data": [...]}`

#### Get Entry Reflections
```http
GET /api/v1/entries/{id}/reflections
```
**Response**: `{"success": true, "data": [...]}`

#### Get Insights
```http
GET /api/v1/insights
```
**Response**: `{"success": true, "data": {...}}`

---

### 👤 **User Management (MVP)**

#### Create User
```http
POST /api/v1/user
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```
**Response**: `{"success": true, "data": {...}, "message": "User created successfully (MVP hardcoded)"}`

---

## 📊 **Data Models**

### Journal Entry
```json
{
  "id": "6868b0a2065bc4e96b88a833",
  "user_id": "user123",
  "title": "My Day",
  "content": "Today was amazing...",
  "tags": ["work", "achievement"],
  "mood": "happy",
  "created_at": "2025-07-05T12:00:00Z",
  "updated_at": "2025-07-05T12:00:00Z"
}
```

### Reflection
```json
{
  "id": "6868b0c1065bc4e96b88a834",
  "entry_id": "6868b0a2065bc4e96b88a833",
  "user_id": "user123",
  "content": "This reflection shows...",
  "type": "insight",
  "keywords": ["excitement", "growth"],
  "sentiment": "positive",
  "created_at": "2025-07-05T12:00:00Z"
}
```

### User (MVP)
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

---

## 🔧 **Configuration**

### Environment Variables
- `PORT=8080`
- `MONGODB_URI=mongodb://localhost:27017`
- `MONGODB_DATABASE=soulprint`
- `USE_LOCAL_MODEL=true`
- `LOCAL_MODEL_URL=http://localhost:11434`
- `LOCAL_MODEL_NAME=llama3:8b`

### CORS Settings
- **Enabled for**: `http://localhost:3000` (typical React dev server)
- **Methods**: GET, POST, PUT, DELETE, OPTIONS
- **Headers**: Content-Type, Authorization

---

## 🚀 **Quick Start for Frontend**

### 1. Test Connection
```bash
curl http://localhost:8080/health
```

### 2. Create a Journal Entry
```bash
curl -X POST http://localhost:8080/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "content": "Hello world", "mood": "happy"}'
```

### 3. Generate AI Reflection
```bash
curl -X POST http://localhost:8080/api/v1/reflect \
  -H "Content-Type: application/json" \
  -d '{"entry_id": "ENTRY_ID_FROM_STEP_2", "type": "insight"}'
```

---

## 🐛 **Common Issues**

### Port 8080 in Use
```bash
lsof -i :8080
kill -9 PID
```

### MongoDB Not Running
```bash
brew services start mongodb-community
```

### Ollama Not Running
```bash
ollama serve
```

---

**Last Updated**: July 5, 2025  
**Backend Status**: Production Ready  
**AI Integration**: Local Llama3 Working 