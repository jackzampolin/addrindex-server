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
