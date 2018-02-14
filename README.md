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