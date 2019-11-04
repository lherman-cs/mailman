package message

import (
	"testing"

	"github.com/pion/webrtc/v2"
)

func TestEncodeAndDecode(t *testing.T) {
	ice := ICEMessage(webrtc.ICECandidateInit{})
	_, ok := ice.Payload.(webrtc.ICECandidateInit)
	if !ok {
		t.Fatal("payload is not set properly")
	}
	encoded, err := ice.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var m Message
	err = m.UnmarshalJSON(encoded)
	if err != nil {
		t.Fatal(err)
	}

	_, ok = m.Payload.(webrtc.ICECandidateInit)
	if !ok {
		t.Fatal("failed in decoding")
	}
}
