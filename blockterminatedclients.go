// Package traefik_block_terminated_clients a Traefik plugin to block HTTP requests from terminated clients.
package traefik_block_terminated_clients

import (
	"context"
	"net/http"
	"strings"
)

type Config struct {
	Usernames string `json:"username,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type Plugin struct {
	next      http.Handler
	name      string
	usernames map[string]bool
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	usernames := make(map[string]bool)

	for _, username := range strings.Split(config.Usernames, ",") {
		usernames[username] = true
	}

	return &Plugin{
		next:      next,
		name:      name,
		usernames: usernames,
	}, nil
}

func (a *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	_, usernameIsBlocked := a.usernames[username]
	if username != "" && usernameIsBlocked {
		rw.WriteHeader(http.StatusNoContent)

		return
	}

	a.next.ServeHTTP(rw, req)
}
