#!/bin/bash

export this="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$this/../tests/test_runner.sh"

test_setup
trap test_teardown EXIT

cd $this
export SIGNING_ADDRESS="$key1_addr"
export BURROW_URL="$BURROW_HOST:$BURROW_GRPC_PORT"

mocha --bail --exit -r ts-node/register "$1"
test_exit=$?
