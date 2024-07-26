package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Address struct {
	Logradouro string `json:"logradouro,omitempty"`
	Bairro     string `json:"bairro,omitempty"`
	Localidade string `json:"localidade,omitempty"`
	Uf         string `json:"uf,omitempty"`
	Api        string `json:"api,omitempty"`
}

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func fetchFromBrasilAPI(ctx context.Context, cep string) (Address, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for BrasilAPI: %v", err)
		return Address{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error fetching from BrasilAPI: %v", err)
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("BrasilAPI returned non-200 status code: %d", resp.StatusCode)
		return Address{}, fmt.Errorf("BrasilAPI returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from BrasilAPI: %v", err)
		return Address{}, err
	}

	var brasilApiResponse BrasilApiResponse
	err = json.Unmarshal(body, &brasilApiResponse)
	if err != nil {
		log.Printf("Error unmarshaling response from BrasilAPI: %v", err)
		return Address{}, err
	}

	address := Address{
		Logradouro: brasilApiResponse.Street,
		Bairro:     brasilApiResponse.Neighborhood,
		Localidade: brasilApiResponse.City,
		Uf:         brasilApiResponse.State,
		Api:        "BrasilAPI",
	}
	return address, nil
}

func fetchFromViaCEP(ctx context.Context, cep string) (Address, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for ViaCEP: %v", err)
		return Address{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error fetching from ViaCEP: %v", err)
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ViaCEP returned non-200 status code: %d", resp.StatusCode)
		return Address{}, fmt.Errorf("ViaCEP returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from ViaCEP: %v", err)
		return Address{}, err
	}

	var address Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		log.Printf("Error unmarshaling response from ViaCEP: %v", err)
		return Address{}, err
	}

	address.Api = "ViaCEP"
	return address, nil
}

func getAddress(cep string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	results := make(chan Address, 2)
	errors := make(chan error, 2)

	go func() {
		address, err := fetchFromBrasilAPI(ctx, cep)
		if err != nil {
			errors <- err
			return
		}
		results <- address
	}()

	go func() {
		address, err := fetchFromViaCEP(ctx, cep)
		if err != nil {
			errors <- err
			return
		}
		results <- address
	}()

	select {
	case res := <-results:
		log.Printf("Fastest API: %s", res.Api)
		log.Printf("Address: %+v\n", res)
	case err := <-errors:
		log.Printf("Error fetching address: %v", err)
	case <-ctx.Done():
		log.Println("Error: Timeout exceeded while fetching addresses")
	}
}

func main() {
	cep := "01153000"
	getAddress(cep)
}
