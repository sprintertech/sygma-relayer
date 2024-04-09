// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: LGPL-3.0-only

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	substrateListener "github.com/ChainSafe/sygma-relayer/chains/substrate/listener"
	substratePallet "github.com/ChainSafe/sygma-relayer/chains/substrate/pallet"
	"github.com/ChainSafe/sygma-relayer/relayer/transfer"
	"github.com/sygmaprotocol/sygma-core/chains/evm/listener"
	"github.com/sygmaprotocol/sygma-core/chains/evm/transactor/gas"
	"github.com/sygmaprotocol/sygma-core/chains/evm/transactor/transaction"
	coreSubstrate "github.com/sygmaprotocol/sygma-core/chains/substrate"
	substrateClient "github.com/sygmaprotocol/sygma-core/chains/substrate/client"
	"github.com/sygmaprotocol/sygma-core/chains/substrate/connection"
	coreSubstrateListener "github.com/sygmaprotocol/sygma-core/chains/substrate/listener"
	"github.com/sygmaprotocol/sygma-core/crypto/secp256k1"
	"github.com/sygmaprotocol/sygma-core/observability"
	"github.com/sygmaprotocol/sygma-core/relayer"
	"github.com/sygmaprotocol/sygma-core/store"
	"github.com/sygmaprotocol/sygma-core/store/lvldb"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/sygmaprotocol/sygma-core/chains/evm/transactor/monitored"
	"github.com/sygmaprotocol/sygma-core/relayer/message"

	"github.com/ChainSafe/sygma-relayer/chains/evm"
	"github.com/ChainSafe/sygma-relayer/chains/substrate"
	substrateExecutor "github.com/ChainSafe/sygma-relayer/chains/substrate/executor"
	"github.com/ChainSafe/sygma-relayer/jobs"
	"github.com/ChainSafe/sygma-relayer/metrics"
	coreEvm "github.com/sygmaprotocol/sygma-core/chains/evm"

	"github.com/ChainSafe/sygma-relayer/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-relayer/chains/evm/calls/events"
	"github.com/ChainSafe/sygma-relayer/chains/evm/executor"
	"github.com/ChainSafe/sygma-relayer/chains/evm/listener/depositHandlers"
	hubEventHandlers "github.com/ChainSafe/sygma-relayer/chains/evm/listener/eventHandlers"
	"github.com/ChainSafe/sygma-relayer/comm/elector"
	"github.com/ChainSafe/sygma-relayer/comm/p2p"
	"github.com/ChainSafe/sygma-relayer/config"
	"github.com/ChainSafe/sygma-relayer/keyshare"
	"github.com/ChainSafe/sygma-relayer/topology"
	"github.com/ChainSafe/sygma-relayer/tss"
	evmClient "github.com/sygmaprotocol/sygma-core/chains/evm/client"
)

