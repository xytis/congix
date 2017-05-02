package config

// NginxConfig contains information to connect to nginx+ API
type NginxConfig struct {
	// Addr is the address of nginx+ server to control
	Addr string `mapstructure:"address"`

	// StatusEndpoint is the path to status endpoint
	StatusEndpoint string `mapstructure:"status_endpoint"`

	// UpstreamEndpoint is the path to upstream_conf endpoint
	UpstreamEndpoint string `mapstructure:"upstream_endpoint"`
}

// DefaultNginxConfig() returns the canonical defaults for the Congix
// `nginx` configuration.
func DefaultNginxConfig() *NginxConfig {
	return &NginxConfig{
		Addr:             "http://localhost:8080",
		StatusEndpoint:   "status",
		UpstreamEndpoint: "upstream_conf",
	}
}

// Merge merges two Nginx Configurations together.
func (a *NginxConfig) Merge(b *NginxConfig) *NginxConfig {
	result := a.Copy()

	if b.Addr != "" {
		result.Addr = b.Addr
	}
	if b.StatusEndpoint != "" {
		result.StatusEndpoint = b.StatusEndpoint
	}
	if b.UpstreamEndpoint != "" {
		result.UpstreamEndpoint = b.UpstreamEndpoint
	}
	return result
}

// Copy returns a copy of this Nginx config.
func (c *NginxConfig) Copy() *NginxConfig {
	if c == nil {
		return nil
	}

	nc := new(NginxConfig)
	*nc = *c

	return nc
}
