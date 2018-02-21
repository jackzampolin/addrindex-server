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

// SpentInfo contains data about spent transaction outputs
type SpentInfo struct {
	Txid   string `json:"txid"`
	Index  int    `json:"index"`
	Height int    `json:"height"`
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

// UTXOIns is an insight representation of a UTXO
type UTXOIns struct {
	Address     string `json:"address"`
	Txid        string `json:"txid"`
	OutputIndex int    `json:"outputIndex"`
	Script      string `json:"script"`
	Satoshis    int    `json:"satoshis"`
	Height      int    `json:"height"`
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
