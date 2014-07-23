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

func startBtcd(cmd *exec.Cmd, homeDir string) {
	err := cmd.Start()
	if err != nil {
		log.Fatalf("btcd start: %v", err)
	}
	// Record the process id.
	file, err := os.Create(filepath.Join(homeDir, "pid"))
	if err != nil {
		log.Fatalf("btcd pid file creation: %v", err)
	}
	pid := fmt.Sprintf("%v", cmd.Process.Pid)
	file.WriteString(pid)
	file.Close()
}

func stopBtcd(homeDir string) {
	// Get contents of pid file.
	pid, err := ioutil.ReadFile(filepath.Join(homeDir, "pid"))
	if err != nil {
		fmt.Printf("stopBtcd get pid: %v\n", err)
		return
	}
	_, err = exec.Command("kill", string(pid)).Output()
	if err != nil {
		fmt.Printf("stopBtcd kill: %v\n", err)
	}
}

func startAllWallets(gopath string, startport int, btcwalletHomeDir string) {
	// Start up 10 btcwallet processes.
	btcwallet := filepath.Join(gopath, "bin", "btcwallet")
	btcctl := filepath.Join(gopath, "bin", "btcctl")
	for i := 0; i < 10; i++ {
		port := startport + i
		listen := fmt.Sprintf("--rpclisten=127.0.0.1:%v", port)
		dir := filepath.Join(btcwalletHomeDir, fmt.Sprintf("%v", port))
		// Make the wallet directory beforehand or the pid file creation will fail.
		err := os.Mkdir(dir, os.ModeDir|0700)
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
		cmd := exec.Command(btcwallet, "--username=user", "--password=pass", listen, datadir)
		err = cmd.Start()
		if err != nil {
			log.Fatalf("%v btcwallet start: %v", port, err)
		}

		// Record the process id.
		file, err := os.Create(filepath.Join(dir, "pid"))
		if err != nil {
			log.Fatalf("btcwallet pid file creation: %v", err)
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

func stopAllWallets(btcwalletHomeDir string) {
	files, _ := ioutil.ReadDir(btcwalletHomeDir)
	for _, file := range files {
		port, err := strconv.Atoi(file.Name())
		if err != nil || port < 8400 || port > 28409 {
			continue
		}
		// Get contents of pid file.
		pid, err := ioutil.ReadFile(filepath.Join(btcwalletHomeDir, file.Name(), "pid"))
		if err != nil {
			fmt.Printf("stopAllWallets get %v pid: %v\n", port, err)
			continue
		}
		_, err = exec.Command("kill", string(pid)).Output()
		if err != nil {
			fmt.Printf("stopAllWallets kill %v: %v\n", port, err)
		}
	}
}

func main() {
	var mainnet bool
	flag.BoolVar(&mainnet, "mainnet", false, "connect to mainnet")
	var simnet bool
	flag.BoolVar(&simnet, "simnet", false, "connect to simnet")
	var stopall bool
	flag.BoolVar(&stopall, "stopall", false, "stop btcd and wallets")
	var stopwallets bool
	flag.BoolVar(&stopwallets, "stopwallets", false, "stop wallet processes")
	flag.Parse()

	gopath := os.Getenv("GOPATH")
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	btcwalletHomeDir := btcutil.AppDataDir("btcwallet", false)

	if stopall {
		stopBtcd(btcdHomeDir)
		stopAllWallets(btcwalletHomeDir)
		return
	}

	if stopwallets {
		stopAllWallets(btcwalletHomeDir)
		return
	}

	btcdPath := filepath.Join(gopath, "bin", "btcd")
	var btcdCmd *exec.Cmd
	if mainnet {
		btcdCmd = exec.Command(btcdPath)
	} else if simnet {
		btcdCmd = exec.Command(btcdPath, "--simnet")
	} else {
		btcdCmd = exec.Command(btcdPath, "--testnet")
	}
	// TODO: check if btcd and btcwallet are installed.
	startBtcd(btcdCmd, btcdHomeDir)

	// The starting port used by the btcwallet processes is:
	// mainnet: 8400 - testnet: 18400 (default) - simnet : 28400
	var startport int
	if mainnet {
		startport = 8400
	} else if simnet {
		startport = 28400
	} else {
		startport = 18400
	}
	startAllWallets(gopath, startport, btcwalletHomeDir)
}
