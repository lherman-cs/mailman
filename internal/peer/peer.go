package peer

import (
	"context"
	"fmt"
	"io"

	msg "github.com/lherman-cs/mailman/internal/message"
	"github.com/pion/webrtc/v2"
	"go.uber.org/zap"
)

// TODO! Replace this hardcoded config
var config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

// Peer is an abstraction that exposes readwritecloser
type Peer struct {
	c      Conn
	l      *zap.SugaredLogger
	pc     *peerConnection
	ctx    context.Context
	cancel context.CancelFunc
}

// NewPeer initializes Peer with c as the adapter connection
func NewPeer(ctx context.Context, c Conn, l *zap.SugaredLogger) (*Peer, error) {
	// Create a SettingEngine and enable Detach
	s := webrtc.SettingEngine{}
	s.DetachDataChannels()

	// Create an API object with the engine
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))
	// Create a new RTCPeerConnection using the API object
	pc, err := api.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &Peer{
		c:      c,
		l:      l,
		pc:     newPeerConnection(pc),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Connect connects to the receiver and finished webrtc flow until connected.
func (p *Peer) Connect() (io.ReadWriteCloser, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := p.startListening(ctx)

	offer, err := p.pc.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	err = p.pc.SetLocalDescription(offer)
	if err != nil {
		return nil, err
	}

	err = p.c.WriteJSON(msg.OfferMessage(offer))
	if err != nil {
		return nil, err
	}

	channel, err := p.pc.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}
	wrappedChannel := newDataChannel(channel)

	p.l.Infof("sender waiting for the data channel to open\n")
	var rwc io.ReadWriteCloser
	select {
	case <-wrappedChannel.chanOpen:
		rwc, err = p.onConnected(wrappedChannel)
	case err = <-errChan:
	}

	if err == nil {
		p.l.Infof("connected to a channel\n")
	} else {
		p.l.Infof("an error occured while waiting for the data channel to open\n")
		p.l.Infof("%s\n", sprintError(err))
	}

	return rwc, err
}

// Wait waits until a sender arrives and finishes webrtc flow until connected.
func (p *Peer) Wait() (io.ReadWriteCloser, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p.l.Infof("waiting for a sender\n")
	errChan := p.startListening(ctx)
	var d *dataChannel
	var err error
	p.l.Infof("waiting for an incoming channel\n")
	var rwc io.ReadWriteCloser
	select {
	case err = <-errChan:
	case d = <-p.pc.chanDataChannel:
		rwc, err = p.onConnected(d)
	}

	if err == nil {
		p.l.Infof("connected to a channel\n")
	} else {
		p.l.Infof("an error occured while waiting for an incoming channel\n")
		p.l.Infof("%s\n", sprintError(err))
	}

	return rwc, err
}

func (p *Peer) Close() {
	p.cancel()
}

func (p *Peer) startListening(ctx context.Context) <-chan error {
	errChan := make(chan error)
	go func() {
		err := p.handleWebRTCMessages(ctx)
		if err != nil {
			errChan <- fmt.Errorf("on handleWebRTCMessages: %w", err)
		}
	}()

	go func() {
		err := p.handleSignalingMessages(ctx)
		if err != nil {
			errChan <- fmt.Errorf("on handleSignalingMessages: %w", err)
		}
	}()

	return errChan
}

func (p *Peer) onConnected(channel *dataChannel) (io.ReadWriteCloser, error) {
	// Detach the data channel
	raw, err := channel.Detach()
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (p *Peer) handleWebRTCMessages(ctx context.Context) error {
	for {
		select {
		case ice := <-p.pc.chanIce:
			err := p.c.WriteJSON(msg.ICEMessage(ice))
			if err != nil {
				return fmt.Errorf("on sending ice message: %w", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("on handleWebRTCMessages done: %w", ctx.Err())
		}
	}
}

func (p *Peer) handleSignalingMessages(ctx context.Context) error {
	var err error

	for {
		// ctx has been cancelled
		if err = ctx.Err(); err != nil {
			return nil
		}

		var m msg.Message
		if err = p.c.ReadJSON(&m); err != nil {
			return fmt.Errorf("on reading signaling message: %w", err)
		}

		switch m.Type {
		case "offer":
			offer := m.Payload.(webrtc.SessionDescription)
			if err = p.handleOffer(offer); err != nil {
				return fmt.Errorf("on handleOffer: %w", err)
			}
		case "answer":
			answer := m.Payload.(webrtc.SessionDescription)
			if err = p.handleAnswer(answer); err != nil {
				return fmt.Errorf("on handleAnswer: %w", err)
			}
		case "ice":
			ice := m.Payload.(webrtc.ICECandidateInit)
			if err = p.handleICE(ice); err != nil {
				return fmt.Errorf("on handleICE: %w", err)
			}
		}
	}
}

func (p *Peer) handleOffer(offer webrtc.SessionDescription) error {
	p.pc.SetRemoteDescription(offer)

	answer, err := p.pc.CreateAnswer(nil)
	if err != nil {
		return err
	}

	p.pc.SetLocalDescription(answer)
	err = p.c.WriteJSON(msg.AnswerMessage(answer))
	if err != nil {
		return err
	}

	return nil
}

func (p *Peer) handleAnswer(answer webrtc.SessionDescription) error {
	err := p.pc.SetRemoteDescription(answer)
	return err
}

func (p *Peer) handleICE(ice webrtc.ICECandidateInit) error {
	return p.pc.AddICECandidate(ice)
}
