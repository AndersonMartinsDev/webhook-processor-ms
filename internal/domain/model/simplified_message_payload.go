package model

type SimplifiedMessagePayload struct {
	OriginMessage string `json:"origin_message"`
	SessionKey    string `json:"session_key"`
	AgentID       uint64 `json:"agent_id"`
	FontNumber    string `json:"font_number"`
	Message       string `json:"message"`
	MessageType   string `json:"message_type"`
	MessageID     string `json:"message_id,omitempty"`
	MediaURL      string `json:"media_url"`
}
