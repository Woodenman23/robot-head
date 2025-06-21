package shared

// custom type based on string
type MessageType string

// assigning constants for the values of various MessageType vars
const (
	MessageTypeUserInput	MessageType = "user_input"
	MessageTypeAIResponse	MessageType = "ai_response"
	MessageTypeStatus		MessageType = "status"
	MessageTypeError		MessageType = "error"
	MessageTypeAudio		MessageType = "audio"
)

// create a message "Class" (called struct in go)
type Message struct {
	Type 	  MessageType  `json:"type"`
	Timestamp int64        `json:"timestamp"`
	Data 	  interface{}  `json:"data"`	
}

type AudioData struct {
	Text      string  `json:"text"`
	AudioData []byte  `json:"audio_data"`
	MimeType  string  `json:"mime_type"`
}
