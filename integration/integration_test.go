package websocket_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateSeries(t *testing.T) {
	path := os.Getenv("GOPATH")
	btcd := filepath.Join(path, "bin", "btcd")
	btcwallet := filepath.Join(path, "bin", "btcwallet")
	btcctl := filepath.Join(path, "bin", "btcctl")

	file, err := ioutil.TempDir("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create db file: %v", err)
	}
	os.Remove(file)

	datadir := fmt.Sprintf("--datadir=%v", file+"btcd")

	// start a new simnet btcd
	btcdCmd := exec.Command(btcd, "--simnet", "--rpcuser=user", "--rpcpass=pass", datadir)

	err = btcdCmd.Start()
	if err != nil {
		t.Fatalf("btcd failed to start", err)
	}

	// start a new simnet btcwallet
	datadir = fmt.Sprintf("--datadir=%v", file+"btcwallet")
	btcwalletCmd := exec.Command(btcwallet, "--simnet", "--username=user", "--password=pass", datadir)
	err = btcwalletCmd.Start()
	if err != nil {
		t.Fatalf("btcwallet failed to start", err)
	}

	time.Sleep(time.Second)

	// create a new encrypted wallet
	createWallet := exec.Command(btcctl, "--simnet", "--rpcuser=user", "--rpcpass=pass", "--wallet", "createencryptedwallet", "test")
	err = createWallet.Run()
	if err != nil {
		t.Fatalf("create wallet failed to run", err)
	}
}
