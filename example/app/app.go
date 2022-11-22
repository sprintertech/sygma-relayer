// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package app

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChainSafe/sygma-relayer/comm"

	"github.com/ChainSafe/chainbridge-core/lvldb"

	"github.com/ethereum/go-ethereum/common"
	secp256k1 "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	coreEvm "github.com/ChainSafe/chainbridge-core/chains/evm"
	coreEvents "github.com/ChainSafe/chainbridge-core/chains/evm/calls/events"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	coreExecutor "github.com/ChainSafe/chainbridge-core/chains/evm/executor"
	coreListener "github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-core/e2e/dummy"
	"github.com/ChainSafe/chainbridge-core/flags"
	"github.com/ChainSafe/chainbridge-core/opentelemetry"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ChainSafe/chainbridge-core/store"

	"github.com/ChainSafe/sygma-relayer/chains/evm"
	"github.com/ChainSafe/sygma-relayer/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-relayer/chains/evm/calls/events"
	"github.com/ChainSafe/sygma-relayer/chains/evm/executor"
	"github.com/ChainSafe/sygma-relayer/chains/evm/listener"
	"github.com/ChainSafe/sygma-relayer/comm/elector"
	"github.com/ChainSafe/sygma-relayer/comm/p2p"
	"github.com/ChainSafe/sygma-relayer/config"
	"github.com/ChainSafe/sygma-relayer/keyshare"
	"github.com/ChainSafe/sygma-relayer/topology"
	"github.com/ChainSafe/sygma-relayer/tss"
)

