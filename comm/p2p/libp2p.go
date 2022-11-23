// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package p2p

import (
	"context"
	"encoding/json"
	"fmt"

	comm "github.com/ChainSafe/sygma-relayer/comm"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Libp2pCommunication struct {
	SessionSubscriptionManager
	h             host.Host
	protocolID    protocol.ID
	streamManager *StreamManager
	logger        zerolog.Logger
}

func NewCommunication(h host.Host, protocolID protocol.ID) Libp2pCommunication {
	logger := log.With().Str("Module", "communication").Str("Peer", h.ID().Pretty()).Logger()
	c := Libp2pCommunication{
		SessionSubscriptionManager: NewSessionSubscriptionManager(),
		h:                          h,
		protocolID:                 protocolID,
		streamManager:              NewStreamManager(),
		logger:                     logger,
	}

	// start processing incoming messages
	c.h.SetStreamHandler(c.protocolID, c.StreamHandlerFunc)
	return c
}

/** Communication interface methods **/

func (c Libp2pCommunication) Broadcast(
	peers peer.IDSlice,
	msg []byte,
	msgType comm.MessageType,
	sessionID string,
	errChan chan error,
) {
	hostID := c.h.ID()
	wMsg := comm.WrappedMessage{
		MessageType: msgType,
		SessionID:   sessionID,
		Payload:     msg,
		From:        hostID,
	}
	marshaledMsg, err := json.Marshal(wMsg)
	if err != nil {
		c.logger.Error().Err(err).Str("SessionID", sessionID).Msg("unable to marshal message")
		return
	}
	c.logger.Debug().Str("MsgType", msgType.String()).Str("SessionID", sessionID).Msg(
		"broadcasting message",
	)
	for _, peerID := range peers {
		if hostID == peerID {
			continue // don't send message to itself
		}
		go func(peerID peer.ID) {
			err := c.sendMessage(peerID, marshaledMsg, msgType, sessionID)
			if err != nil {
				SendError(errChan, err, peerID)
				return
			}
		}(peerID)
	}
}

func (c Libp2pCommunication) Subscribe(
	sessionID string,
	msgType comm.MessageType,
	channel chan *comm.WrappedMessage,
) comm.SubscriptionID {
	subID := c.SubscribeTo(sessionID, msgType, channel)
	c.logger.Trace().Str(
		"SessionID", sessionID).Msgf(
		"subscribed to message type %s", msgType,
	)
	return subID
}

func (c Libp2pCommunication) UnSubscribe(
	subID comm.SubscriptionID,
) {
	c.UnSubscribeFrom(subID)
	c.logger.Trace().Str(
		"SessionID", subID.SessionID()).Str(
		"SubID", subID.SubscriptionIdentifier()).Msgf(
		"unsubscribed from message type %s", subID.MessageType().String(),
	)
}

/** Helper methods **/

func (c Libp2pCommunication) sendMessage(
	to peer.ID,
	msg []byte,
	msgType comm.MessageType,
	sessionID string,
) error {
	pi := c.h.Peerstore().PeerInfo(to)
	resolver, err := madns.NewResolver()
	if err != nil {
		return err
	}
	if len(pi.Addrs) == 0 {
		return fmt.Errorf("peer %s has no defined addresses", to)
	}

	addr, err := resolver.Resolve(context.Background(), pi.Addrs[0])
	if err != nil {
		return err
	}
	err = c.h.Connect(context.TODO(), peer.AddrInfo{
		ID:    to,
		Addrs: addr,
	})
	if err != nil {
		return err
	}

	stream, err := c.h.NewStream(context.TODO(), to, c.protocolID)
	if err != nil {
		c.logger.Error().Err(err).Str("MsgType", msgType.String()).Str("SessionID", sessionID).Msgf(
			"unable to open stream toward %s", to.Pretty(),
		)
		return err
	}

	err = WriteStream(msg, stream)
	if err != nil {
		c.logger.Error().Str("To", string(to)).Err(err).Msg("unable to send message")
		return err
	}
	c.logger.Trace().Str(
		"To", to.Pretty()).Str(
		"MsgType", msgType.String()).Str(
		"SessionID", sessionID).Msg(
		"message sent",
	)
	c.streamManager.AddStream(sessionID, stream)
	return nil
}

func (c Libp2pCommunication) StreamHandlerFunc(s network.Stream) {
	msg, err := c.ProcessMessageFromStream(s)
	if err != nil {
		c.logger.Error().Err(err).Str("StreamID", s.ID()).Msg("unable to process message")
		return
	}

	subscribers := c.GetSubscribers(msg.SessionID, msg.MessageType)
	for _, sub := range subscribers {
		sub := sub
		go func() {
			sub <- msg
		}()
	}
}

func (c Libp2pCommunication) ProcessMessageFromStream(s network.Stream) (*comm.WrappedMessage, error) {
	remotePeerID := s.Conn().RemotePeer()
	msgBytes, err := ReadStream(s)
	if err != nil {
		c.streamManager.AddStream("UNKNOWN", s)
		return nil, err
	}

	var wrappedMsg comm.WrappedMessage
	if err := json.Unmarshal(msgBytes, &wrappedMsg); nil != err {
		c.streamManager.AddStream("UNKNOWN", s)
		return nil, err
	}

	wrappedMsg.From = remotePeerID

	c.streamManager.AddStream(wrappedMsg.SessionID, s)

	c.logger.Trace().Str(
		"From", wrappedMsg.From.Pretty()).Str(
		"MsgType", wrappedMsg.MessageType.String()).Str(
		"SessionID", wrappedMsg.SessionID).Msg(
		"processed message",
	)

	return &wrappedMsg, nil
}
