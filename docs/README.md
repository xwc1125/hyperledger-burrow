# [Hyperledger Burrow](https://hyperledger.github.io/burrow)

[![CI](https://github.com/hyperledger/burrow/workflows/master/badge.svg)](https://launch-editor.github.com/actions?workflowID=master&event=push&nwo=hyperledger%2Fburrow)
[![version](https://img.shields.io/github/tag/hyperledger/burrow.svg)](https://github.com/hyperledger/burrow/releases/latest)
[![GoDoc](https://godoc.org/github.com/burrow?status.png)](https://godoc.org/github.com/hyperledger/burrow)
[![license](https://img.shields.io/github/license/hyperledger/burrow.svg)](../LICENSE.md)
[![LoC](https://tokei.rs/b1/github/hyperledger/burrow?category=lines)](https://github.com/hyperledger/burrow)
[![codecov](https://codecov.io/gh/hyperledger/burrow/branch/master/graph/badge.svg)](https://codecov.io/gh/hyperledger/burrow)

Hyperledger Burrow is a permissioned Ethereum smart-contract blockchain node. It executes Ethereum EVM and WASM smart contract code (usually written in [Solidity](https://solidity.readthedocs.io)) on a permissioned virtual machine. Burrow provides transaction finality and high transaction throughput on a proof-of-stake [Tendermint](https://tendermint.com) consensus engine.

![burrow logo](assets/burrow.png)

## What is Burrow

Burrow is a fully fledged blockchain node and smart contract execution engine -- a distributed database that executes code. Burrow runs Ethereum Virtual Machine (EVM) and Web Assembly (WASM) smart contracts. Burrow networks are synchronised using the [Tendermint](https://github.com/tendermint/tendermint) consensus algorithm.

Highlights include:

- **Tamper-resistant merkle state** - a node can detect if its state is corrupted or if a validator is dishonestly executing the protocol.
- **Proof-of-stake support** - run a private or public permissioned network.
- **On-chain governance primitives** - stakeholders may vote for autonomous smart contract upgrades.
- **Ethereum account world-view** - state and code is organised into cryptographically-addressable accounts.
- **Low-level permissioning** - code execution permissions can be set on a per-account basis.
- **Event streaming** - application state is organised in an event stream and can drive external systems.
- **[SQL mapping layer](reference/vent.md)** - map smart contract event emissions to SQL tables using a projection specification.
- **GRPC interfaces** - all RPC calls can be accessed from any language with GRPC support. Protobuf is used extensively for core objects.
- **Javascript client library** - client library uses code generation to provide access to contracts via statically Typescript objects.
- **Keys service** - provides optional delegating signing at the server or via a local proxy
- **Web3 RPC** - provides compatibility for mainnet Ethereum tooling such as Truffle and Metamask

### What it is not

- An Ethereum mainnet client - we do not speak devp2p and various implementation details are different. We are Ethereum-derived but exploit greater freedom than mainnet compatibility would allow.
- Just a virtual machine
- A research project - we run it in production

### Further introductory material

See [Burrow - the Boring Blockchain](https://wiki.hyperledger.org/display/burrow/Burrow+-+The+Boring+Blockchain) for an introduction to Burrow and its motivating vision.

Watch the [Boring into Burrow](https://www.youtube.com/watch?v=OpbjYaGAP4k) talk from the Hyperledger Global Forum 2020

## JavaScript Client

There is a [JavaScript API](https://github.com/hyperledger/burrow/tree/master/js)

## Project Roadmap

Project information generally updated on a quarterly basis can be found on the [Hyperledger Burrow Wiki](https://wiki.hyperledger.org/display/burrow).

## Documentation
Burrow getting started documentation is available on the [documentation site](https://hyperledger.github.io/burrow) (source markdown files can be found in [docs]()) and programmatic API in [GoDocs](https://godoc.org/github.com/hyperledger/burrow).

## Contribute

We welcome any and all contributions. Read the [contributing file](../.github/CONTRIBUTING.md) for more information on making your first Pull Request to Burrow!

You can find us on:
- [Hyperledger Chat](https://chat.hyperledger.org)
- [Hyperledger Mailing List](https://lists.hyperledger.org/mailman/listinfo)
- [here on Github](https://github.com/hyperledger/burrow/issues)

## License

[Apache 2.0](../LICENSE.md)
