package core

import (
	erismint "github.com/eris-ltd/eris-db/manager/eris-mint"
	stypes "github.com/eris-ltd/eris-db/manager/eris-mint/state/types"

	bc "github.com/tendermint/tendermint/blockchain"
	server "github.com/eris-ltd/eris-db/server"
	"github.com/tendermint/tendermint/consensus"
	mempl "github.com/tendermint/tendermint/mempool"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/tendermint/go-p2p"
)

var blockStore *bc.BlockStore
var consensusState *consensus.ConsensusState
var consensusReactor *consensus.ConsensusReactor
var mempoolReactor *mempl.MempoolReactor
var p2pSwitch *p2p.Switch
var privValidator *tmtypes.PrivValidator
var genDoc *stypes.GenesisDoc // cache the genesis structure
var erisMint *erismint.ErisMint

var config server.ServerConfig = server.ServerConfig{}

func SetConfig(c server.ServerConfig) {
	config = c
}

func SetErisMint(em *erismint.ErisMint) {
	erisMint = em
}

func SetBlockStore(bs *bc.BlockStore) {
	blockStore = bs
}

func SetConsensusState(cs *consensus.ConsensusState) {
	consensusState = cs
}

func SetConsensusReactor(cr *consensus.ConsensusReactor) {
	consensusReactor = cr
}

func SetMempoolReactor(mr *mempl.MempoolReactor) {
	mempoolReactor = mr
}

func SetSwitch(sw *p2p.Switch) {
	p2pSwitch = sw
}

func SetPrivValidator(pv *tmtypes.PrivValidator) {
	privValidator = pv
}

func SetGenDoc(doc *stypes.GenesisDoc) {
	genDoc = doc
}