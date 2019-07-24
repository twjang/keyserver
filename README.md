# Keyserver

This is a basic `keyserver` for terra applications. It contains the following routes:

```
GET     /version
GET     /keys
POST    /keys
GET     /keys/{name}?bech=acc
PUT     /keys/{name}
DELETE  /keys/{name}
POST    /tx/sign
POST    /tx/bank/send
POST    /tx/broadcast
```

First, build and start the server:

```bash
> make install
> keyserver config
> keyserver serve
```

Then you can use the included CLI to create keys, use the mnemonics to create them in `terracli` as well:

```bash
# Create a new key with generated mnemonic
> keyserver keys post yun foobarbaz | jq

# Create another key
> keyserver keys post jim foobarbaz | jq

# Save the mnemonic from the above command and add it to terracli
> terracli keys add yun --recover

# Next create a single node testnet
> terrad init testing --chain-id testing
> terracli config chain-id testing
> terrad add-genesis-account yun 10000000000stake
> terrad add-genesis-account $(keyserver keys show jim | jq -r .address) 100000000stake
> terrad gentx --name yun
> terrad collect-gentxs
> terrad start
```

In another window, generate the transaction to sign, sign it and broadcast:
```bash
> mkdir -p test_data
> keyserver tx bank send $(keyserver keys show yun | jq -r .address) $(keyserver keys show jim | jq -r .address) 10000stake testing "memo" 10stake > test_data/unsigned.json
> keyserver tx sign yun foobarbaz testing 0 1 test_data/unsigned.json > test_data/signed.json
> keyserver tx broadcast test_data/signed.json
{"height":"0","txhash":"84CEF8B7FD04DA6FE9C22A6077D8286FA7775CAA0BB06D1D875AE9527A3D15CB"}
> terracli q txs 84CEF8B7FD04DA6FE9C22A6077D8286FA7775CAA0BB06D1D875AE9527A3D15CB
```
