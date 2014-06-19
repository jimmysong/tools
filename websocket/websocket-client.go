/*
Websocket Client is a utility to connect to btcwallet using websockets

*/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/monetas/btcutil"
	"github.com/monetas/websocket"
)

type T struct {
	Msg   string
	Count int
}

func main() {
	// message is the JSON to be sent to the websocket connection
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Println("Usage: websocket <JSON to send to btcwallet websocket server>")
		return
	}
	message := []byte(arguments[0])

	// get the root cert for connecting to secure websocket
	btcwalletHomeDir := btcutil.AppDataDir("btcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcwalletHomeDir, "rpc.cert"))

	if err != nil {
		log.Fatal(err)
	}
	// Setup TLS
	var tlsConfig *tls.Config
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certs)
	tlsConfig = &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}

	// Create a websocket dialer that will be used to make the connection.
	dialer := websocket.Dialer{TLSClientConfig: tlsConfig}

	// The RPC server requires basic authorization, so create a custom
	// request header with the Authorization header set.
	login := "user:pass"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login))
	requestHeader := make(http.Header)
	requestHeader.Add("Authorization", auth)

	// Dial the connection.
	url := "wss://localhost:18332/frontend"
	conn, resp, err := dialer.Dial(url, requestHeader)

	if err != nil {
		log.Println(resp)
		log.Fatal(err)
	}

	// send message to websocket connection.
	conn.WriteMessage(websocket.TextMessage, message)

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Fatal(err)
		}

		m := string(msg)

		log.Println(m)
	}
}
