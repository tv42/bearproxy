package main

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

type authenticatingProxy struct {
	token   string
	handler http.Handler
}

var _ http.Handler = (*authenticatingProxy)(nil)

func (t *authenticatingProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const authHeader = "Authorization"
	auth := r.Header.Get(authHeader)
	r.Header.Del(authHeader)
	const prefix = "Bearer "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		const status = http.StatusUnauthorized
		http.Error(w, http.StatusText(status), status)
		return
	}
	offered := auth[len(prefix):]
	if subtle.ConstantTimeCompare([]byte(t.token), []byte(offered)) != 1 {
		const status = http.StatusUnauthorized
		http.Error(w, http.StatusText(status), status)
		return
	}
	t.handler.ServeHTTP(w, r)
}
