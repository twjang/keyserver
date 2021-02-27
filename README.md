# Keyserver

This is a basic `keyserver` for terra applications. It contains the following routes:

```
GET     /version
POST    /keys/list
POST    /keys/create
POST    /keys/get/{name}
POST    /keys/delete/{name}
POST    /tx/sign
```

First, build and start the server:

```bash
> make install
> keyserver config
> keyserver serve
```

Then you can use the included CLI to create keys, use the mnemonics to create them in `gaiad` as well:

```bash
# Create a new key with generated mnemonic
# here, "keyringpass" is the keyring password
# this version of keyserver only supports "file" backend
> keyserver keys create keyringpass alice | jq

# Create another key
> keyserver keys create keyringpass bob | jq

# Save the mnemonic from the above command and add it to terracli
> gaiad keys add cynthia --recover --keyring-backend file

# Next create a single node testnet
> gaiad init testing --chain-id testing 
> gaiad add-genesis-account alice 10000000000stake --keyring-backend file
> gaiad add-genesis-account $(keyserver keys show keyringpass bob | jq -r .address) 100000000stake --keyring-backend file
> gaiad gentx alice --keyring-backend file 1000000stake --chain-id testing
> gaiad collect-gentxs
> gaiad start
```

In another window, generate the transaction to sign, sign it and broadcast:
```bash
> gaiad tx bank send $(keyserver keys show keyringpass alice | jq -r .address) $(keyserver keys show keyringpass bob | jq -r .address) 10stake --chain-id testing --memo memo --fees 1stake  --generate-only  > unsigned.json
> keyserver tx sign alice keyringpass testing 1 1 ./unsigned.json > signed.json
> gaiad tx broadcast ./signed.json
{"height":"83","txhash":"2D7135C451931D3B77F2466D617B14AE38DF0153E2D73020F78BD64722523A0E","codespace":"","code":0,"data":"0A060A0473656E64","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"cosmos1zfq2sn550kn347ps27pajktmjdjd4s9e4uut3f\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos15d8yph35s5jddzj967nu8pqcywh9vyyw63p6mg\"},{\"key\":\"sender\",\"value\":\"cosmos1zfq2sn550kn347ps27pajktmjdjd4s9e4uut3f\"},{\"key\":\"amount\",\"value\":\"10stake\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"send"},{"key":"sender","value":"cosmos1zfq2sn550kn347ps27pajktmjdjd4s9e4uut3f"},{"key":"module","value":"bank"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos15d8yph35s5jddzj967nu8pqcywh9vyyw63p6mg"},{"key":"sender","value":"cosmos1zfq2sn550kn347ps27pajktmjdjd4s9e4uut3f"},{"key":"amount","value":"10stake"}]}]}],"info":"","gas_wanted":"200000","gas_used":"61650","tx":null,"timestamp":""}
> gaaid q tx 2D7135C451931D3B77F2466D617B14AE38DF0153E2D73020F78BD64722523A0E
```
