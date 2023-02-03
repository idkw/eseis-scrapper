package eseis

import (
	"fmt"
	"github.com/caarlos0/env/v7"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// EseisClient is a client for the Eseis API
type EseisClient struct {
	config      *Config
	accessToken *authToken
}

// Config is a configuration struct to build an EseisClient
type Config struct {
	ClientId string `env:"ESEIS_CLIENT_ID,required"`
	Username string `env:"ESEIS_USERNAME,required"`
	Password string `env:"ESEIS_PASSWORD,required"`
	BaseURL  string `env:"ESEIS_BASE_URL,required" envDefault:"https://sergic-api-prod.sergic.com"`
}

// NewEseisClient creates a new EseisClient or returns an error
func NewEseisClient() (*EseisClient, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create eseis client config: %w", err)
	}
	return &EseisClient{config: config}, nil
}

// NewEseisClientFatal creates a new EseisClient or panics if an errors occurs
func NewEseisClientFatal() *EseisClient {
	client, err := NewEseisClient()
	if err != nil {
		logrus.Fatalf("failed to create eseis client: %s", err)
	}
	return client
}

func newConfig() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment: %w", err)
	}
	return config, nil
}

func (e *EseisClient) buildURL(path string) string {
	return e.config.BaseURL + path
}

func (e *EseisClient) checkAuthenticated() error {
	if e.accessToken == nil || e.accessToken.expiresAt.Add(-10*time.Minute).Before(time.Now()) {
		token, err := e.Authenticate()
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
		e.accessToken = token
	}
	return nil
}

func (e *EseisClient) setAuthentication(request *http.Request) {
	if e.accessToken != nil {
		request.Header.Set("Authorization", "Bearer "+e.accessToken.accessToken)
	}
}
