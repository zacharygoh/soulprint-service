# Frontend Integration Guide

## 🎯 **What You Need to Know**
Your backend is **production-ready** with local AI capabilities. This guide helps you integrate any frontend framework.

## 🔌 **Backend Connection**
- **API Base URL**: `http://localhost:8080`
- **CORS**: Pre-configured for `http://localhost:3000`
- **Authentication**: Currently hardcoded (`user123`) - ready for real auth
- **Content-Type**: `application/json`

## 🏗️ **Recommended Frontend Architecture**

### **React/Next.js Example**
```typescript
// api/client.ts
const API_BASE = 'http://localhost:8080';

export const api = {
  // Journal entries
  createEntry: (data: JournalEntry) => 
    fetch(`${API_BASE}/api/v1/entries`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    }),
  
  getEntries: () => 
    fetch(`${API_BASE}/api/v1/entries`).then(r => r.json()),
  
  // AI reflections
  generateReflection: (entryId: string, type: 'insight' | 'summary' | 'analysis') =>
    fetch(`${API_BASE}/api/v1/reflect`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ entry_id: entryId, type })
    }),
  
  getReflections: () =>
    fetch(`${API_BASE}/api/v1/reflections`).then(r => r.json())
};
```

### **Vue.js Example**
```typescript
// composables/useApi.ts
import { ref } from 'vue';

export const useApi = () => {
  const loading = ref(false);
  const API_BASE = 'http://localhost:8080';
  
  const createEntry = async (entry: any) => {
    loading.value = true;
    try {
      const response = await fetch(`${API_BASE}/api/v1/entries`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(entry)
      });
      return await response.json();
    } finally {
      loading.value = false;
    }
  };
  
  return { createEntry, loading };
};
```

## 📱 **Key Features to Implement**

### **1. Journal Entry Management**
- ✅ Create/edit/delete entries
- ✅ Rich text editor for content
- ✅ Mood selection (happy, sad, excited, etc.)
- ✅ Tag management
- ✅ Entry listing with search/filter

### **2. AI Reflection Integration**
- ✅ "Generate Reflection" button on entries
- ✅ Three reflection types: insight, summary, analysis
- ✅ Loading states (~30 seconds for local AI)
- ✅ Reflection history display
- ✅ Keyword extraction display

### **3. User Experience**
- 🔄 Loading indicators for AI generation
- 📊 Insights dashboard
- 🏷️ Tag-based filtering
- 🔍 Search functionality
- 📱 Responsive design

## 🎨 **UI/UX Recommendations**

### **Journal Entry Form**
```jsx
<form onSubmit={handleSubmit}>
  <input name="title" placeholder="Entry title..." />
  <textarea name="content" placeholder="What's on your mind?" />
  <select name="mood">
    <option value="happy">😊 Happy</option>
    <option value="sad">😢 Sad</option>
    <option value="excited">🎉 Excited</option>
    <option value="anxious">😰 Anxious</option>
  </select>
  <TagInput name="tags" />
  <button type="submit">Save Entry</button>
</form>
```

### **AI Reflection Component**
```jsx
<div className="reflection-section">
  <div className="reflection-controls">
    <button onClick={() => generateReflection('insight')}>
      💡 Generate Insight
    </button>
    <button onClick={() => generateReflection('summary')}>
      📝 Generate Summary
    </button>
    <button onClick={() => generateReflection('analysis')}>
      🔍 Generate Analysis
    </button>
  </div>
  
  {loading && <LoadingSpinner message="AI is reflecting..." />}
  {reflection && <ReflectionDisplay reflection={reflection} />}
</div>
```

## 🔄 **State Management**

### **Data Flow**
```
User Creates Entry → Save to Backend → Generate Reflection → Display Both
```

### **Recommended State Structure**
```typescript
interface AppState {
  entries: JournalEntry[];
  reflections: Reflection[];
  loading: {
    entries: boolean;
    reflections: boolean;
  };
  filters: {
    mood: string[];
    tags: string[];
    dateRange: [Date, Date];
  };
}
```

## 🔧 **Development Setup**

### **1. Start Backend**
```bash
cd soulprint-service
make run
```

### **2. Test Connection**
```bash
curl http://localhost:8080/health
```

### **3. Create Frontend Project**
```bash
# React
npx create-react-app soulprint-frontend --template typescript
cd soulprint-frontend
npm install axios # or your preferred HTTP client

# Vue
npm create vue@latest soulprint-frontend
cd soulprint-frontend
npm install axios

# Next.js
npx create-next-app@latest soulprint-frontend --typescript
```

## 🚀 **Deployment Considerations**

### **Environment Variables**
```env
# Frontend .env
REACT_APP_API_URL=http://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080
VITE_API_URL=http://localhost:8080
```

### **Production Setup**
- Backend: Consider Docker containerization
- Frontend: Standard deployment (Vercel, Netlify, etc.)
- Database: MongoDB Atlas for production
- AI: Keep local for privacy or switch to cloud for scale

## 🎯 **Quick Start Checklist**

- [ ] Backend running (`make run`)
- [ ] Health check passes (`curl http://localhost:8080/health`)
- [ ] Frontend project created
- [ ] API client implemented
- [ ] Journal entry form working
- [ ] AI reflection integration working
- [ ] Loading states implemented
- [ ] Error handling added

## 📞 **Integration Support**

### **Test These First**
1. **Create Entry**: POST `/api/v1/entries`
2. **List Entries**: GET `/api/v1/entries`
3. **Generate Reflection**: POST `/api/v1/reflect`
4. **View Reflections**: GET `/api/v1/reflections`

### **Common Issues**
- **CORS**: Already configured for localhost:3000
- **Response Time**: AI reflections take ~30 seconds
- **Authentication**: Currently hardcoded, ready for real auth
- **Error Handling**: Backend returns standard HTTP status codes

---

**Your backend is production-ready! Focus on creating a beautiful, intuitive frontend that showcases the AI reflection capabilities.**

**Last Updated**: July 5, 2025  
**Backend Status**: ✅ Production Ready  
**AI Integration**: ✅ Local Llama3 Working 