func Run() error {
	configuration, err := config.GetConfigFromFile(viper.GetString(flags.ConfigFlagName))
	if err != nil {
		panic(err)
	}

	networkTopology, _ := topology.ProcessRawTopology(&topology.RawTopology{
		Peers: []topology.RawPeer{
			{PeerAddress: "/dns4/relayer2/tcp/9001/p2p/QmeTuMtdpPB7zKDgmobEwSvxodrf5aFVSmBXX3SQJVjJaT"},
			{PeerAddress: "/dns4/relayer3/tcp/9002/p2p/QmYAYuLUPNwYEBYJaKHcE7NKjUhiUV8txx2xDXHvcYa1xK"},
			{PeerAddress: "/dns4/relayer1/tcp/9000/p2p/QmcvEg7jGvuxdsUFRUiE4VdrL2P1Yeju5L83BsJvvXz7zX"},
		},
		Threshold: "2",
	})

	db, err := lvldb.NewLvlDB(viper.GetString(flags.BlockstoreFlagName))
	if err != nil {
		panic(err)
	}
	blockstore := store.NewBlockStore(db)

	privBytes, err := crypto.ConfigDecodeKey(configuration.RelayerConfig.MpcConfig.Key)
	if err != nil {
		panic(err)
	}
	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	if err != nil {
		panic(err)
	}

	connectionGate := p2p.NewConnectionGate(networkTopology)
	host, err := p2p.NewHost(priv, networkTopology, connectionGate, configuration.RelayerConfig.MpcConfig.Port)
	if err != nil {
		panic(err)
	}

	healthComm := p2p.NewCommunication(host, "p2p/health")
	go comm.ExecuteCommHealthCheck(healthComm, host.Peerstore().Peers())

	communication := p2p.NewCommunication(host, "p2p/sygma")
	electorFactory := elector.NewCoordinatorElectorFactory(host, configuration.RelayerConfig.BullyConfig)
	coordinator := tss.NewCoordinator(host, communication, electorFactory)
	keyshareStore := keyshare.NewKeyshareStore(configuration.RelayerConfig.MpcConfig.KeysharePath)

	chains := []relayer.RelayedChain{}
	for _, chainConfig := range configuration.ChainConfigs {
		switch chainConfig["type"] {
		case "evm":
			{
				config, err := chain.NewEVMConfig(chainConfig)
				if err != nil {
					panic(err)
				}

				privateKey, err := secp256k1.HexToECDSA(config.GeneralChainConfig.Key)
				if err != nil {
					panic(err)
				}

				client, err := evmclient.NewEVMClient(config.GeneralChainConfig.Endpoint, privateKey)
				if err != nil {
					panic(err)
				}

				mod := big.NewInt(0).Mod(config.StartBlock, config.BlockConfirmations)
				// startBlock % blockConfirmations == 0
				if mod.Cmp(big.NewInt(0)) != 0 {
					config.StartBlock.Sub(config.StartBlock, mod)
				}

				bridgeAddress := common.HexToAddress(config.Bridge)
				dummyGasPricer := dummy.NewStaticGasPriceDeterminant(client, nil)
				t := signAndSend.NewSignAndSendTransactor(evmtransaction.NewTransaction, dummyGasPricer, client)
				bridgeContract := bridge.NewBridgeContract(client, bridgeAddress, t)

				pGenericHandler := chainConfig["permissionlessGenericHandler"].(string)
				depositHandler := coreListener.NewETHDepositHandler(bridgeContract)
				depositHandler.RegisterDepositHandler(config.Erc20Handler, coreListener.Erc20DepositHandler)
				depositHandler.RegisterDepositHandler(config.Erc721Handler, coreListener.Erc721DepositHandler)
				depositHandler.RegisterDepositHandler(config.GenericHandler, coreListener.GenericDepositHandler)
				depositHandler.RegisterDepositHandler(pGenericHandler, listener.PermissionlessGenericDepositHandler)
				depositListener := coreEvents.NewListener(client)
				tssListener := events.NewListener(client)
				eventHandlers := make([]coreListener.EventHandler, 0)
				eventHandlers = append(eventHandlers, listener.NewDepositEventHandler(depositListener, depositHandler, bridgeAddress, *config.GeneralChainConfig.Id))
				eventHandlers = append(eventHandlers, listener.NewKeygenEventHandler(tssListener, coordinator, host, communication, keyshareStore, bridgeAddress, networkTopology.Threshold))
				eventHandlers = append(eventHandlers, listener.NewRefreshEventHandler(nil, nil, tssListener, coordinator, host, communication, connectionGate, keyshareStore, bridgeAddress))
				eventHandlers = append(eventHandlers, listener.NewRetryEventHandler(tssListener, depositHandler, bridgeAddress, *config.GeneralChainConfig.Id, config.BlockConfirmations))
				evmListener := coreListener.NewEVMListener(client, eventHandlers, blockstore, config)

				mh := coreExecutor.NewEVMMessageHandler(bridgeContract)
				mh.RegisterMessageHandler(config.Erc20Handler, coreExecutor.ERC20MessageHandler)
				mh.RegisterMessageHandler(config.Erc721Handler, coreExecutor.ERC721MessageHandler)
				mh.RegisterMessageHandler(config.GenericHandler, coreExecutor.GenericMessageHandler)
				mh.RegisterMessageHandler(pGenericHandler, executor.PermissionlessGenericMessageHandler)
				executor := executor.NewExecutor(host, communication, coordinator, mh, bridgeContract, keyshareStore)

				coreEvmChain := coreEvm.NewEVMChain(evmListener, nil, blockstore, config)
				chain := evm.NewEVMChain(*coreEvmChain, executor)

				chains = append(chains, chain)
			}
		default:
			panic(fmt.Errorf("type '%s' not recognized", chainConfig["type"]))
		}
	}

	r := relayer.NewRelayer(
		chains,
		&opentelemetry.ConsoleTelemetry{},
	)

	errChn := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Start(ctx, errChn)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	select {
	case err := <-errChn:
		log.Error().Err(err).Msg("failed to listen and serve")
		return err
	case sig := <-sysErr:
		log.Info().Msgf("terminating got ` [%v] signal", sig)
		return nil
	}
}
