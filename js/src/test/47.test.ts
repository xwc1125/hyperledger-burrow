import * as assert from 'assert';
import {burrow, compile} from '../test';

describe('#47', function () {it('#47', async () => {
    const source = `
      pragma solidity >=0.0.0;
      contract Test{
        string _withSpace = "  Pieter";
        string _withoutSpace = "Pieter";

        function getWithSpaceConstant() public view returns (string memory) {
          return _withSpace;
        }

        function getWithoutSpaceConstant () public view returns (string memory) {
          return _withoutSpace;
        }
      }
    `
    const {abi, code} = compile(source, 'Test')
    return burrow.contracts.deploy(abi, code)
      .then((contract: any) => Promise.all([contract.getWithSpaceConstant(), contract.getWithoutSpaceConstant()]))
      .then(([withSpace, withoutSpace]) => {
        assert.deepStrictEqual(withSpace, ['  Pieter'])
        assert.deepStrictEqual(withoutSpace, ['Pieter'])
      })
  })
})
