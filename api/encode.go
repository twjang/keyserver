package api

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// EncodeResponse nolint
type EncodeResponse struct {
	TxBytes string `json:"txbytes"`
	TxID    string `json:"txid"`
}

// EncodeTx - no-lint
func (s *Server) EncodeTx(w http.ResponseWriter, r *http.Request) {
	var stdTx auth.StdTx
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	err = cdc.UnmarshalJSON(body, &stdTx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(cdc.MustMarshalJSON(EncodeResponse{
		TxBytes: hex.EncodeToString(txBytes),
		TxID:    hex.EncodeToString(tmhash.Sum(txBytes)),
	}))
	return
}
