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
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/gorilla/mux"
)

// HandleAddrUTXO handles the /addr/<addr>/utxo route
func (as *AddrServer) HandleAddrUTXO(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.fetchAllTransactions(addr)
	if err != nil {
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}

	utxo := txns.UTXO(addr)
	w.Write(utxo.JSON())
}

// HandleAddrBalance handles the /addr/<addr>/balance route
func (as *AddrServer) HandleAddrBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.fetchAllTransactions(addr)
	if err != nil {
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}
	utxo := txns.UTXO(addr)
	w.Write(utxo.Balance())
}

// HandleAddrRecieved handles the /addr/<addr>/totalReceived route
func (as *AddrServer) HandleAddrRecieved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.fetchAllTransactions(addr)
	if err != nil {
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}

	w.Write(txns.Received(addr))
}

// HandleAddrSent handles the /addr/<addr>/totalSent route
func (as *AddrServer) HandleAddrSent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.fetchAllTransactions(addr)
	if err != nil {
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}

	w.Write(txns.Sent(addr))
}

// HandleTxGet handles the /tx/<txid> route
func (as *AddrServer) HandleTxGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	txid := mux.Vars(r)["txid"]

	// Make the chainhash for fetching data
	hash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		w.Write(NewPostError("error parsing txhash", err))
		return
	}

	// fetch transaction details
	raw, err := as.Client.GetRawTransactionVerbose(hash)
	if err != nil {
		w.Write(NewPostError("error fetching transaction details", err))
		return
	}

	txn, _ := json.Marshal(raw)
	w.Write(txn)
}

// HandleRawTxGet handles the /rawtx/<txid> route
func (as *AddrServer) HandleRawTxGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	txid := mux.Vars(r)["txid"]

	// Make the chainhash for fetching data
	hash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		w.Write(NewPostError("error parsing txhash", err))
		return
	}

	// fetch transaction details
	raw, err := as.Client.GetRawTransactionVerbose(hash)
	if err != nil {
		w.Write(NewPostError("error fetching transaction details", err))
		return
	}

	txn, _ := json.Marshal(map[string]string{"rawtx": raw.Hex})
	w.Write(txn)
}

// HandleTransactionSend handles the /tx/send route
// TODO: Test this somehow?
func (as *AddrServer) HandleTransactionSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tx TxPost

	// Read post body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(NewPostError("unable to read post body", err))
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, tx)
	if err != nil {
		w.Write(NewPostError("unable to unmarshall body", err))
		return
	}

	// Convert tansaction to send format
	txn, err := btcutil.NewTxFromBytes([]byte(tx.Tx))
	if err != nil {
		w.Write(NewPostError("unable to parse transaction", err))
		return
	}

	ret, err := as.Client.SendRawTransaction(txn.MsgTx(), true)
	if err != nil {
		w.Write(NewPostError("unable to post transaction to node", err))
		return
	}

	out, _ := json.Marshal(ret)
	w.Write(out)
}

// HandleMessagesVerify handles the /tx/send route
// TODO: Test this somehow?
func (as *AddrServer) HandleMessagesVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tx VerifyPost

	// Read post body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(NewPostError("unable to read post body", err))
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, tx)
	if err != nil {
		w.Write(NewPostError("unable to unmarshall body", err))
		return
	}

	addr, err := btcutil.DecodeAddress(tx.BitcoinAddress, nil)
	if err != nil {
		w.Write(NewPostError("unable to decode bitcoin address", err))
		return
	}

	ret, err := as.Client.VerifyMessage(addr, tx.Signature, tx.Message)
	if err != nil {
		w.Write(NewPostError("unable verify message", err))
		return
	}

	out, _ := json.Marshal(ret)
	w.Write(out)
}

// HandleGetBlock handles the /block/<blockhash> route
func (as *AddrServer) HandleGetBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	blockhash := mux.Vars(r)["blockHash"]

	// Make the chainhash for fetching data
	hash, err := chainhash.NewHashFromStr(blockhash)
	if err != nil {
		w.Write(NewPostError("error parsing txhash", err))
		return
	}

	// paginate through transactions
	block, err := as.Client.GetBlockVerbose(hash)
	if err != nil {
		w.Write(NewPostError("error fetching block", err))
		return
	}
	out, _ := json.Marshal(block)
	w.Write(out)
}

// HandleGetBlockHash handles the /block-index/<height> route
func (as *AddrServer) HandleGetBlockHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	height := mux.Vars(r)["height"]

	h, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		w.Write(NewPostError("error parsing blockheight", err))
		return
	}

	block, err := as.Client.GetBlockHash(h)
	if err != nil {
		w.Write(NewPostError("error fetching blockhash", err))
		return
	}

	bh, _ := json.Marshal(BlockHashReturn{BlockHash: block.String()})
	w.Write(bh)
}

