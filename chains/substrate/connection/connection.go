// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package connection

import (
	"sync"

	"github.com/ChainSafe/sygma-relayer/chains/substrate/events"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain"

	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Connection struct {
	chain.Chain
	*rpc.RPC
	meta        types.Metadata // Latest chain metadata
	metaLock    sync.RWMutex   // Lock metadata for updates, allows concurrent reads
	genesisHash types.Hash     // Chain genesis hash
}

func NewSubstrateConnection(url string) (*Connection, error) {
	c := &Connection{}
	client, err := client.Connect(url)
	if err != nil {
		return nil, err
	}
	rpc, err := rpc.NewRPC(client)
	if err != nil {
		return nil, err
	}
	c.RPC = rpc
	c.Chain = rpc.Chain

	// Fetch metadata
	meta, err := c.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	c.meta = *meta
	// Fetch genesis hash
	genesisHash, err := c.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}
	c.genesisHash = genesisHash
	return c, nil
}

func (c *Connection) GetMetadata() (meta types.Metadata) {
	c.metaLock.RLock()
	meta = c.meta
	c.metaLock.RUnlock()
	return meta
}

func (c *Connection) UpdateMetatdata() error {
	c.metaLock.Lock()
	meta, err := c.RPC.State.GetMetadataLatest()
	if err != nil {
		c.metaLock.Unlock()
		return err
	}
	c.meta = *meta
	c.metaLock.Unlock()
	return nil
}

func (c *Connection) GetBlockEvents(hash types.Hash, target events.Events) error {
	meta := c.GetMetadata()
	key, err := types.CreateStorageKey(&meta, "System", "Events", nil, nil)
	if err != nil {
		return err
	}

	var records types.EventRecordsRaw
	_, err = c.RPC.State.GetStorage(key, &records, hash)
	if err != nil {
		return err
	}

	err = records.DecodeEventRecords(&meta, &target)
	if err != nil {
		return err
	}
	return nil
}
