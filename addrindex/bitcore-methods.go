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
	"io/ioutil"
	"net/http"
)

// BitcoreRequest represents a request to a bitcore node
type BitcoreRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func getAddressTxIDsRequest(addresses []string, start int, end int) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getaddresstxids",
		Params: []interface{}{map[string]interface{}{
			"addresses": addresses,
			"start":     start,
			"end":       end,
		}},
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
func (as *AddrServer) GetAddressTxIDs(addresses []string, start, end int) (GetAddressTxIDsResponse, error) {
	out := GetAddressTxIDsResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getAddressTxIDsRequest(addresses, start, end)))
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

// GetAddressTxIDsResponse wraps the return
type GetAddressTxIDsResponse struct {
	Result []string    `json:"result"`
	Error  interface{} `json:"error"`
	ID     interface{} `json:"id"`
}

func getAddressDeltasRequest(addresses []string, start int, end int) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getaddressdeltas",
		Params:  []interface{}{map[string]interface{}{"addresses": addresses, "start": start, "end": end}},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
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
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getAddressDeltasRequest(addresses, start, end)))
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
type GetAddressDeltasResponse struct {
	Result []AddressDelta `json:"result"`
	Error  interface{}    `json:"error"`
	ID     interface{}    `json:"id"`
}

// AddressDelta represents a balance change for an address
type AddressDelta struct {
	Satoshis   int    `json:"satoshis"`
	Txid       string `json:"txid"`
	Index      int    `json:"index"`
	Blockindex int    `json:"blockindex"`
	Height     int    `json:"height"`
	Address    string `json:"address"`
}

func getAddressBalanceRequest(addresses []string) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getaddressbalance",
		Params:  []interface{}{map[string][]string{"addresses": addresses}},
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
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getAddressBalanceRequest(addresses)))
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
	Result AddressBalance `json:"result"`
	Error  interface{}    `json:"error"`
	ID     interface{}    `json:"id"`
}

// AddressBalance is the balance of an address
type AddressBalance struct {
	Balance  int `json:"balance"`
	Received int `json:"received"`
}

func getAddressUTXORequest(addresses []string) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getaddressutxos",
		Params:  []interface{}{map[string][]string{"addresses": addresses}},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
}

// GetAddressUTXOs returns the list of UTXO for an address sorted by block height
func (as *AddrServer) GetAddressUTXOs(addresses []string) (GetAddressUTXOsResponse, error) {
	out := GetAddressUTXOsResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getAddressUTXORequest(addresses)))
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
type GetAddressUTXOsResponse struct {
	Result []UTXOIns   `json:"result"`
	Error  interface{} `json:"error"`
	ID     interface{} `json:"id"`
}

// UTXOIns is an insight representation of a UTXO
type UTXOIns struct {
	Address     string `json:"address"`
	Txid        string `json:"txid"`
	OutputIndex int    `json:"outputIndex"`
	Script      string `json:"script"`
	Satoshis    int    `json:"satoshis"`
	Height      int    `json:"height"`
}

func getAddressMempoolRequest(addresses []string) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getaddressmempool",
		Params:  []interface{}{map[string][]string{"addresses": addresses}},
	}
	out, err := json.Marshal(srtr)
	if err != nil {
		panic(err)
	}
	return out
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
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getAddressMempoolRequest(addresses)))
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
type GetAddressMempoolResponse struct {
	Result []AddrMempoolTransaction `json:"result"`
	Error  interface{}              `json:"error"`
	ID     interface{}              `json:"id"`
}

// AddrMempoolTransaction represents a transaction in the mempool
// prevtxid and prevout that can be used for marking utxos as spent
// Instead of height there is timestamp that is the time the transaction entered the mempool
type AddrMempoolTransaction []struct {
	Address   string `json:"address"`
	Txid      string `json:"txid"`
	Index     int    `json:"index"`
	Satoshis  int    `json:"satoshis"`
	Timestamp int    `json:"timestamp"`
	Prevtxid  string `json:"prevtxid,omitempty"`
	Prevout   int    `json:"prevout,omitempty"`
}

func getBlockHashesRequest(start, end int) []byte {
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

// GetBlockHashes returns blockhashes between two unix epoch timestamps (seconds)
func (as *AddrServer) GetBlockHashes(start, end int) (GetBlockHashesResponse, error) {
	out := GetBlockHashesResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getBlockHashesRequest(start, end)))
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

// GetBlockHashesResponse is response struct for GetBlockHashes
type GetBlockHashesResponse struct {
	Result []string    `json:"result"`
	Error  interface{} `json:"error"`
	ID     interface{} `json:"id"`
}

func getSpentInfoRequest(txid string, index int) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getspentinfo",
		Params:  []interface{}{map[string]interface{}{"tx": txid, "index": index}},
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
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getSpentInfoRequest(txid, index)))
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
	Result SpentInfo   `json:"result"`
	Error  interface{} `json:"error"`
	ID     interface{} `json:"id"`
}

// SpentInfo contains data about spent transaction outputs
type SpentInfo struct {
	Txid   string `json:"txid"`
	Index  int    `json:"index"`
	Height int    `json:"height"`
}

// getRawTransactionRequest formats the json payload for the getrawtransaction RPC
func getRawTransactionRequest(addresses string) []byte {
	srtr := BitcoreRequest{
		JSONRPC: "1.0",
		Method:  "getrawtransaction",
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
func (as *AddrServer) GetRawTransaction(txn string) (GetRawTransactionResponse, error) {
	out := GetRawTransactionResponse{}
	req, err := http.NewRequest("POST", as.URL(), bytes.NewBuffer(getRawTransactionRequest(txn)))
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

// GetRawTransactionResponse facilitates return
type GetRawTransactionResponse struct {
	Result TransactionIns `json:"result"`
	Error  interface{}    `json:"error"`
	ID     interface{}    `json:"id"`
}

// TransactionIns is the response struct for GetRawTransaction
type TransactionIns struct {
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
