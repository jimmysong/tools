// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/monetas/btcnet"
	"github.com/monetas/btcscript"
	"github.com/monetas/btcutil"
	"github.com/monetas/btcutil/hdkeychain"
)

// This example demonstrates how to generate a cryptographically random seed
// then use it to create a new master node (extended key).
func main() {

	var keys []*hdkeychain.ExtendedKey
	var pks [][]*btcutil.AddressPubKey
	keys = make([]*hdkeychain.ExtendedKey, 3, 3)
	pks = make([][]*btcutil.AddressPubKey, 3, 3)

	params := btcnet.MainNetParams

	for i := range keys {
		pks[i] = make([]*btcutil.AddressPubKey,3,3)
	}

	for i := range keys {
		seed, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		key, _ := hdkeychain.NewMaster(seed)
		
		keys[i], _ = key.Neuter()
		fmt.Println("Public Extended Key:", i, keys[i].String())
		for j := range keys {
			child, _ := keys[i].Child(uint32(j))
			pubkey, _ := child.ECPubKey()
			fmt.Println("Child ", j, hex.EncodeToString(pubkey.SerializeCompressed()))
			pks[j][i], _ = btcutil.NewAddressPubKey(pubkey.SerializeCompressed(), &params)
		}
	}

	for i := range keys {
		script, _ := btcscript.MultiSigScript(pks[i], 2)

		addr, _ := btcutil.NewAddressScriptHash(script, &params)

		fmt.Println("deposit script:", i, addr.EncodeAddress())
	}

	// Show that the generated master node extended key is private.

	// Output:
	// Private Extended Key?: true
}
