// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/bcm"
	"github.com/hyperledger/burrow/consensus/tendermint/codes"
	"github.com/hyperledger/burrow/event"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/hyperledger/burrow/execution/exec"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/logging/structure"
	"github.com/hyperledger/burrow/txs"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/p2p"
	tmTypes "github.com/tendermint/tendermint/types"
)

const (
	SubscribeBufferSize = 10
)

type txChecker func(tx tmTypes.Tx, callback func(*abciTypes.Response), txInfo mempool.TxInfo) error

// Transactor is responsible for helping to formulate, sign, and broadcast transactions to tendermint
//
// The BroadcastTx* methods are able to work against the mempool Accounts (pending) state rather than the
// committed (final) Accounts state and can assign a sequence number based on all of the txs
// seen since the last block - provided these transactions are successfully committed (via DeliverTx) then
// subsequent transactions will have valid sequence numbers. This allows Burrow to coordinate sequencing and signing
// for a key it holds or is provided - it is down to the key-holder to manage the mutual information between transactions
// concurrent within a new block window.
type Transactor struct {
	BlockchainInfo  bcm.BlockchainInfo
	Emitter         *event.Emitter
	MempoolAccounts *Accounts
	checkTxAsync    txChecker
	nodeID          p2p.ID
	txEncoder       txs.Encoder
	logger          *logging.Logger
}

func NewTransactor(tip bcm.BlockchainInfo, emitter *event.Emitter, mempoolAccounts *Accounts,
	checkTxAsync txChecker, id p2p.ID, txEncoder txs.Encoder, logger *logging.Logger) *Transactor {

	return &Transactor{
		BlockchainInfo:  tip,
		Emitter:         emitter,
		MempoolAccounts: mempoolAccounts,
		checkTxAsync:    checkTxAsync,
		nodeID:          id,
		txEncoder:       txEncoder,
		logger:          logger.With(structure.ComponentKey, "Transactor"),
	}
}

func (trans *Transactor) BroadcastTxSync(ctx context.Context, txEnv *txs.Envelope) (*exec.TxExecution, error) {
	// Sign unless already signed - note we must attempt signing before subscribing so we get accurate final TxHash
	unlock, err := trans.MaybeSignTxMempool(txEnv)
	if err != nil {
		return nil, err
	}
	// We will try and call this before the function exits unless we error but it is idempotent
	defer unlock()
	// Subscribe before submitting to mempool
	txHash := txEnv.Tx.Hash()
	subID := event.GenSubID()
	out, err := trans.Emitter.Subscribe(ctx, subID, exec.QueryForTxExecution(txHash), SubscribeBufferSize)
	if err != nil {
		// We do not want to hold the lock with a defer so we must
		return nil, err
	}
	defer trans.Emitter.UnsubscribeAll(context.Background(), subID)
	// Push Tx to mempool
	checkTxReceipt, err := trans.CheckTxSync(ctx, txEnv)
	unlock()
	if err != nil {
		return nil, err
	}
	// Get all the execution events for this Tx
	select {
	case <-ctx.Done():
		syncInfo := bcm.GetSyncInfo(trans.BlockchainInfo)
		bs, err := json.Marshal(syncInfo)
		syncInfoString := string(bs)
		if err != nil {
			syncInfoString = fmt.Sprintf("{error could not marshal SyncInfo: %v}", err)
		}
		return nil, fmt.Errorf("waiting for tx %v, SyncInfo: %s", checkTxReceipt.TxHash, syncInfoString)
	case msg := <-out:
		txe := msg.(*exec.TxExecution)
		callError := txe.CallError()
		if callError != nil && callError.ErrorCode() != errors.Codes.ExecutionReverted {
			return nil, errors.Wrap(callError, "exception during transaction execution")
		}
		return txe, nil
	}
}

// Broadcast a transaction without waiting for confirmation - will attempt to sign server-side and set sequence numbers
// if no signatures are provided
func (trans *Transactor) BroadcastTxAsync(ctx context.Context, txEnv *txs.Envelope) (*txs.Receipt, error) {
	return trans.CheckTxSync(ctx, txEnv)
}

