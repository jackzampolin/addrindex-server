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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/gorilla/mux"
)

// BlockstackStartBlock represents the point on the bitcoin blockchain where blockstack started
const BlockstackStartBlock = 373601

// HandleAddrUTXO handles the /addr/<addr>/utxo route
func (as *AddrServer) HandleAddrUTXO(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// Fetch current block info
	info, err := as.Client.GetInfo()
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("failed to getInfo", err))
		return
	}

	// paginate through transactions
	txns, err := as.GetAddressUTXOs([]string{addr})
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}

	mptxns, err := as.GetAddressMempool([]string{addr})
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching mempool transactions for address", err))
		return
	}

	// If there are no mempool transactions then just return the historical
	if len(mptxns.Result) < 1 {
		o := UTXOInsOuts{}
		for _, utxo := range txns.Result {
			o = append(o, utxo.Enrich(info.Blocks))
		}
		sort.Sort(o)
		out, _ := json.Marshal(o)
		w.Write(out)
		return
	}

	var check UTXOInsOuts
	var out UTXOInsOuts

	for _, tx := range txns.Result {
		check = append(check, tx.Enrich(info.Blocks))
	}

	for _, mptx := range mptxns.Result {
		if mptx.Prevtxid == "" {
			check = append(check, mptx.UTXO())
		}
	}

	for _, toCheck := range check {
		valid := true
		for _, mptx := range mptxns.Result {
			if mptx.Prevtxid == toCheck.Txid && toCheck.OutputIndex == mptx.Prevout {
				valid = false
			}
		}
		if valid {
			out = append(out, toCheck)
		}
	}

	// Sort by confirmations and return
	sort.Sort(out)
	o, _ := json.Marshal(out)
	w.Write(o)
}

// HandleAddrUnconfirmedBalance handles the /addr/<addr>/unconfirmedBalance route
func (as *AddrServer) HandleAddrUnconfirmedBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	mptxns, err := as.GetAddressMempool([]string{addr})
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching mempool transactions for address", err))
		return
	}

	unconfirmed := 0

	for _, mptx := range mptxns.Result {
		unconfirmed += mptx.Satoshis
	}

	out, _ := json.Marshal(unconfirmed)
	w.Write(out)
}

// HandleGetBlocks handles the /blocks route
func (as *AddrServer) HandleGetBlocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	limit := "10"
	if len(query["limit"]) > 0 {
		limit = query["limit"][0]
	}
	lim, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("failed parsing ?limit={val}", err))
		return
	}
	w.Write(as.GetBlocksResponse(lim))
}

// HandleAddrBalance handles the /addr/<addr>/balance route
func (as *AddrServer) HandleAddrBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.GetAddressBalance([]string{addr})
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}
	out, _ := json.Marshal(txns.Result.Balance)
	w.Write(out)
}

// HandleAddrRecieved handles the /addr/<addr>/totalReceived route
func (as *AddrServer) HandleAddrRecieved(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.GetAddressBalance([]string{addr})
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}
	out, _ := json.Marshal(txns.Result.Received)
	w.Write(out)
}

// HandleAddrSent handles the /addr/<addr>/totalSent route
func (as *AddrServer) HandleAddrSent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["addr"]

	// paginate through transactions
	txns, err := as.GetAddressBalance([]string{addr})
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}
	out, _ := json.Marshal(txns.Result.Received - txns.Result.Balance)
	w.Write(out)
}

// HandleTxGet handles the /tx/<txid> route
func (as *AddrServer) HandleTxGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	txid := mux.Vars(r)["txid"]

	// paginate through transactions
	txns, err := as.GetRawTransaction(txid)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}
	out, _ := json.Marshal(txns.Result)
	w.Write(out)
}

// HandleRawTxGet handles the /rawtx/<txid> route
func (as *AddrServer) HandleRawTxGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := mux.Vars(r)["txid"]

	// paginate through transactions
	txns, err := as.GetRawTransaction(addr)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching all transactions for address", err))
		return
	}
	out, _ := json.Marshal(map[string]string{"rawtx": txns.Result.Hex})
	w.Write(out)
}

// HandleTransactionSend handles the /tx/send route
// TODO: Write a test for this
func (as *AddrServer) HandleTransactionSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tx TxPost

	// Read post body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to read post body", err))
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, &tx)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to unmarshall body", err))
		return
	}

	// Convert hex to string
	dec, err := hex.DecodeString(tx.Tx)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to decode hex string", err))
		return
	}

	// Convert tansaction to send format
	txn, err := btcutil.NewTxFromBytes(dec)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to parse transaction", err))
		return
	}

	ret, err := as.Client.SendRawTransaction(txn.MsgTx(), true)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to post transaction to node", err))
		return
	}

	out, _ := json.Marshal(ret)
	w.Write(out)
}

// HandleMessagesVerify handles the /messages/verify route
// TODO: Write a test for this
func (as *AddrServer) HandleMessagesVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tx VerifyPost

	// Read post body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to read post body", err))
		return
	}

	// Unmarshal
	err = json.Unmarshal(b, tx)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to unmarshall body", err))
		return
	}

	addr, err := btcutil.DecodeAddress(tx.BitcoinAddress, nil)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("unable to decode bitcoin address", err))
		return
	}

	ret, err := as.Client.VerifyMessage(addr, tx.Signature, tx.Message)
	if err != nil {
		w.WriteHeader(400)
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
		w.WriteHeader(400)
		w.Write(NewPostError("error parsing txhash", err))
		return
	}

	// paginate through transactions
	block, err := as.Client.GetBlockVerbose(hash)
	if err != nil {
		w.WriteHeader(400)
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
		w.WriteHeader(400)
		w.Write(NewPostError("error parsing blockheight", err))
		return
	}

	block, err := as.Client.GetBlockHash(h)
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching blockhash", err))
		return
	}

	bh, _ := json.Marshal(BlockHashReturn{BlockHash: block.String()})
	w.Write(bh)
}

