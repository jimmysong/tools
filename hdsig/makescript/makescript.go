// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/monetas/btcnet"
	"github.com/monetas/btcscript"
	"github.com/monetas/btcutil"
	"github.com/monetas/btcutil/hdkeychain"
)

func ShowUsage() {
	fmt.Println("makescript -keys=<comma-separated list of extended public keys> -num=<number of keys required to spend")
	os.Exit(1)
}

func GetChild(key *hdkeychain.ExtendedKey, path string) *hdkeychain.ExtendedKey {
	pathcomponents := strings.Split(path, "/")
	current := key
	for _, pc := range pathcomponents {
		var isHardened bool = false
		if strings.HasSuffix(pc, "+") {
			pc = strings.TrimSuffix(pc, "+")
			isHardened = true
		}
		childnum, err := strconv.Atoi(pc)
		if err != nil {
			fmt.Printf("Illegal path component: %v\n", pc)
			os.Exit(1)
		}

		if childnum >= hdkeychain.HardenedKeyStart {
			fmt.Printf("Number %d not in range [0:2^31-1]\n", childnum)
			os.Exit(1)
		}

		childNum32 := uint32(childnum)
		if isHardened {
			childNum32 += hdkeychain.HardenedKeyStart
		}

		current, err = current.Child(childNum32)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return current
}

// This example demonstrates how to generate a cryptographically random seed
// then use it to create a new master node (extended key).
func main() {

	// for creating new keys
	var path, rawkeys, net string
	var m int
	flag.StringVar(&rawkeys, "keys", "", "Public HD Keys to generate a deposit script")
	flag.IntVar(&m, "num", 1, "Number of keys required to spend")
	flag.StringVar(&path, "path", "0", "Child key to derive from each hd key")
	flag.StringVar(&net, "net", "main", "Choice of network: main, test, sim")
	flag.Parse()

	if rawkeys == "" {
		ShowUsage()
	}

	keystrings := strings.Split(rawkeys, ",")
	n := len(keystrings)
	var keys []*hdkeychain.ExtendedKey
	keys = make([]*hdkeychain.ExtendedKey, n, n)

	for i, keystring := range keystrings {
		keys[i], _ = hdkeychain.NewKeyFromString(keystring)
	}
	

	var pks []*btcutil.AddressPubKey
	pks = make([]*btcutil.AddressPubKey, n, n)

	var params btcnet.Params
	switch net {
	case "sim":
		params = btcnet.SimNetParams
	case "test":
		params = btcnet.TestNet3Params
	case "main":
		params = btcnet.MainNetParams
	}

	for i := range keys {
		child := GetChild(keys[i], path)
		pubkey, _ := child.ECPubKey()
		fmt.Println("Key", i, "Child", path, hex.EncodeToString(pubkey.SerializeCompressed()))
		pks[i], _ = btcutil.NewAddressPubKey(pubkey.SerializeCompressed(), &params)
	}

	script, _ := btcscript.MultiSigScript(pks, m)
	addr, _ := btcutil.NewAddressScriptHash(script, &params)
	fmt.Println("deposit script:", 0, addr.EncodeAddress())

}
