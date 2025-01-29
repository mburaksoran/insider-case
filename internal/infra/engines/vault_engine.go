package engines

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mburaksoran/insider-case/internal/app/config"
)

type VaultEngine struct {
	Client  *api.Client
	Secrets *api.Secret
	config  *config.VaultConfig
}

var vaultEngine *VaultEngine

func GetVaultEngine() *VaultEngine {
	return vaultEngine
}

func SetVaultEngine() (*VaultEngine, error) {
	if vaultEngine == nil {
		vaultConfig, err := config.InitVaultConfig("./configuration/config.yml")
		if err != nil {
			log.Fatalf("Vault config error: %v", err)
		}

		configuration := &api.Config{
			Address:    vaultConfig.Address,
			HttpClient: &http.Client{Timeout: 30 * time.Second},
		}
		client, err := api.NewClient(configuration)
		if err != nil {
			return nil, fmt.Errorf("error while creating vault client: %v", err)
		}
		client.SetToken(vaultConfig.Token)
		vaultEngine = &VaultEngine{Client: client, config: vaultConfig}
		return &VaultEngine{Client: client, config: vaultConfig}, nil
	}
	return vaultEngine, nil
}

func (v *VaultEngine) GetSecret() (map[string]interface{}, error) {
	secret, err := v.Client.Logical().Read(v.config.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot gathered from vault secrets: %v", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("cannot find any secrets: %s", v.config.Path)
	}
	v.Secrets = secret
	return secret.Data, nil
}
