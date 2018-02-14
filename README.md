# `addrindex-server`

This is a server to pair with an individual `bitcoind` node with the `addrindex` patch enabled. 

### Configuration

Configuration file lives at `$HOME/.addrindex-server.yaml`. You can also pass one by running `./addrindex-server --config=/path/to/config.yaml`

```yaml
host: localhost:8332
usr: user
pass: password
ssl: false
port: 18332

# The number of transactions to return with each request
# This affects the /addr/ routes and no other ones
transactions: 50
```

### Build

To build the project have a working gopath and run `make`

### Notes:

- `/tx/{txid}` route has some differences with the insight api. Looks like theres some data enrichment going on

> curl -s utxo.technofractal.com:18332/tx/f7cbe7871abc534e5fa287e26973a1ba076783fbad392665211c9b8bfca79d7c | jq

```json
{
  "hex": "0200000001440e64146e72439da521dc20f246a8d4b3bcf5ead4dc6a09523fa498fab10ef10000000088453042021e6fcf15e8d272d1a995af6fcc9d6c0c2f4c0b6b0525142e8af866dd8dad4b022059181242d993aa5101bab12ee86e7df0ce2dc308db930af9278b84d634fdeb45014104af0f6203b43276804c2cbbda0d10c797e61805b56e7abb5dd0cd90f69113dc1bbbaa4bae9c7eee1fde4cab190d54a0da60bca489702a8f7daa895fcaebd2b136ffffffff0134760000000000001976a91428dc4730af1538b45ee1f0e3df7a9c9f31b000a388ac00000000",
  "txid": "f7cbe7871abc534e5fa287e26973a1ba076783fbad392665211c9b8bfca79d7c",
  "hash": "f7cbe7871abc534e5fa287e26973a1ba076783fbad392665211c9b8bfca79d7c",
  "size": 221,
  "vsize": 221,
  "version": 2,
  "locktime": 0,
  "vin": [
    {
      "txid": "f10eb1fa98a43f52096adcd4eaf5bcb3d4a846f220dc21a59d43726e14640e44",
      "vout": 0,
      "scriptSig": {
        "asm": "3042021e6fcf15e8d272d1a995af6fcc9d6c0c2f4c0b6b0525142e8af866dd8dad4b022059181242d993aa5101bab12ee86e7df0ce2dc308db930af9278b84d634fdeb45[ALL] 04af0f6203b43276804c2cbbda0d10c797e61805b56e7abb5dd0cd90f69113dc1bbbaa4bae9c7eee1fde4cab190d54a0da60bca489702a8f7daa895fcaebd2b136",
        "hex": "453042021e6fcf15e8d272d1a995af6fcc9d6c0c2f4c0b6b0525142e8af866dd8dad4b022059181242d993aa5101bab12ee86e7df0ce2dc308db930af9278b84d634fdeb45014104af0f6203b43276804c2cbbda0d10c797e61805b56e7abb5dd0cd90f69113dc1bbbaa4bae9c7eee1fde4cab190d54a0da60bca489702a8f7daa895fcaebd2b136"
      },
      "sequence": 4294967295
    }
  ],
  "vout": [
    {
      "value": 0.0003026,
      "n": 0,
      "scriptPubKey": {
        "asm": "OP_DUP OP_HASH160 28dc4730af1538b45ee1f0e3df7a9c9f31b000a3 OP_EQUALVERIFY OP_CHECKSIG",
        "hex": "76a91428dc4730af1538b45ee1f0e3df7a9c9f31b000a388ac",
        "reqSigs": 1,
        "type": "pubkeyhash",
        "addresses": [
          "14j3us479jTySjG6Tr4uY4Nrx37iTM8QHF"
        ]
      }
    }
  ],
  "blockhash": "00000000000000000063612423ef57aecfa050f4b7bd096024cb09323cea0593",
  "confirmations": 1248,
  "time": 1517865386,
  "blocktime": 1517865386
}
```

> curl -sL explorer.blockstack.org/insight-api/tx/f7cbe7871abc534e5fa287e26973a1ba076783fbad392665211c9b8bfca79d7c | jq

```json
{
  "txid": "f7cbe7871abc534e5fa287e26973a1ba076783fbad392665211c9b8bfca79d7c",
  "version": 2,
  "locktime": 0,
  "vin": [
    {
      "txid": "f10eb1fa98a43f52096adcd4eaf5bcb3d4a846f220dc21a59d43726e14640e44",
      "vout": 0,
      "sequence": 4294967295,
      "n": 0,
      "scriptSig": {
        "hex": "453042021e6fcf15e8d272d1a995af6fcc9d6c0c2f4c0b6b0525142e8af866dd8dad4b022059181242d993aa5101bab12ee86e7df0ce2dc308db930af9278b84d634fdeb45014104af0f6203b43276804c2cbbda0d10c797e61805b56e7abb5dd0cd90f69113dc1bbbaa4bae9c7eee1fde4cab190d54a0da60bca489702a8f7daa895fcaebd2b136",
        "asm": "3042021e6fcf15e8d272d1a995af6fcc9d6c0c2f4c0b6b0525142e8af866dd8dad4b022059181242d993aa5101bab12ee86e7df0ce2dc308db930af9278b84d634fdeb45[ALL] 04af0f6203b43276804c2cbbda0d10c797e61805b56e7abb5dd0cd90f69113dc1bbbaa4bae9c7eee1fde4cab190d54a0da60bca489702a8f7daa895fcaebd2b136"
      },
      "addr": "1FLAMEN6rq2BqMnkUmsJBqCGWdwgVKcegd",
      "valueSat": 35600,
      "value": 0.000356,
      "doubleSpentTxID": null
    }
  ],
  "vout": [
    {
      "value": "0.00030260",
      "n": 0,
      "scriptPubKey": {
        "hex": "76a91428dc4730af1538b45ee1f0e3df7a9c9f31b000a388ac",
        "asm": "OP_DUP OP_HASH160 28dc4730af1538b45ee1f0e3df7a9c9f31b000a3 OP_EQUALVERIFY OP_CHECKSIG",
        "addresses": [
          "14j3us479jTySjG6Tr4uY4Nrx37iTM8QHF"
        ],
        "type": "pubkeyhash"
      },
      "spentTxId": null,
      "spentIndex": null,
      "spentHeight": null
    }
  ],
  "blockhash": "00000000000000000063612423ef57aecfa050f4b7bd096024cb09323cea0593",
  "blockheight": 507850,
  "confirmations": 1248,
  "time": 1517865386,
  "blocktime": 1517865386,
  "valueOut": 0.0003026,
  "size": 221,
  "valueIn": 0.000356,
  "fees": 5.34e-05
}
```