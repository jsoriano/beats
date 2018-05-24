package jolokia

import (
	"github.com/elastic/beats/libbeat/autodiscover"
	"github.com/elastic/beats/libbeat/autodiscover/template"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/bus"
	"github.com/elastic/beats/libbeat/common/cfgwarn"
)

func init() {
	autodiscover.Registry.AddProvider("jolokia", AutodiscoverBuilder)
}

// DiscoveryProber implements discovery probes
type DiscoveryProber interface {
	Start()
	Stop()
	Events() <-chan Event
}

// Provider is the Jolokia Discovery autodiscover provider
type Provider struct {
	config    *Config
	bus       bus.Bus
	builders  autodiscover.Builders
	appenders autodiscover.Appenders
	templates *template.Mapper
	discovery DiscoveryProber
}

// AutodiscoverBuilder builds a Jolokia Discovery autodiscover provider, it fails if
// there is some problem with the configuration
func AutodiscoverBuilder(bus bus.Bus, c *common.Config) (autodiscover.Provider, error) {
	cfgwarn.Experimental("The Jolokia Discovery autodiscover is experimental")

	config, err := getConfig(c)
	if err != nil {
		return nil, err
	}

	discovery := &Discovery{
		Interfaces: config.Interfaces,
	}

	mapper, err := template.NewConfigMapper(config.Templates)
	if err != nil {
		return nil, err
	}

	builders, err := autodiscover.NewBuilders(config.Builders, false)
	if err != nil {
		return nil, err
	}

	appenders, err := autodiscover.NewAppenders(config.Appenders)
	if err != nil {
		return nil, err
	}

	return &Provider{
		bus:       bus,
		templates: mapper,
		builders:  builders,
		appenders: appenders,
		discovery: discovery,
	}, nil
}

// Start starts autodiscover provider
func (p *Provider) Start() {
	p.discovery.Start()
	go func() {
		for event := range p.discovery.Events() {
			p.publish(event.BusEvent())
		}
	}()
}

func (p *Provider) publish(event bus.Event) {
	if config := p.templates.GetConfig(event); config != nil {
		event["config"] = config
	} else if config := p.builders.GetConfig(event); config != nil {
		event["config"] = config
	}

	p.appenders.Append(event)
	p.bus.Publish(event)
}

// Stop stops autodiscover provider
func (p *Provider) Stop() {
	p.discovery.Stop()
}

// String returns the name of the provider
func (p *Provider) String() string {
	return "jolokia"
}
