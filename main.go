package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func fetchAddress(ctx context.Context, url string, ch chan<- result) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- result{err: err}
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- result{err: err}
		return
	}
	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		ch <- result{err: err}
		return
	}

	ch <- result{address: address}
}

type result struct {
	address Address
	err     error
}

func main() {
	cep := "01153000"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan result, 2)

	go fetchAddress(ctx, fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), ch)
	go fetchAddress(ctx, fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), ch)

	select {
	case res := <-ch:
		if res.err != nil {
			log.Fatalf("Error fetching address: %v", res.err)
		} else {
			fmt.Printf("Address: %+v\n", res.address)
		}
	case <-ctx.Done():
		log.Fatalf("Timeout exceeded")
	}
}
