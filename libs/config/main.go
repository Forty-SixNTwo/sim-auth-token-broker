package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	EnvKey        = "ENV"
	EnvDev        = "development"
	EnvProd       = "production"
	PortKey       = "PORT"
	PrefixMapPath = "PREFIX_MAP_PATH"
	SigningKey    = "SIGNING_KEY"
)

type Telco struct {
	BaseURL      string `yaml:"base_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type BrokerConfig struct {
	PrefixMap  map[string]Telco
	SigningKey string
	ListenAddr string
}

type TelcoConfig struct {
	TelcoKeyID        string
	TelcoClientID     string
	TelcoClientSecret string
	TelcoIssuerURL    string
}

func LoadTelcoConfig(keyIDKey, clientIDKey, clientSecretKey, issuerKey string) (*TelcoConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env: %w", err)
	}
	env := os.Getenv(EnvKey)
	if env == EnvProd {
		if err := fetchSecrets(); err != nil {
			return nil, fmt.Errorf("fetching secrets: %w", err)
		}
	}

	kid, err := require(keyIDKey)
	if err != nil {
		return nil, err
	}
	cid, err := require(clientIDKey)
	if err != nil {
		return nil, err
	}
	secret, err := require(clientSecretKey)
	if err != nil {
		return nil, err
	}
	issuer, err := require(issuerKey)
	if err != nil {
		return nil, err
	}

	return &TelcoConfig{
		TelcoKeyID:        kid,
		TelcoClientID:     cid,
		TelcoClientSecret: secret,
		TelcoIssuerURL:    issuer,
	}, nil
}

func LoadBrokerConfig() (*BrokerConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env: %w", err)
	}
	env := os.Getenv(EnvKey)
	if env == EnvProd {
		if err := fetchSecrets(); err != nil {
			return nil, fmt.Errorf("fetching secrets: %w", err)
		}
	}

	path, err := require(PrefixMapPath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading prefix map: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("prefix map file is empty")
	}
	var raw struct {
		Prefixes map[string]Telco `yaml:"prefixes"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing prefix map: %w", err)
	}

	for prefix, telco := range raw.Prefixes {
		cid, err := require(telco.ClientID)
		if err != nil {
			return nil, fmt.Errorf("missing env for telco %s client_id: %w", prefix, err)
		}
		secret, err := require(telco.ClientSecret)
		if err != nil {
			return nil, fmt.Errorf("missing env for telco %s client_secret: %w", prefix, err)
		}
		telco.ClientID = cid
		telco.ClientSecret = secret
		raw.Prefixes[prefix] = telco
	}
	skey, err := require(SigningKey)
	if err != nil {
		return nil, err
	}
	port, err := require(PortKey)
	if err != nil {
		return nil, err
	}

	return &BrokerConfig{
		PrefixMap:  raw.Prefixes,
		SigningKey: skey,
		ListenAddr: port,
	}, nil
}

func require(key string) (string, error) {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v, nil
	}
	return "", fmt.Errorf("environment variable %s is required", key)
}
