# Robot Head :robot: 
### *Converstational NLP Chatbot with LED matrix visualiser for responses*

This is the next step from LED_controller.

I wanted to build something scalable so everything is being re-written in Go...

## Architecture

### Current Components:
- **Go server backend** - WebSocket coordination, hardware control
- **Go client** - Raspberry Pi hardware interface (LEDs, audio)
- **OpenAI API integration** - Direct API calls from Go server
- **TTS with elevenlabs** - Text-to-speech processing
- **Google Speech to text** - Voice recognition
- **LED matrix visualization** - Synchronized with speech responses

### Future Evolution:
- **Python AI Service** - Advanced LLM processing with LangChain
- **RAG Integration** - Vector databases, retrieval-augmented generation
- **Agentic Workflows** - Multi-step reasoning and tool calling
- **Hybrid Architecture** - Go for real-time systems + Python for AI/ML

