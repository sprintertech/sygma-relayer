// The Licensed Work is (c) 2022 Sygma
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"

	"github.com/ChainSafe/sygma-relayer/chains/substrate/connection"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/author"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/rs/zerolog/log"
)

type SubstrateClient struct {
	client.Client
	author.Author
	key       *signature.KeyringPair // Keyring used for signing
	ChainID   *big.Int
	nonceLock sync.Mutex // Locks nonce for updates
	nonce     types.U32  // Latest account nonce
}

func NewSubstrateClient(url string, key *signature.KeyringPair) (*SubstrateClient, error) {
	c := &SubstrateClient{
		key: key,
	}
	client, err := client.Connect(url)
	if err != nil {
		return nil, err
	}
	c.Client = client
	c.ChainID, _ = new(big.Int).SetString("6277101735386680764176071790128604879584176795969512275969", 10)
	return c, nil

}

// SendRawTransaction accepts rlp-encode of signed transaction and sends it via RPC call
func (c *SubstrateClient) sendRawTransaction(ext types.Extrinsic) (types.Hash, error) {
	fmt.Println("sendrawtrannsssss")
	enc, err := codec.EncodeToHex(ext)

	if err != nil {
		return types.Hash{}, err
	}

	var res string
	err = c.Call(&res, "author_submitExtrinsic", enc)
	fmt.Println("erererer")
	fmt.Println(err)
	if err != nil {
		return types.Hash{}, err
	}

	return types.NewHashFromHexString(res)
	/* 	hsh, err := c.Author.SubmitExtrinsic(ext)
	   	fmt.Println(err)
	   	fmt.Println(hsh)
	   	if err != nil {
	   		return types.Hash{}, err
	   	}
	   	return hsh, nil */
}

func (c *SubstrateClient) signAndSendTransaction(opts types.SignatureOptions, ext types.Extrinsic) (types.Hash, error) {

	err := ext.Sign(*c.key, opts)
	if err != nil {
		c.nonceLock.Unlock()
		return types.Hash{}, err
	}
	hash, err := c.sendRawTransaction(ext)

	if err != nil {
		return types.Hash{}, err
	}
	return hash, nil
}

// Transact constructs and submits an extrinsic to call the method with the given arguments.
// All args are passed directly into GSRPC. GSRPC types are recommended to avoid serialization inconsistencies.
func (c *SubstrateClient) Transact(conn *connection.Connection, method string, args ...interface{}) (*types.Hash, error) {
	log.Debug().Msgf("Submitting substrate call... method %s, sender %s", method, c.key.Address)

	meta := conn.GetMetadata()

	// Create call and extrinsic
	fmt.Println("MethooooooodddddddddddddddddddSUbbBBBBBstrateeeeee\nn")
	fmt.Println(method)
	fmt.Println(args...)
	call, err := types.NewCall(
		&meta,
		method,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct call: %w", err)
	}
	ext := types.NewExtrinsic(call)

	// Get latest runtime version
	rv, err := conn.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	c.nonceLock.Lock()

	key, err := types.CreateStorageKey(&meta, "System", "Account", c.key.PublicKey, nil)
	if err != nil {
		return nil, err
	}
	var latestNonce types.U32
	var acct types.AccountInfo
	exists, err := conn.RPC.State.GetStorageLatest(key, &acct)
	if err != nil {
		c.nonceLock.Unlock()
		return nil, err
	}
	if !exists {
		latestNonce = 0
	} else {
		latestNonce = acct.Nonce
	}

	if latestNonce > c.nonce {
		c.nonce = latestNonce
	}

	// Sign the extrinsic
	o := types.SignatureOptions{
		BlockHash:          conn.GenesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        conn.GenesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(c.nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	h, err := c.signAndSendTransaction(o, ext)
	c.nonce++
	c.nonceLock.Unlock()
	if err != nil {
		return nil, fmt.Errorf("submission of extrinsic failed: %w", err)
	}
	log.Trace().Msg("Extrinsic submission succeeded")

	return &h, nil
}
