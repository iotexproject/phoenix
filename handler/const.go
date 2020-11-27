package handler

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/iotexproject/phoenix-gem/json"
	"github.com/pkg/errors"
)

var (
	ErrorPermissionDenied = errors.New("Your don't have permission for this")
	ErrorBodyEmpty        = errors.New("Body must be set")
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

func decodeJSON(r *http.Request, v interface{}) error {
	defer io.Copy(ioutil.Discard, r.Body)
	return json.NewDecoder(r.Body).Decode(v)
}

func renderJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
