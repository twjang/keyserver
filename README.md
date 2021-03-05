# Keyserver

This is a basic `keyserver` for terra applications. It contains the following routes:

```
GET     /version
POST    /keys/list
POST    /keys/create
POST    /keys/get/{name}
POST    /keys/delete/{name}
POST    /tx/sign
POST    /tx/encode
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
> export ACC=$(gaiad  q account  $(keyserver keys show keyringpass alice | jq -r .address) -ojson | jq .account_number | sed 's/"//g')
> export SEQ=$(gaiad  q account  $(keyserver keys show keyringpass alice | jq -r .address) -ojson | jq .sequence | sed 's/"//g')
> keyserver tx sign alice keyringpass testing $ACC $SEQ ./unsigned.json > signed.json
> gaiad tx broadcast ./signed.json --broadcast-mode async
{"height":"0","txhash":"F990053766C5984633EDD827762C063AF64225BB4846215B1313795FA8566371","codespace":"","code":0,"data":"","raw_log":"","logs":[],"info":"","gas_wanted":"0","gas_used":"0","tx":null,"timestamp":""}
> gaiad q tx F990053766C5984633EDD827762C063AF64225BB4846215B1313795FA8566371
```
