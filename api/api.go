package api

import (
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gaia "github.com/cosmos/gaia/v4/app"
	"github.com/cosmos/gaia/v4/app/params"
	"github.com/gorilla/mux"
)

const (
	maxValidAccountValue = int(0x80000000 - 1)
	maxValidIndexalue    = int(0x80000000 - 1)
)

var encodingConfig params.EncodingConfig
var cdc codec.Marshaler
var legacyCdc *codec.LegacyAmino

func init() {
	encodingConfig = gaia.MakeEncodingConfig()
	cdc, legacyCdc = gaia.MakeCodecs()
	config := sdk.GetConfig()
	config.Seal()
}

// Server represents the API server
type Server struct {
	Port int    `json:"port"`
	Node string `json:"node"`

	Version string `yaml:"version,omitempty"`
	Commit  string `yaml:"commit,omitempty"`
	Branch  string `yaml:"branch,omitempty"`

	// Server only supports file backend
	KeyringDir string `json:"keyring_dir"`
}

type KeyringWrapper struct {
	Keyring    keyring.Keyring
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
}

func (k *KeyringWrapper) Cleanup() {
	if k != nil {
		if k.pipeReader != nil {
			k.pipeReader.Close()
			k.pipeReader = nil
		}
		if k.pipeWriter != nil {
			k.pipeWriter.Close()
			k.pipeWriter = nil
		}
	}
}

func (s *Server) GetKeyringWrapper(password string) (*KeyringWrapper, error) {
	os.Stdin = nil
	passwordPipeReader, passwordPipeWriter := io.Pipe()
	keyring, err := keyring.New(sdk.KeyringServiceName(), "file", s.KeyringDir, passwordPipeReader)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			n, err := passwordPipeWriter.Write([]byte(password + "\n"))
			if n == 0 || err != nil {
				break
			}
		}
	}()

	var res = KeyringWrapper{
		Keyring:    keyring,
		pipeReader: passwordPipeReader,
		pipeWriter: passwordPipeWriter,
	}

	return &res, nil
}

// Router returns the router
func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/version", s.VersionHandler).Methods("GET")
	router.HandleFunc("/keys/list", s.GetKeys).Methods("POST")
	router.HandleFunc("/keys/create", s.PostKeys).Methods("POST")
	router.HandleFunc("/keys/get/{name}", s.GetKey).Methods("POST")
	router.HandleFunc("/keys/delete/{name}", s.DeleteKey).Methods("POST")
	router.HandleFunc("/tx/sign", s.Sign).Methods("POST")

	return router
}
