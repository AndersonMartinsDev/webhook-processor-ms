package model

type SimplifiedMessagePayload struct {
	OriginMessage string `json:"origin_message"`
	SessionKey    string `json:"session_key"`
	FontNumber    string `json:"font_number"`
	MessageType   string `json:"message_type"`
	AgentID       uint64 `json:"agent_id"`
	Message       string `json:"message"`
	MediaURL      string `json:"media_url"`
}
