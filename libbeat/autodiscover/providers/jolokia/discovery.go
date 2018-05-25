package jolokia

import (
	"encoding/json"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/common/bus"
	s "github.com/elastic/beats/libbeat/common/schema"
	c "github.com/elastic/beats/libbeat/common/schema/mapstriface"
	"github.com/elastic/beats/libbeat/logp"
)

// Jolokia Discovery query
// {
//   "type": "query"
// }
//
// Example Jolokia Discovery response
// {
//   "agent_version": "1.5.0",
//   "agent_id": "172.18.0.2-7-1322ae88-servlet",
//   "server_product": "tomcat",
//   "type": "response",
//   "server_vendor": "Apache",
//   "server_version": "7.0.86",
//   "secured": false,
//   "url": "http://172.18.0.2:8778/jolokia"
// }
//
// Example discovery probe with socat
//
//   echo '{"type":"query"}' | sudo socat STDIO UDP4-DATAGRAM:239.192.48.84:24884,interface=br0 | jq .
//

// Message contains the information of a Jolokia Discovery message
var messageSchema = s.Schema{
	"agent": s.Object{
		"id":      c.Str("agent_id"),
		"version": c.Str("agent_version", s.Optional),
	},
	"secured": c.Bool("secured", s.Optional),
	"server": s.Object{
		"product": c.Str("server_product", s.Optional),
		"vendor":  c.Str("server_vendor", s.Optional),
		"version": c.Str("server_version", s.Optional),
	},
	"url": c.Str("url"),
}

// Event is a Jolokia Discovery event
type Event struct {
	Type    string
	Message common.MapStr
}

// BusEvent converts a Jolokia Discovery event to a autodiscover bus event
func (e *Event) BusEvent() bus.Event {
	event := bus.Event{
		e.Type:    true,
		"jolokia": e.Message,
		"meta": common.MapStr{
			"jolokia": e.Message,
		},
	}
	return event
}

// Instance is a discovered Jolokia instance, it keeps information of the
// last probe it replied
type Instance struct {
	LastSeen      time.Time
	LastInterface *InterfaceConfig
	Message       common.MapStr
}

// Discovery controls the Jolokia Discovery probes
type Discovery struct {
	sync.Mutex

	Interfaces []InterfaceConfig

	instances map[string]*Instance

	events chan Event
	stop   chan struct{}
}

// Start starts discovery probes
func (d *Discovery) Start() {
	d.instances = make(map[string]*Instance)
	d.events = make(chan Event)
	d.stop = make(chan struct{})
	go d.run()
}

// Stop stops discovery probes
func (d *Discovery) Stop() {
	d.stop <- struct{}{}
	close(d.events)
}

// Events returns a channel with the events of started and stopped Jolokia
// instances discovered
func (d *Discovery) Events() <-chan Event {
	return d.events
}

func (d *Discovery) run() {
	var cases []reflect.SelectCase
	for _, i := range d.Interfaces {
		ticker := time.NewTicker(i.Interval)
		defer ticker.Stop()
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ticker.C),
		})
	}

	// As a last case, place the stop channel so the loop can be gracefuly stopped
	stopIndex := len(cases)
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(d.stop),
	})

	for {
		chosen, _, _ := reflect.Select(cases)
		if chosen == stopIndex {
			return
		}
		d.sendProbe(d.Interfaces[chosen])
		d.checkStopped()
	}
}

func (d *Discovery) interfaces(name string) ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	if name == "any" || name == "*" {
		return interfaces, nil
	}

	var matching []net.Interface
	for _, i := range interfaces {
		if matchInterfaceName(name, i.Name) {
			matching = append(matching, i)
		}
	}
	return matching, nil
}

func matchInterfaceName(name, candidate string) bool {
	if strings.HasSuffix(name, "*") {
		return strings.HasPrefix(candidate, strings.TrimRight(name, "*"))
	}
	return name == candidate
}

func getIPv4Addr(i net.Interface) (net.IP, error) {
	addrs, err := i.Addrs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get addresses for "+i.Name)
	}
	for _, a := range addrs {
		if ip, _, err := net.ParseCIDR(a.String()); err == nil && ip != nil {
			if ipv4 := ip.To4(); ipv4 != nil {
				return ipv4, nil
			}
		}
	}
	return nil, nil
}

var discoveryAddress = net.UDPAddr{IP: net.IPv4(239, 192, 48, 84), Port: 24884}
var queryMessage = []byte(`{"type":"query"}`)

func (d *Discovery) sendProbe(config InterfaceConfig) {
	interfaces, err := d.interfaces(config.Name)
	if err != nil {
		logp.Err("failed to get interfaces: ", err)
		return
	}

	var wg sync.WaitGroup
	for _, i := range interfaces {
		ip, err := getIPv4Addr(i)
		if err != nil {
			logp.Err(err.Error())
			continue
		}
		if ip == nil {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.ListenPacket("udp4", net.JoinHostPort(ip.String(), "0"))
			if err != nil {
				logp.Err(err.Error())
				return
			}
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(config.ProbeTimeout))

			if _, err := conn.WriteTo(queryMessage, &discoveryAddress); err != nil {
				logp.Err(err.Error())
				return
			}

			b := make([]byte, 1500)
			for {
				n, _, err := conn.ReadFrom(b)
				if err != nil {
					if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
						logp.Err(err.Error())
					}
					return
				}
				m := make(map[string]interface{})
				err = json.Unmarshal(b[:n], &m)
				if err != nil {
					logp.Err(err.Error())
					continue
				}
				message, _ := messageSchema.Apply(m)
				d.update(config, message)
			}
		}()
	}
	wg.Wait()
}

func (d *Discovery) update(config InterfaceConfig, message common.MapStr) {
	v, err := message.GetValue("agent.id")
	if err != nil {
		logp.Err("failed to update agent without id: " + err.Error())
		return
	}
	agentID, ok := v.(string)
	if len(agentID) == 0 || !ok {
		logp.Err("empty agent?")
		return
	}

	url, err := message.GetValue("url")
	if err != nil || url == nil {
		// This can happen if Jolokia agent is initializing and it still
		// doesn't know its URL
		logp.Info("agent %s without url, ignoring by now", agentID)
		return
	}

	d.Lock()
	defer d.Unlock()
	i, found := d.instances[agentID]
	if !found {
		i = &Instance{Message: message}
		d.instances[agentID] = i
		d.events <- Event{"start", message}
	}
	i.LastSeen = time.Now()
	i.LastInterface = &config
}

func (d *Discovery) checkStopped() {
	d.Lock()
	defer d.Unlock()

	for id, i := range d.instances {
		if time.Since(i.LastSeen) > i.LastInterface.GracePeriod {
			d.events <- Event{"stop", i.Message}
			delete(d.instances, id)
		}
	}
}
