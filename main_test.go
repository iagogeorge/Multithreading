package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mockServer(response string, delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
}

func TestFetchAddress(t *testing.T) {
	brasilAPIServer := mockServer(`{"cep": "01153000", "logradouro": "Rua Teste", "bairro": "Bairro Teste", "localidade": "Cidade Teste", "uf": "ST"}`, 500*time.Millisecond)
	viaCepServer := mockServer(`{"cep": "01153000", "logradouro": "Avenida Teste", "bairro": "Bairro Teste", "localidade": "Cidade Teste", "uf": "ST"}`, 200*time.Millisecond)
	defer brasilAPIServer.Close()
	defer viaCepServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan result, 2)

	go fetchAddress(ctx, brasilAPIServer.URL, ch)
	go fetchAddress(ctx, viaCepServer.URL, ch)

	select {
	case res := <-ch:
		if res.err != nil {
			t.Fatalf("Error fetching address: %v", res.err)
		}
		if res.address.Logradouro != "Avenida Teste" {
			t.Fatalf("Expected Avenida Teste but got %s", res.address.Logradouro)
		}
	case <-ctx.Done():
		t.Fatalf("Timeout exceeded")
	}
}

func TestTimeout(t *testing.T) {
	slowServer := mockServer(`{"cep": "01153000", "logradouro": "Rua Teste", "bairro": "Bairro Teste", "localidade": "Cidade Teste", "uf": "ST"}`, 2*time.Second)
	defer slowServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan result, 2)

	go fetchAddress(ctx, slowServer.URL, ch)
	go fetchAddress(ctx, slowServer.URL, ch)

	select {
	case res := <-ch:
		if res.err == nil {
			t.Fatalf("Expected timeout but got address: %+v", res.address)
		}
	case <-ctx.Done():
		t.Log("Timeout exceeded as expected")
	}
}
