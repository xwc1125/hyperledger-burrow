// Copyright Monax Industries Limited
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/hyperledger/burrow/execution/native"
	"github.com/iancoleman/strcase"
)

const contractTemplateText = `pragma solidity [[.SolidityPragmaVersion]];

/**
[[.Comment]]
* @dev These functions can be accessed as if this contract were deployed at a special address ([[.Address]]).
* @dev This special address is defined as the last 20 bytes of the sha3 hash of the the contract name.
* @dev To instantiate the contract use:
* @dev [[.Name]] [[.InstanceName]] = [[.Name]](address(uint256(keccak256("[[.Name]]"))));
*/
interface [[.Name]] {[[range .Functions]]
[[.SolidityIndent 1]]
[[end]]}
`
const functionTemplateText = `/**
[[.Comment]]
*/
function [[.Name]]([[.ArgList]]) external returns ([[.RetList]]);`

// Solidity style guide recommends 4 spaces per indentation level
// (see: http://solidity.readthedocs.io/en/develop/style-guide.html)
const indentString = "    "

var contractTemplate *template.Template
var functionTemplate *template.Template

func init() {
	var err error
	functionTemplate, err = template.New("SolidityFunctionTemplate").
		Delims("[[", "]]").
		Parse(functionTemplateText)
	if err != nil {
		panic(fmt.Errorf("couldn't parse native function template: %s", err))
	}
	contractTemplate, err = template.New("SolidityContractTemplate").
		Delims("[[", "]]").
		Parse(contractTemplateText)
	if err != nil {
		panic(fmt.Errorf("couldn't parse native contract template: %s", err))
	}
}

type solidityContract struct {
	SolidityPragmaVersion string
	*native.Contract
}

type solidityFunction struct {
	*native.Function
}

//
// Contract
//

// Create a templated solidityContract from an native contract description
func NewSolidityContract(contract *native.Contract) *solidityContract {
	return &solidityContract{
		SolidityPragmaVersion: ">=0.4.24",
		Contract:              contract,
	}
}

func (contract *solidityContract) Comment() string {
	return comment(contract.Contract.Comment)
}

// Get a version of the contract name to be used for an instance of the contract
func (contract *solidityContract) InstanceName() string {
	// Hopefully the contract name is UpperCamelCase. If it's not, oh well, this
	// is meant to be illustrative rather than cast iron compilable
	instanceName := strings.ToLower(contract.Name[:1]) + contract.Name[1:]
	if instanceName == contract.Name {
		return "contractInstance"
	}
	return instanceName
}

func (contract *solidityContract) Address() string {
	return fmt.Sprintf("0x%s",
		contract.Contract.Address())
}

// Generate Solidity code for this native contract
func (contract *solidityContract) Solidity() (string, error) {
	buf := new(bytes.Buffer)
	err := contractTemplate.Execute(buf, contract)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (contract *solidityContract) Functions() []*solidityFunction {
	functions := contract.Contract.Functions()
	solidityFunctions := make([]*solidityFunction, len(functions))
	for i, function := range functions {
		solidityFunctions[i] = NewSolidityFunction(function)
	}
	return solidityFunctions
}

//
// Function
//

// Create a templated solidityFunction from an native function description
func NewSolidityFunction(function *native.Function) *solidityFunction {
	return &solidityFunction{function}
}

func (function *solidityFunction) ArgList() string {
	abi := function.Abi()
	argList := make([]string, len(abi.Inputs))
	for i, arg := range abi.Inputs {
		storage := ""
		if arg.EVM.Dynamic() {
			storage = " calldata"
		}
		argList[i] = fmt.Sprintf("%s%s %s", arg.EVM.GetSignature(), storage, param(arg.Name))
	}
	return strings.Join(argList, ", ")
}

func (function *solidityFunction) RetList() string {
	abi := function.Abi()
	argList := make([]string, len(abi.Outputs))
	for i, arg := range abi.Outputs {
		argList[i] = fmt.Sprintf("%s %s", arg.EVM.GetSignature(), param(arg.Name))
	}
	return strings.Join(argList, ", ")
}

func (function *solidityFunction) Comment() string {
	return comment(function.Function.Comment)
}

func (function *solidityFunction) SolidityIndent(indentLevel uint) (string, error) {
	return function.solidity(indentLevel)
}

func (function *solidityFunction) Solidity() (string, error) {
	return function.solidity(0)
}

func (function *solidityFunction) solidity(indentLevel uint) (string, error) {
	buf := new(bytes.Buffer)
	iw := NewIndentWriter(indentLevel, indentString, buf)
	err := functionTemplate.Execute(iw, function)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

//
// Utility
//

func comment(comment string) string {
	commentLines := make([]string, 0, 5)
	for _, line := range strings.Split(comment, "\n") {
		trimLine := strings.TrimLeft(line, " \t\n")
		if trimLine != "" {
			commentLines = append(commentLines, trimLine)
		}
	}
	return strings.Join(commentLines, "\n")
}

func param(name string) string {
	return "_" + strcase.ToSnake(name)
}
