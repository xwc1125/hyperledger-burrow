// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/validator"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/consensus/tendermint"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/names"
	"github.com/hyperledger/burrow/execution/registry"
	"github.com/hyperledger/burrow/genesis"
	"github.com/hyperledger/burrow/txs"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/consensus"
	ctypes "github.com/tendermint/tendermint/consensus/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// When using Tendermint types like Block and Vote we are forced to wrap the outer object and use amino marshalling
var aminoCodec = NewAminoCodec()

func NewAminoCodec() *amino.Codec {
	aminoCodec := amino.NewCodec()
	consensus.RegisterMessages(aminoCodec)
	core_types.RegisterAmino(aminoCodec)
	return aminoCodec
}

type ResultStorage struct {
	Key   binary.HexBytes
	Value binary.HexBytes
}

type ResultAccounts struct {
	BlockHeight uint64
	Accounts    []*acm.Account
}

type ResultDumpStorage struct {
	StorageItems []StorageItem
}

type StorageItem struct {
	Key   binary.HexBytes
	Value binary.HexBytes
}

type ResultBlocks struct {
	LastHeight uint64
	BlockMetas []*tmTypes.BlockMeta
}

type ResultBlock struct {
	BlockMeta *BlockMeta
	Block     *Block
}

type BlockMeta struct {
	*tmTypes.BlockMeta
}

func (bm BlockMeta) MarshalJSON() ([]byte, error) {
	return aminoCodec.MarshalJSON(bm.BlockMeta)
}

func (bm *BlockMeta) UnmarshalJSON(data []byte) (err error) {
	return aminoCodec.UnmarshalJSON(data, &bm.BlockMeta)
}

// Needed for go-amino handling of interface types
type Block struct {
	*tmTypes.Block
}

func (b Block) MarshalJSON() ([]byte, error) {
	return aminoCodec.MarshalJSON(b.Block)
}

func (b *Block) UnmarshalJSON(data []byte) (err error) {
	return aminoCodec.UnmarshalJSON(data, &b.Block)
}

type ResultChainId struct {
	ChainName   string
	ChainId     string
	GenesisHash binary.HexBytes
}

type ResultSubscribe struct {
	EventID        string
	SubscriptionID string
}

type ResultUnsubscribe struct {
	SubscriptionID string
}

type ResultNetwork struct {
	ThisNode *tendermint.NodeInfo
	*core_types.ResultNetInfo
}

type ResultNetworkRegistry struct {
	Address crypto.Address
	registry.NodeIdentity
}

type ResultValidators struct {
	BlockHeight         uint64
	BondedValidators    []*validator.Validator
	UnbondingValidators []*validator.Validator
}

type ResultConsensusState struct {
	*core_types.ResultDumpConsensusState
}

// TODO use round state in ResultConsensusState - currently there are some part of RoundState have no Unmarshal
type RoundState struct {
	*ctypes.RoundState
}

func (rs RoundState) MarshalJSON() ([]byte, error) {
	return aminoCodec.MarshalJSON(rs.RoundState)
}

func (rs *RoundState) UnmarshalJSON(data []byte) (err error) {
	return aminoCodec.UnmarshalJSON(data, &rs.RoundState)
}

type ResultPeers struct {
	Peers []core_types.Peer
}

type ResultNames struct {
	BlockHeight uint64
	Names       []*names.Entry
}

type ResultGeneratePrivateAccount struct {
	PrivateAccount *acm.ConcretePrivateAccount
}

type ResultAccount struct {
	Account *acm.Account
}

type AccountHumanReadable struct {
	Address     crypto.Address
	PublicKey   crypto.PublicKey
	Sequence    uint64
	Balance     uint64
	Code        []string
	Permissions []string
	Roles       []string
}

type ResultAccountHumanReadable struct {
	Account *AccountHumanReadable
}

type ResultAccountStats struct {
	AccountsWithCode    uint64
	AccountsWithoutCode uint64
}

type ResultUnconfirmedTxs struct {
	NumTxs int
	Txs    []*txs.Envelope
}

type ResultName struct {
	Entry *names.Entry
}

type ResultGenesis struct {
	Genesis genesis.GenesisDoc
}

type ResultSignTx struct {
	Tx *txs.Envelope
}
