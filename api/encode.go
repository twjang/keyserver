package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// EncodeBody is the body for a encode request
type EncodeBody struct {
	Tx json.RawMessage `json:"tx"`
}

// EncodeResponse is the response for a encode request
type EncodeResponse struct {
	Tx string `json:"tx"`
}

// Encode encodes json tx to binary format
func (s *Server) Encode(w http.ResponseWriter, r *http.Request) {
	var m EncodeBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tx, err := encodingConfig.TxConfig.TxJSONDecoder()(m.Tx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	bytes, err := encodingConfig.TxConfig.TxEncoder()(tx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(newError(err).marshal())
		return
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	output := EncodeResponse{
		Tx: encoded,
	}

	out, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(newError(err).marshal())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return
}
