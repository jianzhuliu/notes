package xclient

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type RegistryDiscovery struct {
	*MultiServersDiscovery

	registry   string        //注册中心地址
	timeout    time.Duration //服务列表过期时间，默认 10s
	lastUpdate time.Time     //最后从注册中心更新服务列表时间
}

const defaultUpdateTimeout = time.Second * 10

func NewRegistryDiscovery(registry string, timeout time.Duration) *RegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	return &RegistryDiscovery{
		MultiServersDiscovery: NewMultiServersDiscovery(),
		timeout:               timeout,
		registry:              registry,
	}
}

func (d *RegistryDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	//还没有达到刷新时间
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}

	log.Println("rpc registry: refresh servers from registry", d.registry)
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return nil
	}

	servers := strings.Split(resp.Header.Get("X-gorpc-servers"), ",")
	if len(servers) == 0 {
		log.Println("rpc register refresh get no alive servers")
		return fmt.Errorf("rpc register refresh get no alive servers")
	}

	d.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		if strings.TrimSpace(server) != "" {
			d.servers = append(d.servers, strings.TrimSpace(server))
		}
	}
	d.lastUpdate = time.Now()

	return nil
}

func (d *RegistryDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *RegistryDiscovery) Get(mode SelectMode) (string, error) {
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

func (d *RegistryDiscovery) GetAll() ([]string, error) {
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}
