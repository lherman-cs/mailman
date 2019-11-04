package message

import (
	"encoding/json"
	"fmt"

	"github.com/pion/webrtc/v2"
)

type Message struct {
	Type           string      `json:"type"`
	EncodedPayload []byte      `json:"payload"`
	Payload        interface{} `json:"-"`
}

func (m *Message) UnmarshalJSON(b []byte) error {
	type Alias Message
	tmp := struct{ *Alias }{Alias: (*Alias)(m)}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch m.Type {
	case "offer", "answer":
		var payload webrtc.SessionDescription
		err = json.Unmarshal(m.EncodedPayload, &payload)
		m.Payload = payload
	case "ice":
		var payload webrtc.ICECandidateInit
		err = json.Unmarshal(m.EncodedPayload, &payload)
		m.Payload = payload
	default:
		err = fmt.Errorf("invalid Message type")
	}
	return err
}

func (m *Message) MarshalJSON() ([]byte, error) {
	encoded, err := json.Marshal(&m.Payload)
	if err != nil {
		return nil, err
	}

	m.EncodedPayload = encoded
	type Alias Message

	return json.Marshal(&struct{ *Alias }{
		Alias: (*Alias)(m),
	})
}

func ICEMessage(ice webrtc.ICECandidateInit) *Message {
	return &Message{
		Type:    "ice",
		Payload: ice,
	}
}

func OfferMessage(offer webrtc.SessionDescription) *Message {
	return &Message{
		Type:    "offer",
		Payload: offer,
	}
}

func AnswerMessage(answer webrtc.SessionDescription) *Message {
	return &Message{
		Type:    "answer",
		Payload: answer,
	}
}

func StateMessage(state string) *Message {
	return &Message{
		Type:    "state",
		Payload: state,
	}
}
