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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SearchRawTransactions returns the result of a searchrawtransactions RPC call against the configured bitcoin node
func (as *AddrServer) SearchRawTransactions(addr string, offset int, count int) (SearchRawTransactionsResult, error) {
	out := SearchRawTransactionsResult{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newSearchRawTransactionRequest(addr, offset, count)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

type BitcoreRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func newSearchRawTransactionRequest(addr string, offset int, count int) []byte {
	srtr := BitcoreRequest{
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

// JSON returns the UTXO in a format for return to client
func (utxos UTXOs) JSON() []byte {
	utxo, err := json.Marshal(utxos)
	if err != nil {
		panic(err)
	}
	return utxo
}

// Balance returns the balance for a address
func (utxos UTXOs) Balance() []byte {
	var val int64
	// var byt = []byte{}
	for _, utxo := range utxos {
		val += (utxo.Value)
	}
	return []byte(fmt.Sprintf("%v", val))
}

// UTXO models an unspent transaction output
type UTXO struct {
	TransactionHash string   `json:"transaction_hash"`
	Outpoint        Outpoint `json:"outpoint"`
	Value           int64    `json:"value"`
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
		Value:         int64(vout.Value * 100000000),
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

// Transactions is a group of Transaction
type Transactions []Transaction

// Received calculates the recieved btc by the address
func (txns Transactions) Received(addr string) []byte {
	var val int64
	for _, txn := range txns {
		for _, vout := range txn.Vout {
			if vout.contains(addr) {
				val += int64(vout.Value * 100000000)
			}
		}
	}
	return []byte(fmt.Sprintf("%v", val))
}

// Sent calculates the sent btc by the address
func (txns Transactions) Sent(addr string) []byte {
	var val int64
	outputs := []UTXO{}

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
		if unspent == false {
			val += txo.Value
		}
	}
	return []byte(fmt.Sprintf("%v", val))
}

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

// JSON returns the Transaction in a format for return to client
func (tx Transaction) JSON() []byte {
	txn, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return txn
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

func newBitcoreAddressesStartEndRequest(addresses []string, start int, end int) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "searchrawtransactions",
		Params:  []interface{}{addresses, start, end},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// GetAddressTxIDs searches for all txid associated with an address.
//   - Most recient last
//   - Only confirmed
func (as *AddrServer) GetAddressTxIDs(addresses []string, start, end int) ([]string, error) {
	out := []string{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreAddressesStartEndRequest(addresses, start, end)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetAddressDeltas searches for all inputs, outputs and top level detail for transactions
//   - Only confirmed
//   - Negative "satoshis" = Vin
//   - Positve "satoshis"  = Vout
// {
//   "satoshis": 30000,
//   "txid": "20fb69a94413637cb50f65e473f91d2599a04d5a0bf9bf6a5e9e843df2710ea4",
//   "index": 0,
//   "blockindex": 165,
//   "height": 228208,
//   "address": "12cbQLTFMXRnSzktFkuoG3eHoMeFtpTu3S"
// }
func (as *AddrServer) GetAddressDeltas(addresses []string, start, end int) (GetAddressDeltasResponse, error) {
	out := GetAddressDeltasResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreAddressesStartEndRequest(addresses, start, end)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetAddressDeltasResponse is the response struct for GetAddressDeltas
type GetAddressDeltasResponse []struct {
	Satoshis   int    `json:"satoshis"`
	Txid       string `json:"txid"`
	Index      int    `json:"index"`
	Blockindex int    `json:"blockindex"`
	Height     int    `json:"height"`
	Address    string `json:"address"`
}

func newBitcoreAddessesRequest(addresses []string) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "searchrawtransactions",
		Params:  []interface{}{addresses},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// GetAddressBalance returns the balance of confirmed transactions
func (as *AddrServer) GetAddressBalance(addresses []string) (GetAddressBalanceResult, error) {
	out := GetAddressBalanceResult{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreAddessesRequest(addresses)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetAddressBalanceResult is the response struct for GetAddressBalance
type GetAddressBalanceResult struct {
	Balance  int `json:"balance"`
	Received int `json:"received"`
}

// GetAddressUTXOs returns the list of UTXO for an address sorted by block height
func (as *AddrServer) GetAddressUTXOs(addresses []string) (GetAddressUTXOsResponse, error) {
	out := GetAddressUTXOsResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreAddessesRequest(addresses)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetAddressUTXOsResponse is the response struct for GetAddressUTXOs
type GetAddressUTXOsResponse []struct {
	Address     string `json:"address"`
	Txid        string `json:"txid"`
	OutputIndex int    `json:"outputIndex"`
	Script      string `json:"script"`
	Satoshis    int    `json:"satoshis"`
	Height      int    `json:"height"`
}

// GetAddressMempool returns GetAddressDeltas but for the mempool:
//   - Only Mempool
//   - Negative "satoshis" = Vin
//   - Positve "satoshis"  = Vout
// {
// 	"address": "3M366gYcKHbvun6YYF1Xim6sDeT5JSTVUy",
// 	"txid": "ff21363aa331f2dc7bbf70acc7eefb7a4080645d30b4e319ca190ceaecbcce42",
// 	"index": 0,
// 	"satoshis": -10684303,
// 	"timestamp": 1463602662,
// 	"prevtxid": "0c15f067d6b082f4dcc2740f039d33bb4f47b23c79ceae880ca759268389f82a",
// 	"prevout": 1
// },
func (as *AddrServer) GetAddressMempool(addresses []string) (GetAddressMempoolResponse, error) {
	out := GetAddressMempoolResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreAddessesRequest(addresses)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetAddressMempoolResponse is the response struct for GetAddressMempool
type GetAddressMempoolResponse []struct {
	Address   string `json:"address"`
	Txid      string `json:"txid"`
	Index     int    `json:"index"`
	Satoshis  int    `json:"satoshis"`
	Timestamp int    `json:"timestamp"`
	Prevtxid  string `json:"prevtxid,omitempty"`
	Prevout   int    `json:"prevout,omitempty"`
}

func newBitcoreStartEndRequest(start, end int) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "searchrawtransactions",
		Params:  []interface{}{start, end},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// GetBlockHashes returns blockhashes between two unix epoch timestamps
func (as *AddrServer) GetBlockHashes(start, end int) ([]string, error) {
	out := []string{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreStartEndRequest(start, end)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

func newBitcoreTxWithIndexRequest(txid string, index int) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "searchrawtransactions",
		Params:  []interface{}{txid, index},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// GetSpentInfo returns the txid and input index that has spent the output
func (as *AddrServer) GetSpentInfo(txid string, index int) (GetSpentInfoResponse, error) {
	out := GetSpentInfoResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreTxWithIndexRequest(txid, index)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetSpentInfoResponse is response struct for GetSpentInfo
type GetSpentInfoResponse struct {
	Txid   string `json:"txid"`
	Index  int    `json:"index"`
	Height int    `json:"height"`
}

func newBitcoreRawTxRequest(addresses string) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "searchrawtransactions",
		Params:  []interface{}{addresses, 1},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// GetRawTransaction verbose result will now has some additional fields added when spentindex is enabled.
// The vin values will include value (a float in BTC) and valueSat (an integer in satoshis) with the
// previous output value as well as the  address. The vout values will also now include a valueSat
// (an integer in satoshis). It will also include  spentTxId, spentIndex and spentHeight that corresponds
// with the input that spent the output.
func (as *AddrServer) GetRawTransaction(addr string) (GetRawTransactionResponse, error) {
	out := GetRawTransactionResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(newBitcoreRawTxRequest(addr)))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "text/plain")

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

// GetRawTransactionResponse is the response struct for GetRawTransaction
type GetRawTransactionResponse struct {
	Hex           string    `json:"hex"`
	Txid          string    `json:"txid"`
	Size          int       `json:"size"`
	Version       int       `json:"version"`
	Locktime      int       `json:"locktime"`
	Vin           []VinIns  `json:"vin"`
	Vout          []VoutIns `json:"vout"`
	Blockhash     string    `json:"blockhash"`
	Height        int       `json:"height"`
	Confirmations int       `json:"confirmations"`
	Time          int       `json:"time"`
	Blocktime     int       `json:"blocktime"`
}

// VinIns is the bitcore representation of a Vin
type VinIns struct {
	Txid      string    `json:"txid"`
	Vout      int       `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Value     float64   `json:"value"`
	ValueSat  int       `json:"valueSat"`
	Address   string    `json:"address"`
	Sequence  int64     `json:"sequence"`
}

// ScriptPubKeyIns is the bitcore representation of a ScriptPubKey
type ScriptPubKeyIns struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

// VoutIns is the bitcore representation of a Vout
type VoutIns struct {
	Value        float64         `json:"value"`
	ValueSat     int             `json:"valueSat"`
	N            int             `json:"n"`
	ScriptPubKey ScriptPubKeyIns `json:"scriptPubKey"`
	SpentTxID    string          `json:"spentTxId,omitempty"`
	SpentIndex   int             `json:"spentIndex,omitempty"`
	SpentHeight  int             `json:"spentHeight,omitempty"`
}
