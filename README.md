# `addrindex-server`

This is a server to pair with an individual `bitcoind` with the [extra `bitcore` methods](https://bitcore.io/guides/bitcoin/) enabled. 

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
To build the docker image have docker installed and run `make docker`

### Notes:

- `GET /addr/<addr>/utxo` returns utxo from addresses that have been included in the last block. To prove, get the block at chain tip, pull a transaction from the list and grab the address from one of the UTXOs there. Then fetch the UTXOs for that address and make sure it's there. The following is a `bitcoin-cli` implementation of that:

```
bitcoin-cli getaddressutxos "{\"addresses\": [\"$(bitcoin-cli getrawtransaction $(bitcoin-cli getblock $(bitcoin-cli getchaintips | jq -r '.[].hash') | jq -r '.tx[0]') 1 | jq -r '.vout[0].scriptPubKey.addresses[0]')\"]}"
```
