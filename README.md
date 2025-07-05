# Soulprint Backend

A journaling application backend with AI-powered reflection capabilities built with Go, MongoDB, and local/cloud AI models.

## Features

- **Journal Entries**: Create, read, update, and delete personal journal entries
- **AI Reflections**: Generate AI-powered insights and reflections on journal entries
- **Local AI Support**: Use local models (Llama3, etc.) for privacy and control
- **Cloud AI Support**: Optional OpenAI integration for advanced capabilities
- **Insights Dashboard**: Get personalized insights and sentiment analysis
- **RESTful API**: Clean, well-documented API endpoints
- **MongoDB Integration**: Robust data storage with MongoDB
- **Privacy First**: Keep your data local with local model support

## Tech Stack

- **Language**: Go 1.21+
- **Database**: MongoDB
- **AI Models**: Local (Ollama/Llama3) + OpenAI GPT API
- **Router**: Gorilla Mux
- **Environment**: godotenv
- **Local AI**: Ollama, LM Studio, or custom model servers

## Prerequisites

- Go 1.21 or higher
- MongoDB (local or Atlas)
- **For Local AI**: Ollama with Llama3 model OR LM Studio OR custom model server
- **For Cloud AI**: OpenAI API key (optional)

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd soulprint-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   MONGODB_URI=mongodb://localhost:27017
   MONGODB_DATABASE=soulprint
   
   # For Local AI (Recommended)
   USE_LOCAL_MODEL=true
   LOCAL_MODEL_URL=http://localhost:11434
   LOCAL_MODEL_NAME=llama3:8b
   
   # For OpenAI (Optional)
   OPENAI_API_KEY=your_openai_api_key_here
   OPENAI_MODEL=gpt-3.5-turbo
   ```

4. **Set up AI Model**
   
   **Option A: Local AI (Recommended)**
   ```bash
   # Install Ollama
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # Start Ollama and pull Llama3
   ollama serve
   ollama pull llama3
   ```
   
   **Option B: OpenAI**
   - Get API key from https://platform.openai.com/account/api-keys
   - Set `USE_LOCAL_MODEL=false` in `.env`
   - Add your `OPENAI_API_KEY` to `.env`

5. **Start MongoDB**
   Make sure MongoDB is running locally or configure your MongoDB Atlas connection string.

6. **Run the application**
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### User Management (MVP)
- `POST /api/v1/user` - Create a user (hardcoded for MVP)

### Journal Entries
- `POST /api/v1/entries` - Create a new journal entry
- `GET /api/v1/entries` - Get all journal entries
- `GET /api/v1/entries/{id}` - Get a specific journal entry
- `PUT /api/v1/entries/{id}` - Update a journal entry
- `DELETE /api/v1/entries/{id}` - Delete a journal entry

### AI Reflections
- `POST /api/v1/reflect` - Generate AI reflection for a journal entry
- `GET /api/v1/reflections` - Get all reflections
- `GET /api/v1/entries/{id}/reflections` - Get reflections for a specific entry
- `GET /api/v1/insights` - Get personalized insights and analytics

## API Examples

### Create a Journal Entry
```bash
curl -X POST http://localhost:8080/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Day",
    "content": "Today was a great day. I accomplished so much and felt really productive.",
    "tags": ["productivity", "happiness"],
    "mood": "positive"
  }'
```

### Generate AI Reflection
```bash
# Generate insight using local Llama3
curl -X POST http://localhost:8080/api/v1/reflect \
  -H "Content-Type: application/json" \
  -d '{
    "entry_id": "64f7b123456789abcdef0123",
    "type": "insight"
  }'

# Generate summary
curl -X POST http://localhost:8080/api/v1/reflect \
  -H "Content-Type: application/json" \
  -d '{
    "entry_id": "64f7b123456789abcdef0123",
    "type": "summary"
  }'
```

### Get Insights
```bash
curl http://localhost:8080/api/v1/insights
```

## Request/Response Examples

### Journal Entry Model
```json
{
  "id": "64f7b123456789abcdef0123",
  "user_id": "user123",
  "title": "My Day",
  "content": "Today was a great day...",
  "tags": ["productivity", "happiness"],
  "mood": "positive",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Reflection Model
```json
{
  "id": "64f7b123456789abcdef0124",
  "entry_id": "64f7b123456789abcdef0123",
  "user_id": "user123",
  "content": "This entry shows a positive mindset...",
  "type": "insight",
  "keywords": ["productivity", "accomplishment"],
  "sentiment": "positive",
  "created_at": "2024-01-15T10:35:00Z"
}
```

## Database Collections

### journal_entries
- Stores user journal entries with metadata
- Indexed by user_id and created_at

### reflections
- Stores AI-generated reflections and insights
- Linked to journal entries via entry_id

## Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Database name | `soulprint` |
| `USE_LOCAL_MODEL` | Use local AI model | `false` |
| `LOCAL_MODEL_URL` | Local model server URL | `http://localhost:11434` |
| `LOCAL_MODEL_NAME` | Local model name | `llama3` |
| `OPENAI_API_KEY` | OpenAI API key (if not using local) | `""` |
| `OPENAI_MODEL` | OpenAI model to use | `gpt-3.5-turbo` |

## AI Reflection Types

- **insight**: Thoughtful reflections and perspectives (default)
- **summary**: Concise summaries of journal entries
- **analysis**: Deep analysis with growth insights

## Local vs Cloud AI

### 🏠 Local AI (Recommended)
- ✅ **Privacy**: Your data never leaves your computer
- ✅ **Cost**: No API fees
- ✅ **Speed**: Fast local processing
- ✅ **Control**: Full ownership of AI responses
- ✅ **Offline**: Works without internet
- 📱 **Models**: Llama3, Code Llama, Mistral, etc.

### ☁️ Cloud AI (OpenAI)
- ✅ **Advanced**: Latest GPT models
- ✅ **Quality**: High-quality responses
- ❌ **Privacy**: Data sent to OpenAI
- ❌ **Cost**: API usage fees
- ❌ **Internet**: Requires internet connection

## Local Model Setup

### Ollama (Recommended)
```bash
# Install
curl -fsSL https://ollama.ai/install.sh | sh

# Start service
ollama serve

# Pull models
ollama pull llama3        # 8B model
ollama pull llama3:70b    # 70B model (requires more RAM)
ollama pull codellama     # Code-focused model
```

### Alternative: LM Studio
1. Download from https://lmstudio.ai/
2. Load a model (Llama3, Mistral, etc.)
3. Start the server
4. Update `LOCAL_MODEL_URL=http://localhost:1234` in `.env`

## Development

### Project Structure
```
soulprint-backend/
├── cmd/main.go           # Application entry point
├── config/config.go      # Configuration management
├── controllers/          # HTTP handlers
│   ├── journal.go
│   └── reflection.go
├── models/journal.go     # Data models
├── routes/router.go      # Route definitions
├── services/             # Business logic
│   ├── journal_service.go
│   └── ai_service.go
├── utils/openai.go       # OpenAI integration
├── go.mod               # Go module dependencies
└── .env                 # Environment variables
```

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o soulprint-backend cmd/main.go
```

## MVP Notes

This is an MVP implementation with the following simplifications:
- User authentication is hardcoded (`user123`)
- No advanced user management
- Basic sentiment analysis
- Simple keyword extraction

## Future Enhancements

- [ ] User authentication and authorization
- [ ] Advanced sentiment analysis
- [ ] Vector embeddings for semantic search
- [ ] Email notifications
- [ ] Export functionality
- [ ] Advanced analytics dashboard
- [ ] Mobile app integration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. 