package rpc_test

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hyperledger/burrow/acm/balance"
	"github.com/hyperledger/burrow/crypto"
	x "github.com/hyperledger/burrow/encoding/hex"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/integration"
	"github.com/hyperledger/burrow/keys"
	"github.com/hyperledger/burrow/logging"
	"github.com/hyperledger/burrow/project"
	"github.com/hyperledger/burrow/rpc"
	"github.com/hyperledger/burrow/rpc/web3"
	"github.com/stretchr/testify/require"
)

func TestWeb3Service(t *testing.T) {
	ctx := context.Background()
	genesisAccounts := integration.MakePrivateAccounts("burrow", 1)
	genesisAccounts = append(genesisAccounts, integration.MakeEthereumAccounts("ethereum", 3)...)
	genesisDoc := integration.TestGenesisDoc(genesisAccounts, 0)

	config, _ := integration.NewTestConfig(genesisDoc)
	logger := logging.NewNoopLogger()
	kern, err := integration.TestKernel(genesisAccounts[0], genesisAccounts, config)
	require.NoError(t, err)
	err = kern.Boot()
	defer kern.Shutdown(ctx)

	dir, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store := keys.NewFilesystemKeyStore(dir, true)
	for _, acc := range genesisAccounts {
		err = store.StoreKeyPlain(&keys.Key{
			CurveType:  acc.PrivateKey().CurveType,
			Address:    acc.GetAddress(),
			PublicKey:  acc.GetPublicKey(),
			PrivateKey: acc.PrivateKey(),
		})
		require.NoError(t, err)
	}

	nodeView, err := kern.GetNodeView()
	require.NoError(t, err)

	accountState := kern.State
	eventsState := kern.State
	validatorState := kern.State
	eth := rpc.NewEthService(accountState, eventsState, kern.Blockchain, validatorState,
		nodeView, kern.Transactor, store, kern.Logger)

	t.Run("Web3Sha3", func(t *testing.T) {
		result, err := eth.Web3Sha3(&web3.Web3Sha3Params{"0x68656c6c6f20776f726c64"}) // hello world
		require.NoError(t, err)
		// hex encoded
		require.Equal(t, "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad", result.HashedData)
	})

	t.Run("NetListening", func(t *testing.T) {
		result, err := eth.NetListening()
		require.NoError(t, err)
		require.Equal(t, true, result.IsNetListening)
	})

	t.Run("NetPeerCount", func(t *testing.T) {
		result, err := eth.NetPeerCount()
		require.NoError(t, err)
		require.Equal(t, "0x0", result.NumConnectedPeers)
	})

	t.Run("Version+ID", func(t *testing.T) {
		t.Run("Web3ClientVersion", func(t *testing.T) {
			result, err := eth.Web3ClientVersion()
			require.NoError(t, err)
			require.Equal(t, project.FullVersion(), result.ClientVersion)
		})

		t.Run("NetVersion", func(t *testing.T) {
			result, err := eth.NetVersion()
			require.NoError(t, err)
			require.Equal(t, x.EncodeNumber(1), result.ChainID)
		})

		t.Run("EthProtocolVersion", func(t *testing.T) {
			result, err := eth.EthProtocolVersion()
			require.NoError(t, err)
			require.NotEmpty(t, result.ProtocolVersion)
		})

		t.Run("EthChainId", func(t *testing.T) {
			result, err := eth.EthChainId()
			require.NoError(t, err)
			doc := config.GenesisDoc
			require.Equal(t, doc.ChainID(), result.ChainId)
		})
	})

	t.Run("EthTransactions", func(t *testing.T) {
		var txHash, contractAddress string

		receivee := genesisAccounts[2].GetPublicKey().GetAddress()
		acc, err := kern.State.GetAccount(receivee)
		require.NoError(t, err)
		before := acc.GetBalance()

		t.Run("EthSendRawTransaction", func(t *testing.T) {
			// see: https://github.com/ethereumjs/ethereumjs-tx/blob/master/examples/transactions.ts#L9
			raw := `0xf867808082520894f97798df751deb4b6e39d4cf998ee7cd4dcb9acc880de0b6b3a76400008025a0f0d2396973296cd6a71141c974d4a851f5eae8f08a8fba2dc36a0fef9bd6440ca0171995aa750d3f9f8e4d0eac93ff67634274f3c5acf422723f49ff09a6885422`
			_, err = eth.EthSendRawTransaction(&web3.EthSendRawTransactionParams{
				SignedTransactionData: raw,
			})
			require.NoError(t, err)
		})

		t.Run("EthGetBalance", func(t *testing.T) {
			result, err := eth.EthGetBalance(&web3.EthGetBalanceParams{
				Address:     x.EncodeBytes(receivee.Bytes()),
				BlockNumber: "latest",
			})
			require.NoError(t, err)
			after, err := x.DecodeToBigInt(result.GetBalanceResult)
			require.NoError(t, err)
			after = balance.WeiToNative(after.Bytes())
			require.Equal(t, after.Uint64(), before+1)
		})

		t.Run("EthGetTransactionCount", func(t *testing.T) {
			result, err := eth.EthGetTransactionCount(&web3.EthGetTransactionCountParams{
				Address: genesisAccounts[1].GetAddress().String(),
			})
			require.NoError(t, err)
			require.Equal(t, x.EncodeNumber(1), result.NonceOrNull)
		})

		// create contract on chain
		t.Run("EthSendTransaction", func(t *testing.T) {
			type ret struct {
				*web3.EthSendTransactionResult
				error
			}
			ch := make(chan ret)
			numSends := 5
			for i := 0; i < numSends; i++ {
				idx := i
				go func() {
					tx := &web3.EthSendTransactionParams{
						Transaction: web3.Transaction{
							From: x.EncodeBytes(genesisAccounts[3].GetAddress().Bytes()),
							Gas:  x.EncodeNumber(uint64(40 + idx)), // make tx unique in mempool
							Data: x.EncodeBytes(rpc.Bytecode_HelloWorld),
						},
					}
					result, err := eth.EthSendTransaction(tx)
					ch <- ret{result, err}
				}()
			}
			for i := 0; i < numSends; i++ {
				select {
				case r := <-ch:
					require.NoError(t, r.error)
					require.NotEmpty(t, r.TransactionHash)
				case <-time.After(2 * time.Second):
					t.Fatalf("timed out waiting for EthSendTransaction result")
				}
			}
		})

		t.Run("EthGetTransactionReceipt", func(t *testing.T) {
			sendResult, err := eth.EthSendTransaction(&web3.EthSendTransactionParams{
				Transaction: web3.Transaction{
					From: x.EncodeBytes(genesisAccounts[3].GetAddress().Bytes()),
					Gas:  x.EncodeNumber(40),
					Data: x.EncodeBytes(rpc.Bytecode_HelloWorld),
				},
			})
			require.NoError(t, err)
			txHash = sendResult.TransactionHash
			require.NotEmpty(t, txHash, "need tx hash to get tx receipt")
			receiptResult, err := eth.EthGetTransactionReceipt(&web3.EthGetTransactionReceiptParams{
				TransactionHash: txHash,
			})
			require.NoError(t, err)
			contractAddress = receiptResult.Receipt.ContractAddress
			require.NotEmpty(t, contractAddress)
		})

		t.Run("EthCall", func(t *testing.T) {
			require.NotEmpty(t, contractAddress, "need contract address to call")

			packed, _, err := abi.EncodeFunctionCall(string(rpc.Abi_HelloWorld), "Hello", logger)
			require.NoError(t, err)

			result, err := eth.EthCall(&web3.EthCallParams{
				Transaction: web3.Transaction{
					From: x.EncodeBytes(genesisAccounts[1].GetAddress().Bytes()),
					To:   contractAddress,
					Data: x.EncodeBytes(packed),
				},
			})
			require.NoError(t, err)

			value, err := x.DecodeToBytes(result.ReturnValue)
			require.NoError(t, err)
			vars, err := abi.DecodeFunctionReturn(string(rpc.Abi_HelloWorld), "Hello", value)
			require.NoError(t, err)
			require.Len(t, vars, 1)
			require.Equal(t, "Hello, World", vars[0].Value)
		})

		t.Run("EthGetCode", func(t *testing.T) {
			require.NotEmpty(t, contractAddress, "need contract address get code")
			result, err := eth.EthGetCode(&web3.EthGetCodeParams{
				Address: contractAddress,
			})
			require.NoError(t, err)
			require.Equal(t, x.EncodeBytes(rpc.DeployedBytecode_HelloWorld), strings.ToLower(result.Bytes))
		})
	})

	t.Run("EthMining", func(t *testing.T) {
		result, err := eth.EthMining()
		require.NoError(t, err)
		require.Equal(t, true, result.Mining)
	})

	t.Run("EthAccounts", func(t *testing.T) {
		result, err := eth.EthAccounts()
		require.NoError(t, err)
		require.Len(t, result.Addresses, len(genesisAccounts)-1)
		for _, acc := range genesisAccounts {
			if acc.PrivateKey().CurveType == crypto.CurveTypeSecp256k1 {
				require.Contains(t, result.Addresses, x.EncodeBytes(acc.GetAddress().Bytes()))
			}
		}
	})

	t.Run("EthSign", func(t *testing.T) {
		result, err := eth.EthSign(&web3.EthSignParams{
			Address: "0x2c2d14a9a3f0d078ac8b38e3043d78ca8bc11029",
			Bytes:   "0xdeadbeaf",
		})
		require.NoError(t, err)
		require.Equal(t, `0x30440220345d17225ac03a575f467cea3a8d5cc2dea42fc89030c42ea175fd5140c542eb02200307004fc21ea592ce5ca013705959292c2de85b71d0fa0c84ebd8b541f505d5`, result.Signature)
	})

	t.Run("EthGetBlock", func(t *testing.T) {
		numberResult, err := eth.EthGetBlockByNumber(&web3.EthGetBlockByNumberParams{BlockNumber: x.EncodeNumber(1)})
		require.NoError(t, err)
		hashResult, err := eth.EthGetBlockByHash(&web3.EthGetBlockByHashParams{BlockHash: numberResult.GetBlockByNumberResult.Hash})
		require.NoError(t, err)
		require.Equal(t, numberResult.GetBlockByNumberResult, hashResult.GetBlockByHashResult)
	})

}
