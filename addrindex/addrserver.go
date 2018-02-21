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
)

// AddrServer is the struct where all methods are defined
type AddrServer struct {
	Host         string
	User         string
	Pass         string
	DisableTLS   bool
	Port         int
	Client       *rpcclient.Client
	Transactions int

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
	Host         string `json:"host"`
	Usr          string `json:"usr"`
	Pass         string `json:"pass"`
	SSL          bool   `json:"ssl"`
	Port         int    `json:"port"`
	Transactions int    `json:"transactions"`
	Version      string
	Commit       string
	Branch       string
}

// NewAddrServer returns a new AddrServer instance
func NewAddrServer(cfg *AddrServerConfig) *AddrServer {
	out := &AddrServer{
		Host:         cfg.Host,
		User:         cfg.Usr,
		Pass:         cfg.Pass,
		DisableTLS:   !cfg.SSL,
		Port:         cfg.Port,
		Transactions: cfg.Transactions,
		versionData: versionData{
			Version: cfg.Version,
			Commit:  cfg.Commit,
			Branch:  cfg.Branch,
		},
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

// fetchAllTransactions pages through the transactions for an address
// NOTE: This can take a long time!
// func (as *AddrServer) fetchAllTransactions(addr string) (Transactions, error) {
// 	txns := []Transaction{}
// 	count := 0
// 	for {
// 		// Fetch page of transactions
// 		searchtxns, err := as.SearchRawTransactions(addr, (count * as.Transactions), as.Transactions)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Increment page count
// 		count++
//
// 		// If the page is full we need to append them and continue to next page
// 		if len(searchtxns.Result) == as.Transactions {
// 			for _, tx := range searchtxns.Result {
// 				txns = append(txns, tx)
// 			}
// 			continue
// 		}
//
// 		// If not we need to save the tansactions and return
// 		for _, tx := range searchtxns.Result {
// 			txns = append(txns, tx)
// 		}
//
// 		return txns, nil
// 	}
// }
