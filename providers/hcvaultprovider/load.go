package hcvaultprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	vault "github.com/hashicorp/vault/api"
)

func (p *provider) Load() error {
	config := vault.DefaultConfig()

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	client.SetToken("dev-only-token")

	secret, err := client.KVv2("secret").Get(context.Background(), "my-secret-password")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	log.Print(secret)

	return nil
}

type fetchSecretResponse struct {
	Secrets []vaultSecret
}
type vaultSecret struct {
	Name    string
	Version struct {
		Value string
	}
}

func (p *provider) fetchSecrets() (*fetchSecretResponse, error) {
	token, err := p.getAuthToken()
	if err != nil {
		// TODO
		return nil, err
	}

	req, err := http.NewRequest("GET", os.Getenv("VAULT_SECRETS_LOCATION"), nil)
	if err != nil {
		// TODO
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := p.client.Do(req)
	if err != nil {
		// TODO
		return nil, err
	}

	if resp.StatusCode >= 300 {
		// TODO
		return nil, errors.New("bad response code from vault secrets get value request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// TODO
		return nil, err
	}

	var result fetchSecretResponse
	json.Unmarshal(body, &result)

	return &result, nil
}

func (p *provider) getAuthToken() (string, error) {

	data, err := json.Marshal(map[string]string{
		"audience":      "https://api.hashicorp.cloud",
		"grant_type":    "client_credentials",
		"client_id":     os.Getenv("HCP_CLIENT_ID"),
		"client_secret": os.Getenv("HCP_CLIENT_SECRET"),
	})
	if err != nil {
		// TODO
		return "", err
	}

	buf := bytes.NewReader(data)
	req, err := http.NewRequest("POST", "https://auth.hashicorp.com/oauth/token", buf)
	if err != nil {
		// TODO
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		// TODO
		return "", err
	}

	if resp.StatusCode >= 300 {
		// TODO
		return "", errors.New("bad response code from vault secrets token request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// TODO
		return "", err
	}

	var result map[string]string
	json.Unmarshal(body, &result)

	return result["access_token"], nil
}
