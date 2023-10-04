// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package executor

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/binance-chain/tss-lib/common"
	"github.com/sourcegraph/conc/pool"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/rs/zerolog/log"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/executor/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/sygma-relayer/chains"
	"github.com/ChainSafe/sygma-relayer/comm"
	"github.com/ChainSafe/sygma-relayer/tss"
	"github.com/ChainSafe/sygma-relayer/tss/signing"
)

const TRANSFER_GAS_COST = 200000

type Batch struct {
	proposals []*chains.Proposal
	gasLimit  uint64
}

var (
	executionCheckPeriod = time.Minute
	signingTimeout       = 30 * time.Minute
)

type MessageHandler interface {
	HandleMessage(m *message.Message) (*proposal.Proposal, error)
}

type BridgeContract interface {
	IsProposalExecuted(p *chains.Proposal) (bool, error)
	ExecuteProposals(proposals []*chains.Proposal, signature []byte, opts transactor.TransactOptions) (*ethCommon.Hash, error)
	ProposalsHash(proposals []*chains.Proposal) ([]byte, error)
}

type Executor struct {
	coordinator       *tss.Coordinator
	host              host.Host
	comm              comm.Communication
	fetcher           signing.SaveDataFetcher
	bridge            BridgeContract
	mh                MessageHandler
	exitLock          *sync.RWMutex
	transactionMaxGas uint64
}

func NewExecutor(
	host host.Host,
	comm comm.Communication,
	coordinator *tss.Coordinator,
	mh MessageHandler,
	bridgeContract BridgeContract,
	fetcher signing.SaveDataFetcher,
	exitLock *sync.RWMutex,
	transactionMaxGas uint64,
) *Executor {
	return &Executor{
		host:              host,
		comm:              comm,
		coordinator:       coordinator,
		mh:                mh,
		bridge:            bridgeContract,
		fetcher:           fetcher,
		exitLock:          exitLock,
		transactionMaxGas: transactionMaxGas,
	}
}

// Execute starts a signing process and executes proposals when signature is generated
func (e *Executor) Execute(msgs []*message.Message) error {
	e.exitLock.RLock()
	defer e.exitLock.RUnlock()

	batches, err := e.proposalBatches(msgs)
	if err != nil {
		return err
	}

	p := pool.New().WithErrors()
	for _, batch := range batches {
		if len(batch.proposals) == 0 {
			continue
		}

		b := batch
		p.Go(func() error {
			propHash, err := e.bridge.ProposalsHash(b.proposals)
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
			executionContext, cancelExecution := context.WithCancel(context.Background())
			watchContext, cancelWatch := context.WithCancel(context.Background())
			ep := pool.New().WithErrors()
			ep.Go(func() error {
				err := e.coordinator.Execute(executionContext, signing, sigChn)
				if err != nil {
					cancelWatch()
				}

				return err
			})
			ep.Go(func() error { return e.watchExecution(watchContext, cancelExecution, b, sigChn, sessionID) })
			return ep.Wait()
		})
	}
	return p.Wait()
}

func (e *Executor) watchExecution(ctx context.Context, cancelExecution context.CancelFunc, batch *Batch, sigChn chan interface{}, sessionID string) error {
	ticker := time.NewTicker(executionCheckPeriod)
	timeout := time.NewTicker(signingTimeout)
	defer ticker.Stop()
	defer timeout.Stop()
	defer cancelExecution()

	for {
		select {
		case sigResult := <-sigChn:
			{
				cancelExecution()
				if sigResult == nil {
					continue
				}

				signatureData := sigResult.(*common.SignatureData)
				hash, err := e.executeBatch(batch, signatureData)
				if err != nil {
					_ = e.comm.Broadcast(e.host.Peerstore().Peers(), []byte{}, comm.TssFailMsg, sessionID)
					return err
				}

				log.Info().Str("SessionID", sessionID).Msgf("Sent proposals execution with hash: %s", hash)
			}
		case <-ticker.C:
			{
				if !e.areProposalsExecuted(batch.proposals, sessionID) {
					continue
				}

				log.Info().Str("SessionID", sessionID).Msgf("Successfully executed proposals")
				return nil
			}
		case <-timeout.C:
			{
				return fmt.Errorf("execution timed out in %s", signingTimeout)
			}
		case <-ctx.Done():
			{
				return nil
			}
		}
	}
}

func (e *Executor) proposalBatches(msgs []*message.Message) ([]*Batch, error) {
	batches := make([]*Batch, 1)
	currentBatch := &Batch{
		proposals: make([]*chains.Proposal, 0),
		gasLimit:  0,
	}
	batches[0] = currentBatch

	for _, m := range msgs {
		prop, err := e.mh.HandleMessage(m)
		if err != nil {
			return nil, err
		}

		evmProposal := chains.NewProposal(prop.Source, prop.Destination, prop.DepositNonce, prop.ResourceId, prop.Data, prop.Metadata)
		isExecuted, err := e.bridge.IsProposalExecuted(evmProposal)
		if err != nil {
			return nil, err
		}
		if isExecuted {
			log.Info().Msgf("Proposal %p already executed", prop)
			continue
		}

		var propGasLimit uint64
		l, ok := evmProposal.Metadata.Data["gasLimit"]
		if ok {
			propGasLimit = l.(uint64)
		} else {
			propGasLimit = uint64(TRANSFER_GAS_COST)
		}
		currentBatch.gasLimit += propGasLimit
		if currentBatch.gasLimit >= e.transactionMaxGas {
			currentBatch = &Batch{
				proposals: make([]*chains.Proposal, 0),
				gasLimit:  0,
			}
			batches = append(batches, currentBatch)
		}

		currentBatch.proposals = append(currentBatch.proposals, evmProposal)
	}

	return batches, nil
}

func (e *Executor) executeBatch(batch *Batch, signatureData *common.SignatureData) (*ethCommon.Hash, error) {
	sig := []byte{}
	sig = append(sig[:], ethCommon.LeftPadBytes(signatureData.R, 32)...)
	sig = append(sig[:], ethCommon.LeftPadBytes(signatureData.S, 32)...)
	sig = append(sig[:], signatureData.SignatureRecovery...)
	sig[len(sig)-1] += 27 // Transform V from 0/1 to 27/28

	hash, err := e.bridge.ExecuteProposals(batch.proposals, sig, transactor.TransactOptions{
		GasLimit: batch.gasLimit,
	})
	if err != nil {
		return nil, err
	}

	return hash, err
}

func (e *Executor) areProposalsExecuted(proposals []*chains.Proposal, sessionID string) bool {
	for _, prop := range proposals {
		isExecuted, err := e.bridge.IsProposalExecuted(prop)
		if err != nil || !isExecuted {
			return false
		}
	}

	return true
}

func (e *Executor) sessionID(hash []byte) string {
	return fmt.Sprintf("signing-%s", ethCommon.Bytes2Hex(hash))
}
