#!/bin/bash
set -e

GOPATH=$HOME/go
BIN=$GOPATH/bin
BTCD0=btcd0
BTCD1=btcd1

mkdir -p $BTCD0
mkdir -p $BTCD1

cp rpc* $BTCD0
cp rpc* $BTCD1

#btcd0

$BIN/btcd --datadir=$BTCD0/ --simnet --listen=:18555 --rpcuser=user --rpcpass=pass &

#btcd1
$BIN/btcd --datadir=$BTCD1 --simnet --listen=:18655 --connect=127.0.0.1:18555 --generate  --miningaddr=SgM2hzgo2tJpsy2rHNL8VkFg7qKCJPMQWD &

#btcwallet
$BIN/btcwallet -D btcwallet -C btcwallet/btcwallet.conf --logdir=btcwallet &

#btcctl
#$BIN/btcctl -C btcctl/btcctl.conf -c rpc.cert COMMAND (i.e. getblockcount)

#to stop 
#killall -9 btcd ; killall -9 btcwallet
