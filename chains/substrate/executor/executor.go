// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package executor

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/sygma-relayer/chains"
	"github.com/ChainSafe/sygma-relayer/chains/substrate/connection"
	"github.com/binance-chain/tss-lib/common"

	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/rs/zerolog/log"

	"github.com/ChainSafe/chainbridge-core/relayer/message"

	"github.com/ChainSafe/sygma-relayer/comm"
	"github.com/ChainSafe/sygma-relayer/tss"
	"github.com/ChainSafe/sygma-relayer/tss/signing"
)

var (
	executionCheckPeriod = time.Minute
	signingTimeout       = 30 * time.Minute
)

type MessageHandler interface {
	HandleMessage(m *message.Message) (*chains.Proposal, error)
}

type BridgePallet interface {
	IsProposalExecuted(p *chains.Proposal) (bool, error)
	ExecuteProposals(proposals []*chains.Proposal, signature []byte) (string, *author.ExtrinsicStatusSubscription, error)
	ProposalsHash(proposals []*chains.Proposal) ([]byte, error)
	TrackExtrinsic(extHash string, sub *author.ExtrinsicStatusSubscription, errChn chan error)
}

type Executor struct {
	coordinator *tss.Coordinator
	host        host.Host
	comm        comm.Communication
	fetcher     signing.SaveDataFetcher
	bridge      BridgePallet
	mh          MessageHandler
	conn        *connection.Connection
}

func NewExecutor(
	host host.Host,
	comm comm.Communication,
	coordinator *tss.Coordinator,
	mh MessageHandler,
	bridgePallet BridgePallet,
	fetcher signing.SaveDataFetcher,
	conn *connection.Connection,
) *Executor {
	return &Executor{
		host:        host,
		comm:        comm,
		coordinator: coordinator,
		mh:          mh,
		bridge:      bridgePallet,
		fetcher:     fetcher,
		conn:        conn,
	}
}

// Execute starts a signing process and executes proposals when signature is generated
func (e *Executor) Execute(msgs []*message.Message) error {
	proposals := make([]*chains.Proposal, 0)
	for _, m := range msgs {
		prop, err := e.mh.HandleMessage(m)
		if err != nil {
			return err
		}

		isExecuted, err := e.bridge.IsProposalExecuted(prop)
		if err != nil {
			return err
		}
		if isExecuted {
			continue
		}

		proposals = append(proposals, prop)
	}
	if len(proposals) == 0 {
		return nil
	}

	propHash, err := e.bridge.ProposalsHash(proposals)
	if err != nil {
		return err
	}

	sessionID := e.sessionID(propHash)
	msg := big.NewInt(0)
	msg.SetBytes(propHash)
	signing, err := signing.NewSigning(
		msg,
		e.sessionID(propHash),
		e.host,
		e.comm,
		e.fetcher)
	if err != nil {
		return err
	}

	sigChn := make(chan interface{})
	statusChn := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	go e.coordinator.Execute(ctx, signing, sigChn, statusChn)

	ticker := time.NewTicker(executionCheckPeriod)
	timeout := time.NewTicker(signingTimeout)
	defer ticker.Stop()
	defer timeout.Stop()
	defer cancel()
	for {
		select {
		case sigResult := <-sigChn:
			{
				signatureData := sigResult.(*common.SignatureData)
				hash, sub, err := e.executeProposal(proposals, signatureData)
				if err != nil {
					go e.comm.Broadcast(e.host.Peerstore().Peers(), []byte{}, comm.TssFailMsg, sessionID, nil)
					return err
				}
				errChn := make(chan error)
				go e.bridge.TrackExtrinsic(hash, sub, errChn)
				err = <-errChn
				return err
			}
		case err := <-statusChn:
			{
				return err
			}
		case <-ticker.C:
			{
				allExecuted := true
				for _, prop := range proposals {
					isExecuted, err := e.bridge.IsProposalExecuted(prop)
					if err != nil {
						return err
					}
					if !isExecuted {
						allExecuted = false
						continue
					}

					log.Info().Msgf("Successfully executed proposal %v", prop)
				}

				if allExecuted {
					return nil
				}
			}
		case <-timeout.C:
			{
				return fmt.Errorf("execution timed out in %s", signingTimeout)
			}
		}
	}
}

func (e *Executor) executeProposal(proposals []*chains.Proposal, signatureData *common.SignatureData) (string, *author.ExtrinsicStatusSubscription, error) {
	sig := []byte{}
	sig = append(sig[:], ethCommon.LeftPadBytes(signatureData.R, 32)...)
	sig = append(sig[:], ethCommon.LeftPadBytes(signatureData.S, 32)...)
	sig = append(sig[:], signatureData.SignatureRecovery...)
	sig[len(sig)-1] += 27 // Transform V from 0/1 to 27/28

	hash, sub, err := e.bridge.ExecuteProposals(proposals, sig)
	if err != nil {
		return "", nil, err
	}

	return hash, sub, err
}

func (e *Executor) sessionID(hash []byte) string {
	return fmt.Sprintf("signing-%s", ethCommon.Bytes2Hex(hash))
}
