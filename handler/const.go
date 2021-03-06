// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package handler

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/iotexproject/phoenix/auth"
	"github.com/iotexproject/phoenix/json"
	"github.com/pkg/errors"
)

var (
	ErrorPermissionDenied = errors.New("You don't have permission for this")
	ErrorBodyEmpty        = errors.New("Body must be set")
	ErrorStoreCtx         = errors.New("Failed to get store in context")
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type podObject struct {
	Name string `json:"name"`
}

type registerObject struct {
	Name     string `json:"name"`
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
	Token    string `json:"token"`
}

func (r *registerObject) Store() auth.Store {
	return auth.NewStore(r.Name, r.Region, r.Endpoint, r.Key, r.Token)
}

func decodeAndCloseRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	return err
}

func renderJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
