// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package execution

import (
	"context"
	"testing"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/bcm"
	"github.com/hyperledger/burrow/consensus/tendermint/codes"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/event"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/keys"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/txs"
	"github.com/hyperledger/burrow/txs/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/mempool"
	tmTypes "github.com/tendermint/tendermint/types"
)

func TestTransactor_BroadcastTxSync(t *testing.T) {
	chainID := "TestChain"
	bc := &bcm.Blockchain{}
	evc := event.NewEmitter()
	evc.SetLogger(logging.NewNoopLogger())
	txCodec := txs.NewProtobufCodec()
	privAccount := acm.GeneratePrivateAccountFromSecret("frogs")
	tx := &payload.CallTx{
		Input: &payload.TxInput{
			Address: privAccount.GetAddress(),
		},
		Address: &crypto.Address{1, 2, 3},
	}
	txEnv := txs.Enclose(chainID, tx)
	err := txEnv.Sign(privAccount)
	require.NoError(t, err)
	height := uint64(35)
	trans := NewTransactor(bc, evc, NewAccounts(acmstate.NewMemoryState(),
		keys.NewLocalKeyClient(keys.NewMemoryKeyStore(privAccount), logger), 100),
		func(tx tmTypes.Tx, cb func(*abciTypes.Response), txInfo mempool.TxInfo) error {
			txe := exec.NewTxExecution(txEnv)
			txe.Height = height
			err := evc.Publish(context.Background(), txe, txe)
			if err != nil {
				return err
			}
			bs, err := txe.Receipt.Encode()
			if err != nil {
				return err
			}
			cb(abciTypes.ToResponseCheckTx(abciTypes.ResponseCheckTx{
				Code: codes.TxExecutionSuccessCode,
				Data: bs,
			}))
			return nil
		}, "", txCodec, logger)
	txe, err := trans.BroadcastTxSync(context.Background(), txEnv)
	require.NoError(t, err)
	assert.Equal(t, height, txe.Height)
}