// Broadcast a transaction and waits for a response from the mempool. Transactions to BroadcastTx will block during
// various mempool operations (managed by Tendermint) including mempool Reap, Commit, and recheckTx.
func (trans *Transactor) CheckTxSync(ctx context.Context, txEnv *txs.Envelope) (*txs.Receipt, error) {
	trans.logger.Trace.Log("method", "CheckTxSync",
		structure.TxHashKey, txEnv.Tx.Hash(),
		"tx", txEnv.String())
	// Sign unless already signed
	unlock, err := trans.MaybeSignTxMempool(txEnv)
	if err != nil {
		return nil, err
	}
	defer unlock()
	err = txEnv.Validate()
	if err != nil {
		return nil, err
	}
	txBytes, err := trans.txEncoder.EncodeTx(txEnv)
	if err != nil {
		return nil, err
	}
	return trans.CheckTxSyncRaw(ctx, txBytes)
}

func (trans *Transactor) MaybeSignTxMempool(txEnv *txs.Envelope) (UnlockFunc, error) {
	// Sign unless already signed
	if len(txEnv.Signatories) == 0 {
		var err error
		var unlock UnlockFunc
		// We are writing signatures back to txEnv so don't shadow txEnv here
		txEnv, unlock, err = trans.SignTxMempool(txEnv)
		if err != nil {
			return nil, fmt.Errorf("error signing transaction: %v", err)
		}
		// Hash will have change since we signed
		txEnv.Tx.Rehash()
		// Make this idempotent for defer
		var once sync.Once
		return func() { once.Do(unlock) }, nil
	}
	return func() {}, nil
}

func (trans *Transactor) SignTxMempool(txEnv *txs.Envelope) (*txs.Envelope, UnlockFunc, error) {
	inputs := txEnv.Tx.GetInputs()
	signers := make([]acm.AddressableSigner, len(inputs))
	unlockers := make([]UnlockFunc, len(inputs))
	for i, input := range inputs {
		ssa, err := trans.MempoolAccounts.SequentialSigningAccount(input.Address)
		if err != nil {
			return nil, nil, err
		}
		sa, unlock, err := ssa.Lock()
		if err != nil {
			return nil, nil, err
		}
		// Hold lock until safely in mempool - important that this is held until after CheckTxSync returns
		unlockers[i] = unlock
		signers[i] = sa
		// Set sequence number consecutively from mempool
		input.Sequence = sa.Sequence + 1
	}

	err := txEnv.Sign(signers...)
	if err != nil {
		return nil, nil, err
	}
	return txEnv, UnlockFunc(func() {
		for _, unlock := range unlockers {
			defer unlock()
		}
	}), nil
}

func (trans *Transactor) SignTx(txEnv *txs.Envelope) (*txs.Envelope, error) {
	var err error
	inputs := txEnv.Tx.GetInputs()
	signers := make([]acm.AddressableSigner, len(inputs))
	for i, input := range inputs {
		signers[i], err = trans.MempoolAccounts.SigningAccount(input.Address)
		if err != nil {
			return nil, err
		}
	}
	err = txEnv.Sign(signers...)
	if err != nil {
		return nil, err
	}
	return txEnv, nil
}

func (trans *Transactor) CheckTxSyncRaw(ctx context.Context, txBytes []byte) (*txs.Receipt, error) {
	responseCh := make(chan *abciTypes.Response, 3)
	err := trans.CheckTxAsyncRaw(txBytes, func(res *abciTypes.Response) {
		responseCh <- res
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("waiting for CheckTx response in CheckTxSyncRaw: %v", ctx.Err())
	case response := <-responseCh:
		checkTxResponse := response.GetCheckTx()
		if checkTxResponse == nil {
			return nil, fmt.Errorf("application did not return CheckTx response")
		}

		switch checkTxResponse.Code {
		case codes.TxExecutionSuccessCode:
			receipt, err := txs.DecodeReceipt(checkTxResponse.Data)
			if err != nil {
				return nil, fmt.Errorf("could not deserialise transaction receipt: %s", err)
			}
			return receipt, nil
		default:
			return nil, fmt.Errorf("error %d returned by Tendermint in BroadcastTxSync ABCI log: %v",
				checkTxResponse.Code, checkTxResponse.Log)
		}
	}
}

func (trans *Transactor) CheckTxAsyncRaw(txBytes []byte, callback func(res *abciTypes.Response)) error {
	return trans.checkTxAsync(txBytes, callback, mempool.TxInfo{SenderP2PID: trans.nodeID})
}

func (trans *Transactor) CheckTxAsync(txEnv *txs.Envelope, callback func(res *abciTypes.Response)) error {
	err := txEnv.Validate()
	if err != nil {
		return err
	}
	txBytes, err := trans.txEncoder.EncodeTx(txEnv)
	if err != nil {
		return fmt.Errorf("error encoding transaction: %v", err)
	}
	return trans.CheckTxAsyncRaw(txBytes, callback)
}