// HandleGetSync handles the /sync route
func (as *AddrServer) HandleGetSync(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chainInfo, err := as.Client.GetBlockChainInfo()
	if err != nil {
		w.WriteHeader(400)
		w.Write(NewPostError("error fetching blockchain info", err))
		return
	}

	w.Write(NewSyncResponse(chainInfo))
}

// HandleGetStatus handles the /status route
func (as *AddrServer) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	method := "getInfo"

	if len(query["q"]) > 0 {
		method = query["q"][0]
	}

	switch method {
	case "getDifficulty":
		info, err := as.Client.GetDifficulty()
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("failed to getDifficulty", err))
			return
		}
		w.Write(NewGetDifficultyReturn(info))
	case "getBestBlockHash":
		info, err := as.Client.GetBestBlockHash()
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("failed to getBestBlockHash", err))
			return
		}
		w.Write(NewGetBestBlockHashReturn(info.String()))
	default:
		info, err := as.Client.GetInfo()
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("failed to getInfo", err))
			return
		}
		out, _ := json.Marshal(info)
		w.Write(out)
	}
}

// HandleGetTransactions handles the /txs route
func (as *AddrServer) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	var (
		page    int
		address string
		block   string
	)

	if len(query["page"]) > 0 {
		pg, err := strconv.ParseInt(query["page"][0], 10, 64)
		if err != nil {
			page = 0
		} else {
			page = int(pg)
		}
	} else {
		page = 0
	}

	if len(query["address"]) > 0 {
		address = query["address"][0]
	}

	if len(query["block"]) > 0 {
		block = query["block"][0]
	} else if len(query["block"]) > 1 {
		w.WriteHeader(400)
		w.Write(NewPostError("only one block accepted in query", fmt.Errorf("")))
		return
	}

	if address != "" {
		// Fetch Block Height
		info, err := as.Client.GetInfo()
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("failed to getInfo", err))
			return
		}

		// paginate through transactions
		txns, err := as.GetAddressTxIDs([]string{address}, BlockstackStartBlock, int(info.Blocks))
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("error fetching page of transactions for address", err))
			return
		}

		var retTxns []string
		var out []TransactionIns

		// Pull off a page of transactions
		if len(txns.Result) < 10 {
			retTxns = txns.Result
		} else if len(txns.Result) > ((page + 1) * 10) {
			retTxns = []string{}
		} else if len(txns.Result) > (page*10) && len(txns.Result) < ((page+1)*10) {
			retTxns = txns.Result[page*10:]
		} else {
			retTxns = txns.Result[page*10 : (page+1)*10]
		}

		for _, txid := range retTxns {
			tx, err := as.GetRawTransaction(txid)
			if err != nil {
				w.WriteHeader(400)
				w.Write(NewPostError("error fetching page of transactions for address", err))
				return
			}
			out = append(out, tx.Result)
		}

		o, _ := json.Marshal(out)
		w.Write(o)
		return
	}

	if block != "" {
		// Make the chainhash for fetching data
		blockhash, err := chainhash.NewHashFromStr(block)
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("error parsing blockhash", err))
			return
		}

		// Fetch block data
		blockData, err := as.Client.GetBlockVerbose(blockhash)
		if err != nil {
			w.WriteHeader(400)
			w.Write(NewPostError("failed to fetch block transactions", err))
			return
		}

		// Initialize output
		var txns = []*btcjson.TxRawResult{}

		// fetch proper slice of transactions
		var txs []string

		// Pick the proper slice from the txs array
		if len(blockData.Tx) < ((page) * 10) {
			// If there is no data left to fetch, return error
			w.WriteHeader(400)
			w.Write(NewPostError("Out of bounds", fmt.Errorf("page %v doesn't exist", page)))
			return
			// If it's the last page, just return the last few transactions
		} else if len(blockData.Tx)-((page+1)*10) <= 0 {
			txs = blockData.Tx[int(page)*10:]
			// Otherwise return a full page
		} else {
			txs = blockData.Tx[int(page)*10 : int(page+1)*10]
		}

		// Fetch individual transaction data and append it to the txns array
		for _, tx := range txs {
			txhash, err := chainhash.NewHashFromStr(tx)
			if err != nil {
				w.WriteHeader(400)
				w.Write(NewPostError(fmt.Sprintf("error parsing transaction %v", tx), err))
				return
			}

			txData, err := as.Client.GetRawTransactionVerbose(txhash)
			if err != nil {
				w.WriteHeader(400)
				w.Write(NewPostError(fmt.Sprintf("error fetching transaction details: %v", tx), err))
				return
			}
			txns = append(txns, txData)
		}

		// Return the JSON
		out, _ := json.Marshal(txns)
		w.Write(out)
		return
	}
	w.WriteHeader(400)
	w.Write(NewPostError("Need to pass ?block=BLOCKHASH or ?address=ADDR", fmt.Errorf("")))
}

// HandleGetVersion handles the /version route
func (as *AddrServer) HandleGetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write(as.version())
}

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
	Error   string `json:"error"`
}

// NewPostError is a convinence function for returning errors to clients
func NewPostError(msg string, err error) []byte {
	out, _ := json.Marshal(PostError{
		Message: msg,
		Error:   err.Error(),
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

// HandleGetCurrency handles the /currency route
func (as *AddrServer) HandleGetCurrency(w http.ResponseWriter, r *http.Request) {
	cd := NewCurrencyData()
	cd.Status = 200
	w.Write(cd.JSON())
}
