/*

This tests the "getdepositscript" command via websocket. Assumes btcwallet/btcd are running in simnet mode with rpc enabled on both.

DO NOT RUN more than 25 times without restarting btcwallet. There is a 25 connection limit that imposes itself.

*/

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetDepositScript(t *testing.T) {
	// create the command for the websocket connection
	path := os.Getenv("GOPATH")
	binary := filepath.Join(path, "bin", "websocket")
	json := "{\"jsonrpc\":\"1.0\",\"id\":7,\"method\":\"getdepositscript\",\"params\":[]}"
	cmd := exec.Command(binary, "-simnet", json)

	// grab the std out pipe for the command
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Error(err)
	}

	// start the command
	err = cmd.Start()
	if err != nil {
		t.Error(err)
	}

	results1 := make([]byte, 69)
	results2 := make([]byte, 1000)

	// first message should be the connected message
	stdout.Read(results1)
	r := string(results1)
	want := "{\"jsonrpc\":\"1.0\",\"id\":null,\"method\":\"btcdconnected\",\"params\":[true]}\n"
	if r != want {
		t.Error(r)
	}

	// second message should be the expected stub function
	stdout.Read(results2)
	r = string(results2[:61])
	want = "{\"result\":\"someasyetunimplementedscript\",\"error\":null,\"id\":7}"

	if r != want {
		t.Error(string(results2))
	}
}
