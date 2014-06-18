package main

import (
	"github.com/gorilla/websocket"
	"github.com/conformal/btcutil"
	"github.com/monetas/btcws"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type T struct {
	Msg string
	Count int
}

func main() {
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

	// Create getdepositscript command.
	id := 1
	cmd, err := btcws.NewGetDepositScriptCmd(id)

	// JSON marshal and send request to websocket connection.
	conn.WriteJSON(cmd)

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Fatal(err)
		}

		m := string(msg)
		
		log.Println(m)
	}
}
