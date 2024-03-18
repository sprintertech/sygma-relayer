package erc20_test

import (
	"math/big"
	"testing"

	erc20 "github.com/ChainSafe/sygma-relayer/e2e/evm/contracts/erc20"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"github.com/sygmaprotocol/sygma-core/chains/evm/transactor"
	"github.com/sygmaprotocol/sygma-core/mock"
	"go.uber.org/mock/gomock"
)

type ERC20ContractCallsTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
	mockClient       *mock.MockClient
	mockTransactor   *mock.MockTransactor
	erc20contract    *erc20.ERC20Contract
}

var (
	testContractAddress   = "0x5f75ce92326e304962b22749bd71e36976171285"
	testInteractorAddress = "0x8362bbbd6d987895E2A4630a55e69Dd8C7b9f87B"
)

func TestRunERC20ContractCallsTestSuite(t *testing.T) {
	suite.Run(t, new(ERC20ContractCallsTestSuite))
}

func (s *ERC20ContractCallsTestSuite) SetupSuite()    {}
func (s *ERC20ContractCallsTestSuite) TearDownSuite() {}
func (s *ERC20ContractCallsTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockClient = mock.NewMockClient(s.gomockController)
	s.mockTransactor = mock.NewMockTransactor(s.gomockController)
	s.erc20contract = erc20.NewERC20Contract(
		s.mockClient, common.HexToAddress(testContractAddress), s.mockTransactor,
	)
}
func (s *ERC20ContractCallsTestSuite) TearDownTest() {}

func (s *ERC20ContractCallsTestSuite) TestErc20Contract_GetBalance_Success() {
	s.mockClient.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5}, nil)
	res, err := s.erc20contract.GetBalance(common.HexToAddress(testInteractorAddress))
	s.Equal(
		big.NewInt(5),
		res,
	)
	s.Nil(err)
}

func (s *ERC20ContractCallsTestSuite) TestErc20Contract_MintTokens_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3, 4, 5}, nil)
	res, err := s.erc20contract.MintTokens(common.HexToAddress(testInteractorAddress), big.NewInt(10), transactor.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3, 4, 5},
		res,
	)
	s.Nil(err)
}

func (s *ERC20ContractCallsTestSuite) TestErc20Contract_ApproveTokens_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)
	res, err := s.erc20contract.ApproveTokens(common.HexToAddress(testInteractorAddress), big.NewInt(100), transactor.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9},
		res,
	)
	s.Nil(err)
}

func (s *ERC20ContractCallsTestSuite) TestErc20Contract_MinterRole_Success() {
	s.mockClient.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10}, nil)
	res, err := s.erc20contract.MinterRole()
	s.Equal(
		[32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10},
		res,
	)
	s.Nil(err)
}

func (s *ERC20ContractCallsTestSuite) TestErc20Contract_AddMinter_Success() {
	s.mockClient.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockClient.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 25}, nil)
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3}, nil)
	res, err := s.erc20contract.AddMinter(common.HexToAddress(testInteractorAddress), transactor.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3},
		res,
	)
	s.Nil(err)
}
