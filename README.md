# Robot Head :robot: 
### *Converstational NLP Chatbot with LED matrix visualiser for responses*

This is the next step from LED_controller.

I wanted to build something scalable so everything is being re-written in Go...

## Architecture

### Current Components:
- **Go WebSocket Server** - Real-time client coordination with graceful shutdown
- **Go WebSocket Client** - Interactive messaging with connection retry logic
- **OpenAI API Integration** - GPT-4o powered conversations via raw HTTP calls
- **Shared Message Protocol** - Type-safe JSON communication with full test coverage
- **Production Features** - Environment configuration, error handling, modular architecture

### In Development:
- **Text-to-Speech** - Voice output for robot responses
- **Speech-to-Text** - Voice input recognition
- **LED Matrix Visualization** - Synchronized with speech responses
- **Hardware Integration** - Raspberry Pi deployment

### Future Evolution:
- **Python AI Service** - Advanced LLM processing with LangChain
- **RAG Integration** - Vector databases, retrieval-augmented generation
- **Agentic Workflows** - Multi-step reasoning and tool calling
- **Hybrid Architecture** - Go for real-time systems + Python for AI/ML

## Features

✅ **Real-time Communication** - WebSocket client-server with retry logic  
✅ **AI-Powered Conversations** - OpenAI GPT-4o integration  
✅ **Production Ready** - Graceful shutdown, env config, comprehensive testing  
✅ **Modular Architecture** - Clean separation of concerns, swappable components  
🚧 **Voice Interface** - TTS/STT integration in progress  
🚧 **Hardware Integration** - LED visualization and Pi deployment  

## Quick Start

```bash
# Terminal 1 - Start server
cd server && go run .

# Terminal 2 - Start client
cd client && go run .

# Start chatting with your AI robot!
```

**Requirements:** OpenAI API key in `~/.api_keys/openai_key` or `OPENAI_API_KEY` environment variable.

