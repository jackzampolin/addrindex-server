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

package addrindex

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/gorilla/mux"
	"github.com/jackzampolin/addrindex-server/cache"
)

// AddrServer is the struct where all methods are defined
type AddrServer struct {
	Host            string
	User            string
	Pass            string
	DisableTLS      bool
	Port            int
	Client          *rpcclient.Client
	Blocks          *Blocks
	RedisConnection string

	versionData versionData
}

func (as *AddrServer) version() []byte {
	out, _ := json.Marshal(as.versionData)
	return out
}

type versionData struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
}

// AddrServerConfig configures the AddrServer
type AddrServerConfig struct {
	Host            string `json:"host"`
	Usr             string `json:"usr"`
	Pass            string `json:"pass"`
	SSL             bool   `json:"ssl"`
	Port            int    `json:"port"`
	RedisConnection string `json:"redis"`
	Version         string
	Commit          string
	Branch          string
}

// NewAddrServer returns a new AddrServer instance
func NewAddrServer(cfg *AddrServerConfig) *AddrServer {
	out := &AddrServer{
		Host:            cfg.Host,
		User:            cfg.Usr,
		Pass:            cfg.Pass,
		DisableTLS:      !cfg.SSL,
		Port:            cfg.Port,
		RedisConnection: cfg.RedisConnection,
		versionData: versionData{
			Version: cfg.Version,
			Commit:  cfg.Commit,
			Branch:  cfg.Branch,
		},
		Blocks: &Blocks{},
	}
	client, err := rpcclient.New(out.connCfg(), nil)
	if err != nil {
		panic(err)
	}
	out.Client = client
	return out
}

// NewTestAddrServer returns a new AddrServer instance
func NewTestAddrServer(cfg *AddrServerConfig) *AddrServer {
	out := &AddrServer{
		Host:       cfg.Host,
		User:       cfg.Usr,
		Pass:       cfg.Pass,
		DisableTLS: !cfg.SSL,
		Port:       cfg.Port,
		versionData: versionData{
			Version: cfg.Version,
			Commit:  cfg.Commit,
			Branch:  cfg.Branch,
		},
		Blocks: &Blocks{},
	}
	client, err := rpcclient.New(out.connCfg(), nil)
	if err != nil {
		panic(err)
	}
	out.Client = client
	return out
}

// URL returns the backend server's URL
func (as *AddrServer) URL() string {
	if as.DisableTLS {
		return fmt.Sprintf("http://%s:%s@%v", as.User, as.Pass, as.Host)
	}
	return fmt.Sprintf("https://%s:%s@%v", as.User, as.Pass, as.Host)
}

func (as *AddrServer) connCfg() *rpcclient.ConnConfig {
	return &rpcclient.ConnConfig{
		Host:         as.Host,
		User:         as.User,
		Pass:         as.Pass,
		HTTPPostMode: true,
		DisableTLS:   as.DisableTLS,
	}
}

// Router holds the routing table for the AddrServer
func (as *AddrServer) Router() *mux.Router {
	router := mux.NewRouter()
	c := cache.NewMemoryCache()
	cacheTime := "1m"

	router.HandleFunc("/addr/{addr}/utxo", as.HandleAddrUTXO).Methods("GET")
	router.HandleFunc("/addr/{addr}/balance", as.HandleAddrBalance).Methods("GET")
	router.HandleFunc("/addr/{addr}/totalReceived", as.HandleAddrRecieved).Methods("GET")
	router.HandleFunc("/addr/{addr}/totalSent", as.HandleAddrSent).Methods("GET")
	router.HandleFunc("/addr/{addr}/unconfirmedBalance", as.HandleAddrUnconfirmedBalance).Methods("GET")
	router.HandleFunc("/tx/{txid}", as.HandleTxGet).Methods("GET")
	router.HandleFunc("/txs", as.HandleGetTransactions).Methods("GET")
	router.HandleFunc("/rawtx/{txid}", as.HandleRawTxGet).Methods("GET")
	router.HandleFunc("/tx/send", as.HandleTransactionSend).Methods("POST")
	router.HandleFunc("/messages/verify", as.HandleMessagesVerify).Methods("POST")
	router.HandleFunc("/block/{blockHash}", as.HandleGetBlock).Methods("GET")
	router.HandleFunc("/blocks", cache.Middleware(cacheTime, c, as.HandleGetBlocks)).Methods("GET")
	router.HandleFunc("/block-index/{height}", as.HandleGetBlockHash).Methods("GET")
	router.HandleFunc("/status", as.HandleGetStatus).Methods("GET")
	router.HandleFunc("/sync", as.HandleGetSync).Methods("GET")
	router.HandleFunc("/version", as.HandleGetVersion).Methods("GET")
	router.HandleFunc("/currency", cache.Middleware(cacheTime, c, as.HandleGetCurrency)).Methods("GET")
	return router
}
