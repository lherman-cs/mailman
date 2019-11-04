package peer

import "github.com/pion/webrtc/v2"

type peerConnection struct {
	*webrtc.PeerConnection
	chanIce         <-chan webrtc.ICECandidateInit
	chanDataChannel <-chan *dataChannel
}

func newPeerConnection(raw *webrtc.PeerConnection) *peerConnection {
	chanIce := make(chan webrtc.ICECandidateInit)
	chanDataChannel := make(chan *dataChannel)

	raw.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice == nil {
			return
		}
		chanIce <- ice.ToJSON()
	})

	raw.OnDataChannel(func(channel *webrtc.DataChannel) {
		chanDataChannel <- newDataChannel(channel)
	})

	p := peerConnection{
		PeerConnection:  raw,
		chanIce:         chanIce,
		chanDataChannel: chanDataChannel,
	}
	return &p
}

type dataChannel struct {
	*webrtc.DataChannel
	chanOpen <-chan struct{}
}

func newDataChannel(raw *webrtc.DataChannel) *dataChannel {
	chanOpen := make(chan struct{})
	raw.OnOpen(func() {
		chanOpen <- struct{}{}
	})

	c := dataChannel{DataChannel: raw, chanOpen: chanOpen}
	return &c
}
