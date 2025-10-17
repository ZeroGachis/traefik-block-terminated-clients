package traefik_block_terminated_clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFowardRequestIfNoUsernameQueryParamExists(t *testing.T) {
	cfg := CreateConfig()
	cfg.Usernames = "client1,client2"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler, err := New(ctx, next, cfg, "sw-block-terminated-clients-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Errorf("Got status code %d", recorder.Code)
	}
}

func TestFowardRequestIfUsernameIsNotBlocked(t *testing.T) {
	cfg := CreateConfig()
	cfg.Usernames = "client1,client2"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler, err := New(ctx, next, cfg, "sw-block-terminated-clients-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost?username=client3", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Errorf("Got status code %d", recorder.Code)
	}
}

func TestBlockRequestIfUsernameIsBlocked(t *testing.T) {
	cfg := CreateConfig()
	cfg.Usernames = "client1,client2"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler, err := New(ctx, next, cfg, "sw-block-terminated-clients-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost?username=client2", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusNoContent {
		t.Errorf("Got status code %d", recorder.Code)
	}
}
