// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package common_test

import (
	"testing"

	"github.com/ChainSafe/sygma-relayer/tss/common"
	"github.com/stretchr/testify/suite"
)

type TssMessageTestSuite struct {
	suite.Suite
}

func TestRunTssMessageTestSuite(t *testing.T) {
	suite.Run(t, new(TssMessageTestSuite))
}

func (s *TssMessageTestSuite) Test_UnmarshaledMessageShouldBeEqual() {
	originalMsg := &common.TssMessage{
		MsgBytes:    []byte{1},
		IsBroadcast: true,
	}
	msgBytes, err := common.MarshalTssMessage(originalMsg.MsgBytes, originalMsg.IsBroadcast)
	s.Nil(err)

	unmarshaledMsg, err := common.UnmarshalTssMessage(msgBytes)
	s.Nil(err)

	s.Equal(originalMsg, unmarshaledMsg)
}

type StartMessageTestSuite struct {
	suite.Suite
}

func TestRunStartMessageTestSuite(t *testing.T) {
	suite.Run(t, new(StartMessageTestSuite))
}

func (s *StartMessageTestSuite) Test_UnmarshaledMessageShouldBeEqual() {
	originalMsg := &common.StartMessage{
		Params: []byte("test"),
	}
	msgBytes, err := common.MarshalStartMessage(originalMsg.Params)
	s.Nil(err)

	unmarshaledMsg, err := common.UnmarshalStartMessage(msgBytes)
	s.Nil(err)

	s.Equal(originalMsg, unmarshaledMsg)
}
