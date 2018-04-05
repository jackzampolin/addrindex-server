# `addrindex-server`

This is a server to pair with an individual `bitcoind` with the [extra `bitcore` methods](https://bitcore.io/guides/bitcoin/) enabled.

### API Methods

This server aims to replicate the interface provided by the [Insight API](https://github.com/bitpay/insight-api). The following routes are available:

```
/addr/{addr}/utxo
/addr/{addr}/balance
/addr/{addr}/totalReceived
/addr/{addr}/totalSent
/addr/{addr}/unconfirmedBalance
/tx/{txid}
/txs
/rawtx/{txid}
/tx/send
/messages/verify
/block/{blockHash}
/blocks
/block-index/{height}
/status
/sync
/version
/currency
```

### Configuration

Configuration file lives at `$HOME/.addrindex-server.yaml`. You can also pass one by running `./addrindex-server --config=/path/to/config.yaml`

```yaml
# This is connection information to the addrindex bitcoin node you are running.
# This service is tested to work with the version (v0.14-bitcore) specified in the docker-compose file
host: localhost:8332
usr: user
pass: password
ssl: false
port: 18332
```

### Build

To build the project have a working gopath and run `make`
To build the docker image have docker installed and run `make docker`

### Running

To run this configuration just:

```
$ cd deploy
$ docker-compose up -d
```

That docker-compose file will spin up the address index node. At time of writing that node requires ~1-2 days to sync and ~400 GB of disk. You may need to pass the `-reindex` argument the first time you run `docker-compose up -d`. If you have any issues with this configuration please open an issue.

### Notes:

- `GET /addr/<addr>/utxo` returns utxo from addresses that have been included in the last block. To prove, get the block at chain tip, pull a transaction from the list and grab the address from one of the UTXOs there. Then fetch the UTXOs for that address and make sure it's there. The following is a `bitcoin-cli` implementation of that:

```
bitcoin-cli getaddressutxos "{\"addresses\": [\"$(bitcoin-cli getrawtransaction $(bitcoin-cli getblock $(bitcoin-cli getchaintips | jq -r '.[].hash') | jq -r '.tx[0]') 1 | jq -r '.vout[0].scriptPubKey.addresses[0]')\"]}"
```
