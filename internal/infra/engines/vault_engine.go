package engines

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/mburaksoran/insider-case/internal/app/config"
)

type VaultEngine struct {
	client  *api.Client
	Secrets *api.Secret
}

var vaultEngine *VaultEngine

func GetVaultEngine() *VaultEngine {
	return vaultEngine
}

func SetVaultEngine(cfg *config.AppConfig) (*VaultEngine, error) {
	if vaultEngine == nil {
		configuration := &api.Config{
			Address: cfg.Vault.Address,
		}
		client, err := api.NewClient(configuration)
		if err != nil {
			return nil, fmt.Errorf("error while creating vault client: %v", err)
		}
		client.SetToken(cfg.Vault.Token)
		return &VaultEngine{client: client}, nil
	}
	return vaultEngine, nil
}

func (v *VaultEngine) GetSecret(cfg *config.AppConfig) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(cfg.Vault.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot gathered from vault secrets: %v", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("cannot find any secrets: %s", cfg.Vault.Path)
	}
	v.Secrets = secret
	return secret.Data, nil
}
