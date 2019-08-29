package api

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/terra-project/core/app"

	"github.com/terra-project/core/x/treasury"
)

const (
	maxValidAccountValue = int(0x80000000 - 1)
	maxValidIndexalue    = int(0x80000000 - 1)
)

var cdc *codec.Codec

const (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "terra"
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = "terrapub"
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = "terravaloper"
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = "terravaloperpub"
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = "terravalcons"
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = "terravalconspub"
)

func init() {
	cdc = app.MakeCodec()
	config := sdk.GetConfig()
	config.SetCoinType(330)
	config.SetFullFundraiserPath("44'/330'/0'/0/0")
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}

// Server represents the API server
type Server struct {
	Port   int    `json:"port"`
	KeyDir string `json:"key_dir"`
	Node   string `json:"node"`

	Version string `yaml:"version,omitempty"`
	Commit  string `yaml:"commit,omitempty"`
	Branch  string `yaml:"branch,omitempty"`
}

// Router returns the router
func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/version", s.VersionHandler).Methods("GET")
	router.HandleFunc("/keys", s.GetKeys).Methods("GET")
	router.HandleFunc("/keys", s.PostKeys).Methods("POST")
	router.HandleFunc("/keys/{name}", s.GetKey).Methods("GET")
	router.HandleFunc("/keys/{name}", s.PutKey).Methods("PUT")
	router.HandleFunc("/keys/{name}", s.DeleteKey).Methods("DELETE")
	router.HandleFunc("/tx/sign", s.Sign).Methods("POST")
	router.HandleFunc("/tx/broadcast", s.Broadcast).Methods("POST")
	router.HandleFunc("/tx/bank/send", s.BankSend).Methods("POST")
	router.HandleFunc("/tx/encode", s.EncodeTx).Methods("POST")

	return router
}

// SimulateGas simulates gas for a transaction
func (s *Server) SimulateGas(txbytes []byte) (res uint64, err error) {
	result, err := rpcclient.NewHTTP(s.Node, "/websocket").ABCIQueryWithOptions(
		"/app/simulate",
		cmn.HexBytes(txbytes),
		rpcclient.ABCIQueryOptions{},
	)

	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		return 0, errors.New(result.Response.Log)
	}

	var simulationResult sdk.Result
	if err := cdc.UnmarshalBinaryLengthPrefixed(result.Response.Value, &simulationResult); err != nil {
		return 0, err
	}

	return simulationResult.GasUsed, nil
}

// LoadCurrentEpoch load current epoch
func (s *Server) LoadCurrentEpoch() (res int64, err error) {
	result, err := rpcclient.NewHTTP(s.Node, "/websocket").ABCIQueryWithOptions(
		"custom/treasury/currentEpoch",
		[]byte{},
		rpcclient.ABCIQueryOptions{},
	)
	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		return 0, errors.New(result.Response.Log)
	}

	var currentEpoch int64
	if err := cdc.UnmarshalJSON(result.Response.Value, &currentEpoch); err != nil {
		return 0, err
	}

	return currentEpoch, nil
}

// LoadTaxRate load tax-rate
func (s *Server) LoadTaxRate(epoch int64) (res sdk.Dec, err error) {
	params := treasury.NewQueryTaxRateParams(epoch)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	result, err := rpcclient.NewHTTP(s.Node, "/websocket").ABCIQueryWithOptions(
		"custom/treasury/taxRate",
		cmn.HexBytes(bz),
		rpcclient.ABCIQueryOptions{},
	)
	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		return sdk.ZeroDec(), errors.New(result.Response.Log)
	}

	var taxRate sdk.Dec
	if err := cdc.UnmarshalJSON(result.Response.Value, &taxRate); err != nil {
		return sdk.ZeroDec(), err
	}

	return taxRate, nil
}

// LoadTaxCap load tax-cap
func (s *Server) LoadTaxCap(denom string) (res sdk.Int, err error) {
	params := treasury.NewQueryTaxCapParams(denom)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	result, err := rpcclient.NewHTTP(s.Node, "/websocket").ABCIQueryWithOptions(
		"custom/treasury/taxCap",
		cmn.HexBytes(bz),
		rpcclient.ABCIQueryOptions{},
	)
	if err != nil {
		return
	}

	if !result.Response.IsOK() {
		return sdk.ZeroInt(), errors.New(result.Response.Log)
	}

	var taxRate sdk.Int
	if err := cdc.UnmarshalJSON(result.Response.Value, &taxRate); err != nil {
		return sdk.ZeroInt(), err
	}

	return taxRate, nil
}
