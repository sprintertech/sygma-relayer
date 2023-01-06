package executor

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-core/chains/evm/executor/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
)

func PermissionlessGenericMessageHandler(msg *message.Message, handlerAddr, bridgeAddress common.Address) (*proposal.Proposal, error) {
	executeFunctionSignature, ok := msg.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong function signature format")
	}
	executeContractAddress, ok := msg.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong contract address format")
	}
	maxFee, ok := msg.Payload[2].([]byte)
	if !ok {
		return nil, errors.New("wrong max fee format")
	}
	executionData, ok := msg.Payload[3].([]byte)
	if !ok {
		return nil, errors.New("wrong execution data format")
	}

	data := bytes.Buffer{}
	data.Write(common.LeftPadBytes(maxFee, 32))

	data.Write(common.LeftPadBytes(big.NewInt(int64(len(executeFunctionSignature))).Bytes(), 2))
	data.Write(executeFunctionSignature)

	data.Write([]byte{byte(len(executeContractAddress))})
	data.Write(executeContractAddress)

	data.Write([]byte{byte(len(executionData))})
	data.Write(executionData)

	return proposal.NewProposal(msg.Source, msg.Destination, msg.DepositNonce, msg.ResourceId, data.Bytes(), handlerAddr, bridgeAddress, msg.Metadata), nil
}
