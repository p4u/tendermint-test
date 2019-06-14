package main

import (
	"bufio"
	"fmt"
	"os"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	abcicli "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/vocdoni/go-dvote-chain/counter"
)

var (
	client abcicli.Client
	logger log.Logger
)

var app *counter.CounterApplication

func startABCI() {

	flagAddress := "0.0.0.0:26658"
	flagSerial := true
	flagAbci := "socket"
	//	flabLogLevel := "debug"

	app = counter.NewCounterApplication(flagSerial)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// Start the listener
	srv, err := server.NewServer(flagAddress, flagAbci, app)
	if err != nil {
		return
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		return
	}

	// Stop upon receiving SIGTERM or CTRL-C.
	cmn.TrapSignal(logger, func() {
		// Cleanup
		srv.Stop()
	})

	// Run forever.
	select {}
}

func main() {
	go startABCI()
	reader := bufio.NewReader(os.Stdin)
	var reqQuery types.RequestQuery

	for {
		fmt.Println("---------------------")
		fmt.Print("tx: ")
		tx, _ := reader.ReadString('\n')
		response := app.DeliverTx([]byte(tx))
		fmt.Printf("Response: %t\n", response.Log)
		//fmt.Printf("Comit: %s\n", app.Commit().Data)

		reqQuery.Path = "hash"
		fmt.Printf("Hashes: %s\n", app.Query(reqQuery).Value)

		reqQuery.Path = "tx"
		fmt.Printf("TXs: %s\n", app.Query(reqQuery).Value)
	}
}
