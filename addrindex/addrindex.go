package addrindex

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
)

// const v013 = "104.198.100.79:8332"
// const v014 = "35.185.225.4:8332"
// const v015 = "35.203.175.172:8332"

// func GetPeople(w http.ResponseWriter, r *http.Request) {
//     json.NewEncoder(w).Encode(people)
// }

func txids(txs []*btcjson.SearchRawTransactionsResult) []string {
	out := []string{}
	for _, txn := range txs {
		out = append(out, txn.Txid)
	}
	return out
}

// SearchRawTransactions returns the result of a searchrawtransactions RPC call against the configured bitcoin node
func (as *AddrServer) SearchRawTransactions(addr string, offset int, count int) (SearchRawTransactionsResult, error) {
	out := SearchRawTransactionsResult{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newSearchRawTransactionRequest(addr, offset, count)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "tesxt/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(body, &out)
	if err != nil {
		return out, err
	}
	return out, nil
}

type searchRawTransactionsRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func newSearchRawTransactionRequest(addr string, offset int, count int) []byte {
	srtr := searchRawTransactionsRequest{
		JSONRPC: "1.0",
		Method:  "searchrawtransactions",
		Params:  []interface{}{addr, 1, offset, count},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// UTXOs are a set of UTXO
type UTXOs []UTXO

func (utxos UTXOs) JSON() []byte {
	utxo, err := json.Marshal(utxos)
	if err != nil {
		panic(err)
	}
	return utxo
}

// UTXO models an unspent transaction output
type UTXO struct {
	TransactionHash string   `json:"transaction_hash"`
	Outpoint        Outpoint `json:"outpoint"`
	Value           float64  `json:"value"`
	OutScript       string   `json:"out_script"`
	Confirmations   int      `json:"confirmations"`
}

// Outpoint models the outpoint of a transaction
type Outpoint struct {
	Hash  string `json:"txid"`
	Index int    `json:"vout"`
}

// NewUTXO takes the data from a transaction and creates a UTXO
func newUTXO(txid string, conf int, vout Vout) UTXO {
	return UTXO{
		TransactionHash: txid,
		Outpoint: Outpoint{
			Hash:  txid,
			Index: vout.N,
		},
		Value:         vout.Value,
		OutScript:     vout.ScriptPubKey.Hex,
		Confirmations: conf,
	}
}

// UTXO returns the transaction outputs for a transaction for an address in a format for return
func (txns Transactions) UTXO(addr string) UTXOs {
	outputs := []UTXO{}
	out := []UTXO{}

	// First gather all the outputs from all the transactions that apply to the address
	for _, tx := range txns {
		for _, vout := range tx.Vout {
			if vout.contains(addr) {
				outputs = append(outputs, newUTXO(tx.Txid, tx.Confirmations, vout))
			}
		}
	}

	// Next, filter out spent outputs
	for _, txo := range outputs {
		unspent := true
		for _, tx := range txns {
			for _, vin := range tx.Vin {
				if vin.Txid == txo.TransactionHash {
					unspent = false
				}
			}
		}
		if unspent == true {
			out = append(out, txo)
		}
	}
	return out
}

func (vout Vout) contains(addr string) bool {
	for _, ad := range vout.ScriptPubKey.Addresses {
		if ad == addr {
			return true
		}
	}
	return false
}

// SearchRawTransactionsResult models the raw result from a SearchRawTransactions call
type SearchRawTransactionsResult struct {
	Result Transactions `json:"result"`
	Error  interface{}  `json:"error"`
	ID     interface{}  `json:"id"`
}

type Transactions []Transaction

// Transaction models a bitcoin tansaction
type Transaction struct {
	Txid          string `json:"txid"`
	Hash          string `json:"hash"`
	Size          int    `json:"size"`
	Vsize         int    `json:"vsize"`
	Version       int    `json:"version"`
	Locktime      int    `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	Blockhash     string `json:"blockhash"`
	Confirmations int    `json:"confirmations"`
	Time          int    `json:"time"`
	Blocktime     int    `json:"blocktime"`
	Hex           string `json:"hex"`
}

// Vin models a vin in a transaction
type Vin struct {
	Txid      string    `json:"txid"`
	Vout      int       `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  int64     `json:"sequence"`
}

// ScriptSig models the scriptSig portion of a vin
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// Vout models a vout in a transaction
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// ScriptPubKey models the scriptPubKey portion of a vout
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}
