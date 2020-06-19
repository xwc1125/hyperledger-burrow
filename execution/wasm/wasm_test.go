package wasm

import (
	"fmt"
	"testing"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/evm/abi"

	"github.com/hyperledger/burrow/crypto"
	"github.com/stretchr/testify/require"
)

func TestStaticCallWithValue(t *testing.T) {
	cache := acmstate.NewMemoryState()

	gas := uint64(0)

	params := engine.CallParams{
		Origin: crypto.ZeroAddress,
		Caller: crypto.ZeroAddress,
		Callee: crypto.ZeroAddress,
		Input:  []byte{},
		Value:  0,
		Gas:    &gas,
	}

	// run constructor
	runtime, cerr := RunWASM(cache, params, Bytecode_storage_test)
	require.NoError(t, cerr)

	// run getFooPlus2
	spec, err := abi.ReadSpec(Abi_storage_test)
	require.NoError(t, err)
	calldata, _, err := spec.Pack("getFooPlus2")

	params.Input = calldata

	returndata, cerr := RunWASM(cache, params, runtime)
	require.NoError(t, cerr)

	data := abi.GetPackingTypes(spec.Functions["getFooPlus2"].Outputs)

	err = spec.Unpack(returndata, "getFooPlus2", data...)
	require.NoError(t, err)
	returnValue := *data[0].(*uint64)
	var expected uint64
	expected = 104
	require.Equal(t, expected, returnValue)

	// call incFoo
	calldata, _, err = spec.Pack("incFoo")

	params.Input = calldata

	returndata, cerr = RunWASM(cache, params, runtime)
	require.NoError(t, cerr)

	require.Equal(t, returndata, []byte{})

	// run getFooPlus2
	calldata, _, err = spec.Pack("getFooPlus2")
	require.NoError(t, err)

	params.Input = calldata

	returndata, cerr = RunWASM(cache, params, runtime)
	require.NoError(t, cerr)

	spec.Unpack(returndata, "getFooPlus2", data...)
	expected = 105
	returnValue = *data[0].(*uint64)

	require.Equal(t, expected, returnValue)
}

func blockHashGetter(height uint64) []byte {
	return binary.LeftPadWord256([]byte(fmt.Sprintf("block_hash_%d", height))).Bytes()
}
