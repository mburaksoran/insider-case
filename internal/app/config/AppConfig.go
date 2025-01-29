package config

import (
	"encoding/json"
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// InitFromConfigFile is to populate VaultConfig with values from file
func InitVaultConfig(path string) (*VaultConfig, error) {
	// read vault yaml
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// map to vault config struct
	var vaultConfig VaultConfig
	if err := yaml.Unmarshal(data, &vaultConfig); err != nil {
		return nil, err
	}
	return &vaultConfig, nil
}

// AppConfig is a structure, instance of which will be populated from provided yaml config

type Data struct {
	AppConfig `json:"data"`
}
type AppConfig struct {
	Redis      RedisConfig      `json:"redis_config"`
	Postgres   PostgresConfig   `json:"postgres_config"`
	HttpClient HttpClientConfig `json:"http_client_config"`
	Vault      VaultConfig      `yaml:"vault_path"`
	WorkerPool WorkerPoolConfig `json:"worker_pool"`
}

type HttpClientConfig struct {
	Url    string `json:"http_client_url"`
	ApiKey string `json:"api_key"`
}

type RedisConfig struct {
	Host     string `json:"redis_host"`
	Port     string `json:"redis_port"`
	Password string `json:"redis_password"`
}

type PostgresConfig struct {
	SqlDatabaseName string `json:"sql_database_name"`
	SqlHost         string `json:"sql_host"`
	SqlPassword     string `json:"sql_password"`
	SqlPort         string `json:"sql_port"`
	SqlUser         string `json:"sql_user"`
	SqlSslMode      string `json:"sql_ssl_mode"`
}
type WorkerPoolConfig struct {
	WorkerCount int `json:"worker_count"`
	ChannelSize int `json:"channel_size"`
}
type VaultConfig struct {
	Address string `yaml:"vault_address"`
	Token   string `yaml:"vault_token"`
	Path    string `yaml:"vault_path"`
}

func (v *AppConfig) MapVaultSecretToConfig(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data["data"])
	if err != nil {
		return errors.New("failed to marshal data to JSON")
	}
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return errors.New("failed to unmarshal JSON to struct")
	}

	return nil
}

func (v *AppConfig) ChangeVaultConfig(vaultConfig VaultConfig) {
	v.Vault = vaultConfig
}
