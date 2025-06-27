#!/bin/bash

# Robot Head Client - Voice Input/Output
# This script runs the Go client for continuous voice recording and audio playback

echo "Starting Robot Head Client..."
echo "Make sure the server is running first with: ./run-server.sh"
echo ""
echo "Voice Mode: The client will continuously record 3-second audio chunks"
echo "and send them to the server for speech recognition and AI processing."
echo ""
echo "Press Ctrl+C to quit"
echo ""

# Run the client (suppress ALSA audio warnings)
go run ./client 2>/dev/null