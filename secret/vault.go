package secret

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// vaultProvider defines type for object to interact with vault api
type vaultProvider struct {
	path    string
	client  *api.Logical
	results map[string]map[string]string
}

// NewVaultProvider creates new instance of vault provider
func NewVaultProvider(token, addr, path string) (Provider, error) {
	config := &api.Config{
		Address: addr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	client.SetToken(token)

	return &vaultProvider{
		path:    path,
		client:  client.Logical(),
		results: make(map[string]map[string]string),
	}, nil
}

// Get retrieves a value from vault. 
// The actual key selected is determined by the value separated by the colon. 
// For example "database:password" will retrieve the key "password" from the path "database".
func (v *vaultProvider) Get(keyPath string) (string, error) {
	split := strings.Split(keyPath, ":")
	if len(split) == 1 {
		return "", errors.New("missing key value")
	}

	pathSecret := split[0]
	key := split[1]

	res, ok := v.results[pathSecret]
	if ok {
		val, ok := res[key]
		if !ok {
			return "", errors.New("key not found in cached data")
		}

		return val, nil
	}

	secret, err := v.client.Read(fmt.Sprintf("%s/data/%s", v.path, pathSecret))
	if err != nil {
		return "", fmt.Errorf("reading: %w", err)
	}

	if secret == nil {
		return "", errors.New("secret not found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid data in secret")
	}

	secrets := make(map[string]string)

	for k, v := range data {
		val, ok := v.(string)
		if !ok {
			return "", errors.New("secret value in data is not string")
		}

		secrets[k] = val
	}

	val, ok := secrets[key]
	if !ok {
		return "", errors.New("key not found in retrieved data")
	}

	v.results[pathSecret] = secrets

	return val, nil
}
