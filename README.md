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
```