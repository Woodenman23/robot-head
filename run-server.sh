#!/bin/bash

# Robot Head Server - Whisper.cpp Speech-to-Text
# This script sets up the required environment variables for CGO linking
# and runs the Go server with Whisper.cpp integration

echo "Starting Robot Head Server with Whisper.cpp..."

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set CGO environment variables for Whisper.cpp
export CGO_CFLAGS="-I${PROJECT_ROOT}/whisper.cpp/include -I${PROJECT_ROOT}/whisper.cpp/ggml/include"
export CGO_LDFLAGS="-L${PROJECT_ROOT}/whisper.cpp/build/src -L${PROJECT_ROOT}/whisper.cpp/build/ggml/src -lwhisper -lggml"

# Set runtime library path
export LD_LIBRARY_PATH="${PROJECT_ROOT}/whisper.cpp/build/src:${PROJECT_ROOT}/whisper.cpp/build/ggml/src"

# Check if whisper libraries exist
if [ ! -f "${PROJECT_ROOT}/whisper.cpp/build/src/libwhisper.so" ]; then
    echo "Error: Whisper library not found. Please build whisper.cpp first:"
    echo "  cd whisper.cpp && make"
    exit 1
fi

# Check if model exists
if [ ! -f "${PROJECT_ROOT}/models/ggml-base.en.bin" ]; then
    echo "Error: Whisper model not found at ${PROJECT_ROOT}/models/ggml-base.en.bin"
    echo "Please download the model first."
    exit 1
fi

echo "Environment configured. Starting server..."
echo "CGO_CFLAGS: $CGO_CFLAGS"
echo "CGO_LDFLAGS: $CGO_LDFLAGS"
echo "LD_LIBRARY_PATH: $LD_LIBRARY_PATH"
echo ""

# Run the server
go run ./server