// GetSync handles the /sync route
func (as *AddrServer) GetSync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chainInfo, err := as.Client.GetBlockChainInfo()
	if err != nil {
		w.Write(NewPostError("error fetching blockchain info", err))
		return
	}

	w.Write(NewSyncResponse(chainInfo))
}

// GetStatus handles the /status route
func (as *AddrServer) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	method := r.URL.Query()["q"]
	switch method[0] {
	case "getDifficulty":
		info, err := as.Client.GetDifficulty()
		if err != nil {
			w.Write(NewPostError("failed to get info", err))
			return
		}
		w.Write(NewGetDifficultyReturn(info))
	case "getBestBlockHash":
		info, err := as.Client.GetBestBlockHash()
		if err != nil {
			w.Write(NewPostError("failed to get info", err))
			return
		}
		w.Write(NewGetBestBlockHashReturn(info.String()))
	default:
		info, err := as.Client.GetInfo()
		if err != nil {
			w.Write(NewPostError("failed to get info", err))
			return
		}
		out, _ := json.Marshal(info)
		w.Write(out)
	}
}

// /insight-api/txs?address=<address>
// /insight-api/txs?block=<blockhash>
// GET /txs
// router.HandleFunc("/txs", as.GetTransactions).Methods("GET")

// /insight-api/version
// GET /version
// router.HandleFunc("/version", as.GetVersion).Methods("GET")

// router.HandleFunc("/addr/{addr}/unconfirmedBalance", as.HandleAddrUnconfirmed).Methods("GET")

// NOTE: This pulls data from outside price APIs. Might want to implement a couple
// NOTE: Lets cache this data server side the same way we are doing with the block index
// GET /currency
// router.HandleFunc("/currency", as.GetCurrency).Methods("GET")
// curl https://www.bitstamp.net/api/v2/ticker/btcusd/
// curl https://blockchain.info/ticker
// curl https://api.coindesk.com/v1/bpi/currentprice/usd.json

// /insight-api/blocks?limit=3&blockDate=2016-04-22
// NOTE: We are going to need to keep a cache of this data on the server
// router.HandleFunc("/blocks", as.HandleGetBlocksByDate).Methods("GET")

// TxPost models a post request for sending a transaction
type TxPost struct {
	Tx string `json:"tx"`
}

// VerifyPost models a post request for verifying a transaction
type VerifyPost struct {
	BitcoinAddress string `json:"bitcoinaddress"`
	Signature      string `json:"signature"`
	Message        string `json:"message"`
}

// BlockHashReturn handles the return for the /block-index/<int>
type BlockHashReturn struct {
	BlockHash string `json:"blockHash"`
}

// BestBlockHashReturn handles the return for the /sync?q=getBestBlockHash
type BestBlockHashReturn struct {
	BlockHash string `json:"bestblockhash"`
}

// NewGetBestBlockHashReturn gets the bytes
func NewGetBestBlockHashReturn(blockhash string) []byte {
	out, _ := json.Marshal(BestBlockHashReturn{BlockHash: blockhash})
	return out
}

// PostError models an error returned to a client during the post
type PostError struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

// NewPostError is a convinence function for returning errors to clients
func NewPostError(msg string, err error) []byte {
	out, _ := json.Marshal(PostError{
		Message: msg,
		Error:   err,
	})
	return out
}

// SyncResponse models a response to the sync command
type SyncResponse struct {
	Status           string      `json:"status"`
	BlockChainHeight int         `json:"blockChainHeight"`
	SyncPercentage   int         `json:"syncPercentage"`
	Height           int         `json:"height"`
	Error            interface{} `json:"error"`
	Type             string      `json:"type"`
}

// GetDifficultyReturn models a response to the /status endpoint
type GetDifficultyReturn struct {
	Difficulty float64 `json:"difficulty"`
}

// NewGetDifficultyReturn gets the bytes
func NewGetDifficultyReturn(dif float64) []byte {
	out, _ := json.Marshal(GetDifficultyReturn{Difficulty: dif})
	return out
}

// NewSyncResponse returns a response for the /sync route
func NewSyncResponse(bc *btcjson.GetBlockChainInfoResult) []byte {
	status := "finished"
	if bc.Headers != bc.Blocks {
		status = "syncing"
	}
	out, _ := json.Marshal(SyncResponse{
		Status:           status,
		BlockChainHeight: int(bc.Blocks),
		SyncPercentage:   int((bc.Blocks / bc.Headers) * 100),
		Height:           int(bc.Headers),
		Error:            nil,
		Type:             "addrindex-server",
	})
	return out
}
