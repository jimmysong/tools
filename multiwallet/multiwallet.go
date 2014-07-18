// The multiwallet tool spawns a btcwallet server process for each
//  series in the voting pool.

package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// parse which network we're using
	var simnet bool
	var mainnet bool
	flag.BoolVar(&simnet, "simnet", false, "connect to simnet")
	flag.BoolVar(&mainnet, "mainnet", false, "connect to mainnet")
	flag.Parse()

	// The starting port used by the btcwallet processes are determined
	//  as follows:
	//  mainnet: 8400
	//  testnet: 18400
	//  simnet : 28400
	startport := 18400
	if mainnet {
		startport = 8400
	} else if simnet {
		startport = 28400
	}

	// TODO: check if btcd and btcwallet are installed

	// start up btcd
	path := os.Getenv("GOPATH")
	btcd := filepath.Join(path, "bin", "btcd")
	cmd := exec.Command(btcd)
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	// start up 10 btcwallet processes
	btcwallet := filepath.Join(path, "bin", "btcwallet")
	btcctl := filepath.Join(path, "bin", "btcctl")
	for i := 0; i < 10; i++ {
		port := startport + i
		listen := fmt.Sprintf("--rpclisten=127.0.0.1:%v", port)
		datadir := fmt.Sprintf("--datadir=/tmp/test%v", i)
		cmd = exec.Command(btcwallet, "--username=user", "--password=pass", listen, datadir)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		// create an encrypted wallet
		server := fmt.Sprintf("--rpcserver=localhost:%v", port)
		exec.Command(btcctl, server, "createencryptedwallet", "test").Output()
		
	}
	
	

}
