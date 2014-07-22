// The multiwallet tool spawns a btcwallet server process for each
// series in the voting pool.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/monetas/btcutil"
)

func main() {
	// Parse which network we're using.
	var simnet bool
	var mainnet bool
	var closewallet bool
	flag.BoolVar(&simnet, "simnet", false, "connect to simnet")
	flag.BoolVar(&mainnet, "mainnet", false, "connect to mainnet")
	flag.BoolVar(&closewallet, "closewallet", false, "close wallet processes")
	flag.Parse()

	// Get the root cert for connecting to secure websocket.
	btcwalletHomeDir := btcutil.AppDataDir("btcwallet", false)

	if closewallet {
		// Close all the btcwallet processes.
		files, _ := ioutil.ReadDir(btcwalletHomeDir)
		for _, file := range files {
			port, err := strconv.Atoi(file.Name())
			if err != nil || port < 8400 || port > 28409 {
				continue
			}
			// Get contents of pid file.
			pid, err := ioutil.ReadFile(filepath.Join(btcwalletHomeDir, file.Name(), "pid"))
			if err != nil {
				fmt.Printf("closewallet get %v pid: %v\n", port, err)
				continue
			}
			_, err = exec.Command("kill", string(pid)).Output()
			if err != nil {
				fmt.Printf("closewallet kill %v: %v\n", port, err)
			}
		}
		return
	}

	// The starting port used by the btcwallet processes are determined as follows:
	//  mainnet: 8400
	//  testnet: 18400 (default)
	//  simnet : 28400
	var startport int
	if mainnet {
		startport = 8400
	} else if simnet {
		startport = 28400
	} else {
		startport = 18400
	}

	// TODO: check if btcd and btcwallet are installed.

	// Start up btcd.
	path := os.Getenv("GOPATH")
	btcd := filepath.Join(path, "bin", "btcd")
	var cmd *exec.Cmd
	if mainnet {
		cmd = exec.Command(btcd)
	} else if simnet {
		cmd = exec.Command(btcd, "--simnet")
	} else {
		cmd = exec.Command(btcd, "--testnet")
	}

	err := cmd.Start()
	if err != nil {
		log.Fatalf("btcd start: %v", err)
	}

	// Start up 10 btcwallet processes.
	btcwallet := filepath.Join(path, "bin", "btcwallet")
	btcctl := filepath.Join(path, "bin", "btcctl")
	for i := 0; i < 10; i++ {
		port := startport + i
		listen := fmt.Sprintf("--rpclisten=127.0.0.1:%v", port)
		dir := filepath.Join(btcwalletHomeDir, fmt.Sprintf("%v", port))
		// Directory creation or the pid file creation will fail.
		err = os.Mkdir(dir, os.ModeDir|0700)
		if e, ok := err.(*os.PathError); ok && e.Err != syscall.EEXIST {
			log.Fatalf("home dir creation: %v", err)
		}

		// Check that the right cert files exist.
		cert := filepath.Join(dir, "rpc.cert")
		key := filepath.Join(dir, "rpc.key")
		_, err1 := os.Stat(cert)
		_, err2 := os.Stat(key)
		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			cert_s := filepath.Join(btcwalletHomeDir, "rpc.cert")
			key_s := filepath.Join(btcwalletHomeDir, "rpc.key")
			exec.Command("ln", "-s", cert_s, cert).Output()
			exec.Command("ln", "-s", key_s, key).Output()
		}

		datadir := fmt.Sprintf("--datadir=%v", dir)
		cmd = exec.Command(btcwallet, "--username=user", "--password=pass", listen, datadir)
		err = cmd.Start()
		if err != nil {
			log.Fatalf("%v btcwallet start: %v", port, err)
		}

		// Record the process id.
		file, err := os.Create(filepath.Join(dir, "pid"))
		if err != nil {
			log.Fatalf("pid file creation: %v", err)
		}
		pid := fmt.Sprintf("%v", cmd.Process.Pid)
		file.WriteString(pid)
		file.Close()

		// Create an encrypted wallet.
		server := fmt.Sprintf("--rpcserver=localhost:%v", port)
		_, err = exec.Command(btcctl, server, "createencryptedwallet", "test").Output()
		if err != nil {
			fmt.Printf("createencryptedwallet %v: %v\n", port, err)
		}
	}
}
