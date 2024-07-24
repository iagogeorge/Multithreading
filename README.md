# Multithreading

# Go CEP Lookup Service

This project implements a concurrent CEP (Postal Code) lookup service in Go, demonstrating skills in multithreading and API interaction.

## Project Overview

The project performs simultaneous requests to two different APIs to retrieve address information based on a provided CEP (Postal Code) and returns the result from the fastest responding API.

## Detailed Requirements

1. **Make concurrent requests** to the following APIs:
   - `https://brasilapi.com.br/api/cep/v1/{cep}`
   - `http://viacep.com.br/ws/{cep}/json/`
2. **Use the response from the fastest API** and discard the slower response.
3. **Display the address information** on the command line along with the name of the API that provided the response.
4. **Limit the response time to 1 second**. If neither API responds within this time frame, display a timeout error.
