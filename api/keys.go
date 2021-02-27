package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
	"github.com/gorilla/mux"
)

var notFoundMsg = "The specified item could not be found in the keyring"

// GetKeyBody describes input format of some rest api endpoints
type GetKeyBody struct {
	Password string `json:"password"`
}

// GetKeys is the handler for the POST /keys
func (s *Server) GetKeys(w http.ResponseWriter, r *http.Request) {
	var m GetKeyBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	keyringWrapper, err := s.GetKeyringWrapper(m.Password)
	defer keyringWrapper.Cleanup()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	k := keyringWrapper.Keyring
	infos, err := k.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	if len(infos) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	keysOutput, err := keyring.Bech32KeysOutput(infos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	out, err := json.Marshal(keysOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

// AddNewKey is the necessary data for adding a new key
type AddNewKey struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Mnemonic string `json:"mnemonic,omitempty"`
	Account  int    `json:"account,string,omitempty"`
	Index    int    `json:"index,string,omitempty"`
}

// Marshal - no-lint
func (ak AddNewKey) Marshal() []byte {
	out, err := json.Marshal(ak)
	if err != nil {
		panic(err)
	}
	return out
}

// PostKeys is the handler for the POST /keys
func (s *Server) PostKeys(w http.ResponseWriter, r *http.Request) {
	var m AddNewKey

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	if m.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("must include both password and name with request")).marshal())
		return
	}

	// if mnemonic is empty, generate one
	mnemonic := m.Mnemonic
	if mnemonic == "" {
		fullFundraiserPath := sdk.FullFundraiserPath
		_, mnemonic, _ = keyring.NewInMemory().NewMnemonic("inmemorykey", keyring.English, fullFundraiserPath, hd.Secp256k1)
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("invalid mnemonic")).marshal())
		return
	}

	if m.Account < 0 || m.Account > maxValidAccountValue {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("invalid account number")).marshal())
		return
	}

	if m.Index < 0 || m.Index > maxValidIndexalue {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("invalid index number")).marshal())
		return
	}

	keyringWrapper, err := s.GetKeyringWrapper(m.Password)
	defer keyringWrapper.Cleanup()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}
	k := keyringWrapper.Keyring
	_, err = k.Key(m.Name)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(fmt.Errorf("key %s already exists", m.Name)).marshal())
		return
	}

	if err != nil && err.Error() != notFoundMsg {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
	}

	account := uint32(m.Account)
	index := uint32(m.Index)

	hdpath := fmt.Sprintf("44'/118'/%d'/0/%d", account, index)
	info, err := k.NewAccount(m.Name, mnemonic, keyring.DefaultBIP39Passphrase, hdpath, hd.Secp256k1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	keyOutput, err := keyring.Bech32KeyOutput(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	keyOutput.Mnemonic = mnemonic

	out, err := json.Marshal(keyOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

// GetKey is the handler for the POST /keys/{name}
func (s *Server) GetKey(w http.ResponseWriter, r *http.Request) {
	var m GetKeyBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	keyringWrapper, err := s.GetKeyringWrapper(m.Password)
	defer keyringWrapper.Cleanup()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}
	k := keyringWrapper.Keyring

	vars := mux.Vars(r)
	name := vars["name"]
	bechPrefix := r.URL.Query().Get("bech")

	if bechPrefix == "" {
		bechPrefix = "acc"
	}

	bechKeyOut, err := getBechKeyOut(bechPrefix)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	info, err := k.Key(name)

	if err != nil {
		if err.Error() == notFoundMsg {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(newError(err).marshal())
		return
	}

	keyOutput, err := bechKeyOut(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	out, err := json.Marshal(keyOutput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}

type bechKeyOutFn func(keyInfo keyring.Info) (keyring.KeyOutput, error)

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case "acc":
		return keyring.Bech32KeyOutput, nil
	case "val":
		return keyring.Bech32ValKeyOutput, nil
	case "cons":
		return keyring.Bech32ConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}

// DeleteKeyBody request
type DeleteKeyBody struct {
	Password string `json:"password"`
}

// Marshal - no-lint
func (u DeleteKeyBody) Marshal() []byte {
	out, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return out
}

// DeleteKey is the handler for the DELETE /keys/{name}
func (s *Server) DeleteKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	var m DeleteKeyBody
	var err error

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	keyringWrapper, err := s.GetKeyringWrapper(m.Password)
	defer keyringWrapper.Cleanup()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}
	k := keyringWrapper.Keyring

	err = k.Delete(name)

	if err != nil {
		if err.Error() == notFoundMsg {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
