// Package traefik_block_terminated_clients a Traefik plugin to block HTTP requests from terminated clients.
package traefik_block_terminated_clients

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	Username []string `json:"username,omitempty"`
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
	fmt.Println("Blocklist: ", config.Username) //nolint
	for _, username := range config.Username {
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
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	a.next.ServeHTTP(rw, req)
}
