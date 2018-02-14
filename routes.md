#### `GET /addr/{addr}/utxo`
#### `GET /addr/{addr}/balance`
#### `GET /addr/{addr}/totalReceived`
#### `GET /addr/{addr}/totalSent`
#### `GET /tx/{txid}`
#### `GET /rawtx/{txid}`
#### `POST /messages/verify`

```json
{
  "bitcoinaddress": "string",
  "signature": "string",
  "message": "string",
}
```

#### `POST /tx/send`

```json
{
  "tx": "rawtxstring"
}
```

#### `GET /block/{blockHash}`
#### `GET /block-index/{height}`
#### `GET /status`

```
GET /status?q=getInfo
GET /status?q=getDifficulty
GET /status?q=getBestBlockHash
```

#### `GET /sync`
#### `GET /txs`

```
GET /txs?block=<blockhash>&page=<page>
GET /txs?address=<addr>&page=<page>
```

#### `GET /version`