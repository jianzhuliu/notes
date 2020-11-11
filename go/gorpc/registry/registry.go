/*
注册中心的好处在于，客户端和服务端都只需要感知注册中心的存在，而无需感知对方的存在。

1、服务端启动后，向注册中心发送注册消息，注册中心得知该服务已经启动，处于可用状态。
一般来说，服务端还需要定期向注册中心发送心跳，证明自己还活着。

2、客户端向注册中心询问，当前哪个服务是可用的，注册中心将可用的服务列表返回客户端。
3、客户端根据注册中心得到的服务列表，选择其中一个发起调用。

*/
package registry

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

//注册中心，记录各个服务器数据，包括超时数据
type Registry struct {
	timeout time.Duration
	mu      sync.Mutex
	servers map[string]*ServerItem
}

type ServerItem struct {
	Addr  string    //地址
	start time.Time //注册时间
}

const (
	defaultPath    = "/gorpc/registry"
	defaultTimeout = time.Minute * 5 //默认超时数据
)

func New(timeout time.Duration) *Registry {
	return &Registry{
		timeout: timeout,
		servers: make(map[string]*ServerItem),
	}
}

var DefaultRegister = New(defaultTimeout)

//添加服务
func (r *Registry) addServer(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s := r.servers[addr]
	if s != nil {
		//已经存在， 则更新开始时间，用于保存活动状态
		s.start = time.Now()
		return
	}

	serverItem := &ServerItem{
		Addr:  addr,
		start: time.Now(),
	}

	r.servers[addr] = serverItem
}

//所有可用的服务列表,考虑超时情况
func (r *Registry) aliveServers() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.servers) == 0 {
		return nil
	}

	servers := make([]string, len(r.servers))
	for addr, serverItem := range r.servers {
		if r.timeout == 0 || serverItem.start.Add(r.timeout).After(time.Now()) {
			servers = append(servers, serverItem.Addr)
		} else {
			//主动删除过期的服务
			delete(r.servers, addr)
		}
	}

	sort.Strings(servers)
	return servers
}

///////////////////http
func (r *Registry) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		//获取可用服务列表
		rw.Header().Set("X-gorpc-servers", strings.Join(r.aliveServers(), ","))
	case http.MethodPost:
		//设置服务
		addr := req.Header.Get("X-gorpc-server")
		if addr == "" {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.addServer(addr)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *Registry) HandleHTTP(registryPath string) {
	http.Handle(registryPath, r)
	log.Println("rpc registry path:", registryPath)
}

func HandleHTTP() {
	DefaultRegister.HandleHTTP(defaultPath)
}

////////心跳
//间隔时间，发送心跳，心跳间隔需要在 注册中心过期时间内
func Heartbeat(registry, addr string, duration time.Duration) {
	if duration == 0 {
		duration = defaultTimeout - time.Minute
	}

	var err error
	err = sendHeartbeat(registry, addr)

	go func() {
		t := time.NewTicker(duration)
		for err == nil {
			<-t.C
			err = sendHeartbeat(registry, addr)
		}
	}()
}

func sendHeartbeat(registry, addr string) error {
	log.Printf("%s send heart beat ro registry %s", addr, registry)
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodPost, registry, nil)
	req.Header.Set("X-gorpc-server", addr)

	if _, err := client.Do(req); err != nil {
		log.Println("rpc server:heart beat error:", err)
		return err
	}

	return nil
}
