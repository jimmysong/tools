// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/monetas/btcutil/hdkeychain"
)

func ShowUsage() {
	fmt.Println("Usage: create <number of keys to create>")
	os.Exit(1)
}

// This example demonstrates how to generate a cryptographically random seed
// then use it to create a new master node (extended key).
func main() {

	flag.Parse()
	arguments := flag.Args()
	if len(arguments) != 1 {
		ShowUsage()
	}
	num, err := strconv.Atoi(arguments[0])
	if err != nil {
		fmt.Println(err)
		ShowUsage()
	}

	for i := 0; i < num; i++ {
		seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			fmt.Println(err)
			return
		}
		key, err := hdkeychain.NewMaster(seed)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("-------------------------------")
		fmt.Println("Private Extended Key:", key.String())
		pubkey, _ := key.Neuter()
		fmt.Println("Public Extended Key:", pubkey.String())
	}
}
