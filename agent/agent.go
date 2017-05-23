package agent

import (
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
)

type Agent struct {
	n *config.NginxConfig
	m *mapping.Mapping
	c *consulapi.Client
	h *http.Client
}

func NewAgent(config *Config) (*Agent, error) {
	agent := &Agent{
		n: config.Nginx,
		m: config.Mapping,
	}

	//Create consul connection
	if client, err := consulApi.NewClient(consulApi.DefaultConfig()); err == nil {
		agent.c = client
	} else {
		return nil, err
	}

	//Create nginx connection
	agent.h = &http.Client{}

	return agent, nil
}

func (a *Agent) Run() {
	for name, e := range m {
		if entry, ok := e.(*mapping.ConsulEntry); !ok {
			continue
		} else {

		}
	}
}

func (a *Agent) Watch(service string) {
	health := a.c.Health()
	index := 0

}