func Run() error {
	configuration := &config.Config{}
	configuration, err := config.GetConfigFromFile(viper.GetString(config.ConfigFlagName), configuration)
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

	db, err := lvldb.NewLvlDB(viper.GetString(config.BlockstoreFlagName))
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

	communication := p2p.NewCommunication(host, "p2p/sygma")
	electorFactory := elector.NewCoordinatorElectorFactory(host, configuration.RelayerConfig.BullyConfig)
	coordinator := tss.NewCoordinator(host, communication, electorFactory)
	keyshareStore := keyshare.NewKeyshareStore(configuration.RelayerConfig.MpcConfig.KeysharePath)

	// wait until executions are done and then stop further executions before exiting
	exitLock := &sync.RWMutex{}
	defer exitLock.Lock()

	mp, err := observability.InitMetricProvider(context.Background(), configuration.RelayerConfig.OpenTelemetryCollectorURL)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Error().Msgf("Error shutting down meter provider: %v", err)
		}
	}()
	sygmaMetrics, err := metrics.NewSygmaMetrics(mp.Meter("relayer-metric-provider"), configuration.RelayerConfig.Env, configuration.RelayerConfig.Id)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	msgChan := make(chan []*message.Message)
	chains := make(map[uint8]relayer.RelayedChain)
	for _, chainConfig := range configuration.ChainConfigs {
		switch chainConfig["type"] {
		case "evm":
			{
				config, err := evm.NewEVMConfig(chainConfig)
				if err != nil {
					panic(err)
				}

				kp, err := secp256k1.NewKeypairFromString(config.GeneralChainConfig.Key)
				if err != nil {
					panic(err)
				}

				client, err := evmClient.NewEVMClient(config.GeneralChainConfig.Endpoint, kp)
				if err != nil {
					panic(err)
				}

				log.Info().Str("domain", config.String()).Msgf("Registering EVM domain")

				bridgeAddress := common.HexToAddress(config.Bridge)
				dummyGasPricer := gas.NewStaticGasPriceDeterminant(client, nil)
				t := monitored.NewMonitoredTransactor(transaction.NewTransaction, dummyGasPricer, client, config.MaxGasPrice, config.GasIncreasePercentage)
				go t.Monitor(ctx, time.Minute*3, time.Minute*10, time.Minute)

				bridgeContract := bridge.NewBridgeContract(client, bridgeAddress, t)

				depositHandler := depositHandlers.NewETHDepositHandler(bridgeContract)
				mh := message.NewMessageHandler()
				mh.RegisterMessageHandler(transfer.TransferMessageType, &executor.TransferMessageHandler{})
				//mh.RegisterMessageHandler("FungibleTransfer", &executor.TransferMessageHandler{})
				for _, handler := range config.Handlers {

					switch handler.Type {
					case "erc20":
						{
							depositHandler.RegisterDepositHandler(handler.Address, &depositHandlers.Erc20DepositHandler{})
						}
					case "permissionedGeneric":
						{
							depositHandler.RegisterDepositHandler(handler.Address, &depositHandlers.GenericDepositHandler{})
						}
					case "permissionlessGeneric":
						{
							depositHandler.RegisterDepositHandler(handler.Address, &depositHandlers.PermissionlessGenericDepositHandler{})
						}
					case "erc721":
						{
							depositHandler.RegisterDepositHandler(handler.Address, &depositHandlers.Erc721DepositHandler{})
						}
					case "erc1155":
						{
							depositHandler.RegisterDepositHandler(handler.Address, &depositHandlers.Erc1155DepositHandler{})
						}
					}
				}
				depositListener := events.NewListener(client)
				tssListener := events.NewListener(client)
				eventHandlers := make([]listener.EventHandler, 0)
				l := log.With().Str("chain", fmt.Sprintf("%v", config.GeneralChainConfig.Name)).Uint8("domainID", *config.GeneralChainConfig.Id)
				eventHandlers = append(eventHandlers, hubEventHandlers.NewDepositEventHandler(depositListener, depositHandler, bridgeAddress, *config.GeneralChainConfig.Id, msgChan))
				eventHandlers = append(eventHandlers, hubEventHandlers.NewKeygenEventHandler(l, tssListener, coordinator, host, communication, keyshareStore, bridgeAddress, networkTopology.Threshold))
				eventHandlers = append(eventHandlers, hubEventHandlers.NewRefreshEventHandler(l, nil, nil, tssListener, coordinator, host, communication, connectionGate, keyshareStore, bridgeAddress))
				eventHandlers = append(eventHandlers, hubEventHandlers.NewRetryEventHandler(l, tssListener, depositHandler, bridgeAddress, *config.GeneralChainConfig.Id, config.BlockConfirmations, msgChan))
				evmListener := listener.NewEVMListener(client, eventHandlers, blockstore, sygmaMetrics, *config.GeneralChainConfig.Id, config.BlockRetryInterval, config.BlockConfirmations, config.BlockInterval)
				executor := executor.NewExecutor(host, communication, coordinator, bridgeContract, keyshareStore, exitLock, config.GasLimit.Uint64())

				chain := coreEvm.NewEVMChain(evmListener, mh, executor, *config.GeneralChainConfig.Id, config.StartBlock)

				chains[*config.GeneralChainConfig.Id] = chain
			}
		case "substrate":
			{
				config, err := substrate.NewSubstrateConfig(chainConfig)
				if err != nil {
					panic(err)
				}

				conn, err := connection.NewSubstrateConnection(config.GeneralChainConfig.Endpoint)
				if err != nil {
					panic(err)
				}

				keyPair, err := signature.KeyringPairFromSecret(config.GeneralChainConfig.Key, config.SubstrateNetwork)
				if err != nil {
					panic(err)
				}

				log.Info().Str("domain", config.String()).Msgf("Registering substrate domain")

				substrateClient := substrateClient.NewSubstrateClient(conn, &keyPair, config.ChainID, config.Tip)
				bridgePallet := substratePallet.NewPallet(substrateClient)

				l := log.With().Str("chain", fmt.Sprintf("%v", config.GeneralChainConfig.Name)).Uint8("domainID", *config.GeneralChainConfig.Id)
				depositHandler := substrateListener.NewSubstrateDepositHandler()
				depositHandler.RegisterDepositHandler(transfer.FungibleTransfer, substrateListener.FungibleTransferHandler)
				eventHandlers := make([]coreSubstrateListener.EventHandler, 0)
				eventHandlers = append(eventHandlers, substrateListener.NewFungibleTransferEventHandler(l, *config.GeneralChainConfig.Id, depositHandler, msgChan, conn))
				eventHandlers = append(eventHandlers, substrateListener.NewRetryEventHandler(l, conn, depositHandler, *config.GeneralChainConfig.Id, msgChan))

				substrateListener := coreSubstrateListener.NewSubstrateListener(conn, eventHandlers, blockstore, sygmaMetrics, *config.GeneralChainConfig.Id, config.BlockRetryInterval, config.BlockInterval)

				mh := message.NewMessageHandler()
				mh.RegisterMessageHandler(transfer.TransferMessageType, &substrateExecutor.SubstrateMessageHandler{})

				sExecutor := substrateExecutor.NewExecutor(host, communication, coordinator, bridgePallet, keyshareStore, conn, exitLock)
				substrateChain := coreSubstrate.NewSubstrateChain(substrateListener, mh, sExecutor, *config.GeneralChainConfig.Id, config.StartBlock)
				chains[*config.GeneralChainConfig.Id] = substrateChain
			}
		default:
			panic(fmt.Errorf("type '%s' not recognized", chainConfig["type"]))
		}
	}

	go jobs.StartCommunicationHealthCheckJob(host, configuration.RelayerConfig.MpcConfig.CommHealthCheckInterval, sygmaMetrics)

	r := relayer.NewRelayer(chains)

	go r.Start(ctx, msgChan)

	sysErr := make(chan os.Signal, 1)
	signal.Notify(sysErr,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT)

	sig := <-sysErr
	log.Info().Msgf("terminating got ` [%v] signal", sig)
	return nil

}
