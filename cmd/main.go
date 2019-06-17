package main

import (
	"bufio"
	"fmt"
	"os"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	abcicli "github.com/tendermint/tendermint/abci/client"
	rpccli "github.com/tendermint/tendermint/rpc/client"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	tmapp "github.com/vocdoni/go-dvote-chain/app"
)

var (
	client abcicli.Client
	logger log.Logger
)

var app *tmapp.CounterApplication
var tmRPC rpccli.Client

func startABCI() {

	flagAddress := "0.0.0.0:26658"
	flagSerial := false
	flagAbci := "socket"

	app = tmapp.NewCounterApplication(flagSerial)
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
	tmRPC = rpccli.NewHTTP("http://localhost:26657", "/websocket")
	reader := bufio.NewReader(os.Stdin)
	var reqQuery types.RequestQuery

	for {
		fmt.Println("---------------------")
		fmt.Print("tx: ")
		tx, _ := reader.ReadString('\n')
		res, err := tmRPC.BroadcastTxCommit([]byte(tx))
		if err != nil {
			fmt.Println(err)
		}
		if res.CheckTx.IsErr() {
			fmt.Println(res.CheckTx.Log)
		}
		if res.DeliverTx.IsErr() {
			fmt.Println(res.DeliverTx.Log)
		}

		reqQuery.Path = "hash"
		fmt.Printf("Hashes: %s\n", app.Query(reqQuery).Value)

		reqQuery.Path = "tx"
		fmt.Printf("TXs: %s\n", app.Query(reqQuery).Value)
	}
}
