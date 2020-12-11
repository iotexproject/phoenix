// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package midware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/config"
	"github.com/pkg/errors"
)

// RateLimit ratelimit can be applied on a per-token、per-ip、per-url basis
func RateLimit(rateLimit config.RateLimit) chi.Middlewares {
	middlewares := chi.Middlewares{}
	if rateLimit.Enable && rateLimit.RequestLimit > 0 && rateLimit.WindowLength > 0 {
		keyFuncs := []httprate.KeyFunc{}
		for _, key := range rateLimit.LimitByKey {
			switch strings.ToLower(key) {
			case "ip", "client":
				keyFuncs = append(keyFuncs, httprate.KeyByIP)
			case "url", "endpoint":
				keyFuncs = append(keyFuncs, httprate.KeyByEndpoint)
			case "token", "user":
				keyFuncs = append(keyFuncs, func(r *http.Request) (string, error) {
					claims, ok := r.Context().Value(auth.TokenCtxKey).(*auth.Claims)
					if !ok {
						return "", errors.New("failed to get claims in context")
					}
					// trustor is the address that registers endpoint with us
					trustor, err := crypto.HexStringToPublicKey(claims.Issuer)
					if err != nil {
						return "", err
					}
					return trustor.Address().Hex(), nil
				})
			}
		}
		if len(keyFuncs) > 0 {
			middlewares = append(middlewares, httprate.Limit(rateLimit.RequestLimit, time.Duration(rateLimit.WindowLength)*time.Second, keyFuncs...))
		}
	}
	return middlewares
}
