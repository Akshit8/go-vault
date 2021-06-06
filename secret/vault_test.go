package secret_test

import (
	"fmt"
	"testing"

	"github.com/Akshit8/go-vault/secret"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// vaultClient defines to interact with vault during the test
type vaultClient struct {
	Token   string
	Address string
	path    string
	Client  *api.Client
}

func newVaultClient(t *testing.T) *vaultClient {
	token := "secrettoken"
	addr := "http://0.0.0.0:8300"
	path := "/secret"

	config := &api.Config{
		Address: addr,
	}

	client, err := api.NewClient(config)
	require.NoError(t, err)

	client.SetToken(token)

	return &vaultClient{
		Token:   token,
		Address: addr,
		path:    path,
		Client:  client,
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	type output struct {
		res     string
		withErr bool
	}

	tests := []struct {
		name   string
		setup  func(*testing.T, *vaultClient)
		input  string
		output output
	}{
		{
			"OK: from vault",
			func(t *testing.T, vc *vaultClient) {
				_, err := vc.Client.Logical().Write(
					fmt.Sprintf("%s/data/ok", vc.path),
					map[string]interface{}{
						"data": map[string]interface{}{
							"key1": "value1",
							"key2": "value2",
						},
					},
				)
				require.NoError(t, err)
			},
			"/ok:key1",
			output{
				res: "value1",
			},
		},
		{
			"OK: cached",
			func(t *testing.T, vc *vaultClient) {},
			"/ok:key2",
			output{
				res: "value2",
			},
		},
		{
			"ERR: key not found in cached data",
			func(t *testing.T, vc *vaultClient) {},
			"/ok:three",
			output{
				withErr: true,
			},
		},
		{
			"ERR: secret not found",
			func(t *testing.T, vc *vaultClient) {},
			"/not:found",
			output{
				withErr: true,
			},
		},
		{
			"ERR: key not found in retreived data from vault",
			func(t *testing.T, vc *vaultClient) {
				_, err := vc.Client.Logical().Write(
					fmt.Sprintf("%s/data/err", vc.path),
					map[string]interface{}{
						"data": map[string]interface{}{
							"hello": "world",
						},
					},
				)
				require.NoError(t, err)
			},
			"/err:random",
			output{
				withErr: true,
			},
		},
	}

	// Using central provider as need to test the cache logic.

	client := newVaultClient(t)
	provider, err := secret.NewVaultProvider(client.Token, client.Address, client.path)
	require.NoError(t, err)

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			// Not calling t.Parallel() because vault.Provider is not goroutine safe.

			test.setup(t, client)

			res, err := provider.Get(test.input)

			isErr := err != nil
			require.Equal(t, test.output.withErr, isErr)
			require.Equal(t, test.output.res, res)
		})
	}
}
