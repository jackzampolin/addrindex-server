package addrindex

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

// ScriptSig models the scriptSig portion of a vin
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// TransactionIns is the response struct for GetRawTransaction
type TransactionIns struct {
	Hex           string    `json:"hex,omitempty"`
	Txid          string    `json:"txid,omitempty"`
	Size          int       `json:"size,omitempty"`
	Version       int       `json:"version,omitempty"`
	Locktime      int       `json:"locktime,omitempty"`
	Vin           []VinIns  `json:"vin,omitempty"`
	Vout          []VoutIns `json:"vout,omitempty"`
	Blockhash     string    `json:"blockhash,omitempty"`
	Height        int       `json:"height,omitempty"`
	Confirmations int       `json:"confirmations,omitempty"`
	Time          int       `json:"time,omitempty"`
	Blocktime     int       `json:"blocktime,omitempty"`
}

// AddrMempoolTransaction represents a transaction in the mempool
// prevtxid and prevout that can be used for marking utxos as spent
// Instead of height there is timestamp that is the time the transaction entered the mempool
type AddrMempoolTransaction struct {
	Address   string `json:"address"`
	Txid      string `json:"txid"`
	Index     int    `json:"index"`
	Satoshis  int    `json:"satoshis"`
	Timestamp int    `json:"timestamp"`
	Prevtxid  string `json:"prevtxid,omitempty"`
	Prevout   int    `json:"prevout,omitempty"`
}

// UTXO takes a mempool transaction and converts it into the output format for /addr/<addr>/utxo
func (amp AddrMempoolTransaction) UTXO() UTXOInsOut {
	return UTXOInsOut{
		Address:       amp.Address,
		Txid:          amp.Txid,
		OutputIndex:   amp.Index,
		Satoshis:      amp.Satoshis,
		Amount:        float64(amp.Satoshis) / 100000000,
		Timestamp:     amp.Timestamp,
		Confirmations: 0,
	}
}

// UTXOIns is an insight representation of a UTXO
type UTXOIns struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	OutputIndex   int     `json:"outputIndex"`
	Script        string  `json:"script,omitempty"`
	Satoshis      int     `json:"satoshis,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	Height        int     `json:"height,omitempty"`
	Timestamp     int     `json:"timestamp,omitempty"`
	Confirmations int     `json:"confirmations"`
}

// UTXOInsOuts is a collection with methods defined for sorting
type UTXOInsOuts []UTXOInsOut

// Len implements sort for UTXOInsOuts
func (s UTXOInsOuts) Len() int {
	return len(s)
}

// Swap implements sort for UTXOInsOuts
func (s UTXOInsOuts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements sort for UTXOInsOuts
func (s UTXOInsOuts) Less(i, j int) bool {
	return s[i].Confirmations < s[j].Confirmations
}

// UTXOInsOut Output representation
type UTXOInsOut struct {
	Address       string  `json:"address"`
	Txid          string  `json:"txid"`
	OutputIndex   int     `json:"vout"`
	Script        string  `json:"script,omitempty"`
	Satoshis      int     `json:"satoshis,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	Height        int     `json:"height,omitempty"`
	Timestamp     int     `json:"ts,omitempty"`
	Confirmations int     `json:"confirmations"`
}

// Enrich adds data to a utxo to make the output format
func (utxo UTXOIns) Enrich(blockHeight int32) UTXOInsOut {
	return UTXOInsOut{
		Address:       utxo.Address,
		Txid:          utxo.Txid,
		OutputIndex:   utxo.OutputIndex,
		Script:        utxo.Script,
		Satoshis:      utxo.Satoshis,
		Amount:        float64(utxo.Satoshis) / 100000000,
		Height:        utxo.Height,
		Confirmations: int(blockHeight) - utxo.Height + 1,
		Timestamp:     utxo.Timestamp,
	}
}

// SpentInfo contains data about spent transaction outputs
type SpentInfo struct {
	Txid   string `json:"txid"`
	Index  int    `json:"index"`
	Height int    `json:"height"`
}

// AddressBalance is the balance of an address
type AddressBalance struct {
	Balance  int `json:"balance"`
	Received int `json:"received"`
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
