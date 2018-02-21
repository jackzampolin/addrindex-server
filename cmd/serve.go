// Copyright © 2018 Jack Zampolin <jack@blockstack.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackzampolin/addrindex-server/addrindex"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the addrindex server",
	Run: func(cmd *cobra.Command, args []string) {
		as := addrindex.NewAddrServer(cfg)
		defer as.Client.Shutdown()

		router := mux.NewRouter()

		router.HandleFunc("/addr/{addr}/utxo", as.HandleAddrUTXO).Methods("GET")
		router.HandleFunc("/addr/{addr}/balance", as.HandleAddrBalance).Methods("GET")
		router.HandleFunc("/addr/{addr}/totalReceived", as.HandleAddrRecieved).Methods("GET")
		router.HandleFunc("/addr/{addr}/totalSent", as.HandleAddrSent).Methods("GET")
		router.HandleFunc("/tx/{txid}", as.HandleTxGet).Methods("GET")
		router.HandleFunc("/rawtx/{txid}", as.HandleRawTxGet).Methods("GET")
		router.HandleFunc("/messages/verify", as.HandleMessagesVerify).Methods("POST")
		router.HandleFunc("/tx/send", as.HandleTransactionSend).Methods("POST")
		router.HandleFunc("/block/{blockHash}", as.HandleGetBlock).Methods("GET")
		router.HandleFunc("/block-index/{height}", as.HandleGetBlockHash).Methods("GET")
		router.HandleFunc("/status", as.GetStatus).Methods("GET")
		router.HandleFunc("/sync", as.GetSync).Methods("GET")
		router.HandleFunc("/txs", as.GetTransactions).Methods("GET")
		router.HandleFunc("/version", as.GetVersion).Methods("GET")
		router.HandleFunc("/test/{addr}", as.HandleTest).Methods("GET")

		// router.HandleFunc("/addr/{addr}/unconfirmedBalance", as.HandleAddrUnconfirmed).Methods("GET")

		// /insight-api/blocks?limit=3&blockDate=2016-04-22
		// NOTE: this should fetch the last n blocks
		// router.HandleFunc("/blocks", as.HandleGetBlocks).Methods("GET")

		// NOTE: This pulls data from outside price APIs. Might want to implement a couple
		// GET /currency
		// router.HandleFunc("/currency", as.GetCurrency).Methods("GET")

		log.Println(fmt.Sprintf("Listening on port ':%v'...", as.Port))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", as.Port), handlers.LoggingHandler(os.Stdout, router)))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
