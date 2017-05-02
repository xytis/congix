package config

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/xytis/congix/helper"
)

// ConsulConfig contains the configuration information necessary to
// communicate with a Consul Agent.
type ConsulConfig struct {
	// Addr is the address of the local Consul agent
	Addr string `mapstructure:"address"`

	// Timeout is used by Consul HTTP Client
	Timeout time.Duration `mapstructure:"timeout"`

	// Token is used to provide a per-request ACL token. This options overrides
	// the agent's default token
	Token string `mapstructure:"token"`

	// Auth is the information to use for http access to Consul agent
	Auth string `mapstructure:"auth"`

	// EnableSSL sets the transport scheme to talk to the Consul agent as https
	EnableSSL *bool `mapstructure:"ssl"`

	// VerifySSL enables or disables SSL verification when the transport scheme
	// for the consul api client is https
	VerifySSL *bool `mapstructure:"verify_ssl"`

	// CAFile is the path to the ca certificate used for Consul communication
	CAFile string `mapstructure:"ca_file"`

	// CertFile is the path to the certificate for Consul communication
	CertFile string `mapstructure:"cert_file"`

	// KeyFile is the path to the private key for Consul communication
	KeyFile string `mapstructure:"key_file"`
}

// DefaultConsulConfig() returns the canonical defaults for the Congix
// `consul` configuration.
func DefaultConsulConfig() *ConsulConfig {
	return &ConsulConfig{
		EnableSSL: helper.BoolToPtr(false),
		VerifySSL: helper.BoolToPtr(false),
		Timeout:   5 * time.Second,
	}
}

// Merge merges two Consul Configurations together.
func (a *ConsulConfig) Merge(b *ConsulConfig) *ConsulConfig {
	result := a.Copy()

	if b.Addr != "" {
		result.Addr = b.Addr
	}
	if b.Timeout != 0 {
		result.Timeout = b.Timeout
	}
	if b.Token != "" {
		result.Token = b.Token
	}
	if b.Auth != "" {
		result.Auth = b.Auth
	}
	if b.EnableSSL != nil {
		result.EnableSSL = helper.BoolToPtr(*b.EnableSSL)
	}
	if b.VerifySSL != nil {
		result.VerifySSL = helper.BoolToPtr(*b.VerifySSL)
	}
	if b.CAFile != "" {
		result.CAFile = b.CAFile
	}
	if b.CertFile != "" {
		result.CertFile = b.CertFile
	}
	if b.KeyFile != "" {
		result.KeyFile = b.KeyFile
	}
	return result
}

// ApiConfig() returns a usable Consul config that can be passed directly to
// hashicorp/consul/api.  NOTE: datacenter is not set
func (c *ConsulConfig) ApiConfig() (*consul.Config, error) {
	config := consul.DefaultConfig()
	if c.Addr != "" {
		config.Address = c.Addr
	}
	if c.Token != "" {
		config.Token = c.Token
	}
	if c.Timeout != 0 {
		config.HttpClient.Timeout = c.Timeout
	}
	if c.Auth != "" {
		var username, password string
		if strings.Contains(c.Auth, ":") {
			split := strings.SplitN(c.Auth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = c.Auth
		}

		config.HttpAuth = &consul.HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}
	if c.EnableSSL != nil && *c.EnableSSL {
		config.Scheme = "https"
		tlsConfig := consul.TLSConfig{
			Address:  config.Address,
			CAFile:   c.CAFile,
			CertFile: c.CertFile,
			KeyFile:  c.KeyFile,
		}
		if c.VerifySSL != nil {
			tlsConfig.InsecureSkipVerify = !*c.VerifySSL
		}

		tlsClientCfg, err := consul.SetupTLSConfig(&tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("error creating tls client config for consul: %v", err)
		}
		config.HttpClient.Transport = &http.Transport{
			TLSClientConfig: tlsClientCfg,
		}
	}

	return config, nil
}

// Copy returns a copy of this Consul config.
func (c *ConsulConfig) Copy() *ConsulConfig {
	if c == nil {
		return nil
	}

	nc := new(ConsulConfig)
	*nc = *c

	// Copy the bools
	if nc.EnableSSL != nil {
		nc.EnableSSL = helper.BoolToPtr(*nc.EnableSSL)
	}
	if nc.VerifySSL != nil {
		nc.VerifySSL = helper.BoolToPtr(*nc.VerifySSL)
	}

	return nc
}
