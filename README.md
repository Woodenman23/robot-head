# Robot Head :robot: 
### *Converstational NLP Chatbot with LED matrix visualiser for responses*

This is the next step from LED_controller.

I wanted to build something scalable so everything is being re-written in Go...

## Architecture

### System Overview

The robot head operates as a distributed real-time voice conversation system with clean separation between hardware interface and AI processing:

**Server (AI Processing Hub):**
- WebSocket server continuously listens for incoming audio messages
- Whisper.cpp handles speech-to-text transcription with offline processing
- OpenAI GPT-4o generates conversational responses to detected speech
- ElevenLabs provides text-to-speech conversion for robot voice output
- Only processes audio containing actual speech (filters silence automatically)

**Client (Voice Interface):**
- Continuously records 1-second audio chunks via microphone
- Streams all audio to server without local processing
- Receives and plays TTS audio responses from robot
- Maintains real-time conversation loop with minimal latency

**Communication Flow:**
```
Client: Record Audio → Send to Server → Play Response → Repeat
Server: Receive Audio → STT → LLM → TTS → Send Response → Listen
```

This architecture enables:
- **Real-time Conversations** - Sub-second response times for natural dialogue
- **Scalable Processing** - Server can handle multiple robot clients simultaneously  
- **Offline Capability** - Whisper.cpp runs locally without internet dependency
- **Clean Interfaces** - Modular design allows easy component swapping

### Current Components:
- **Go WebSocket Server** - Real-time audio processing and AI coordination
- **Go WebSocket Client** - Continuous voice recording and audio playback
- **Whisper.cpp Integration** - Local speech-to-text with Go bindings
- **OpenAI API Integration** - GPT-4o powered conversations via HTTP
- **ElevenLabs TTS** - High-quality voice synthesis for robot responses
- **Shared Message Protocol** - Type-safe JSON communication with full test coverage

### Future Evolution:
- **LED Matrix Visualization** - Synchronized with speech responses
- **Python AI Service** - Advanced LLM processing with LangChain
- **RAG Integration** - Vector databases, retrieval-augmented generation
- **Multi-Robot Support** - Cloud server coordinating multiple robot heads

## Features

✅ **Real-time Voice Conversations** - Continuous speech recognition and synthesis  
✅ **Offline Speech Processing** - Whisper.cpp runs locally without internet  
✅ **AI-Powered Responses** - OpenAI GPT-4o integration with conversation context  
✅ **High-Quality Voice Output** - ElevenLabs TTS for natural robot speech  
✅ **Production Ready** - Graceful shutdown, error handling, modular architecture  
✅ **Silence Detection** - Automatically filters background noise and empty audio  
🚧 **LED Matrix Visualization** - Synchronized with speech responses  
🚧 **Multi-Robot Support** - Cloud deployment for multiple robot heads  

## Quick Start

```bash
# Install system dependencies
sudo apt update && sudo apt install -y cmake build-essential portaudio19-dev

# Build Whisper.cpp (first time only)
cd whisper.cpp && make && cd ..

# Terminal 1 - Start server
./run-server.sh

# Terminal 2 - Start client  
./run-client.sh

# Start talking to your robot!
```

**Requirements:** 
- OpenAI API key in `~/.api_keys/openai_key` or `OPENAI_API_KEY` environment variable
- ElevenLabs API key in `~/.api_keys/elevenlabs_key` or `ELEVENLABS_API_KEY` environment variable

