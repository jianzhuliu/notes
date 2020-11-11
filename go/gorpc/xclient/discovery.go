/*
假设有多个服务实例，每个实例提供相同的功能，为了提高整个系统的吞吐量，每个实例部署在不同的机器上
客户端可以选择任意一个实例进行调用，获取想要的结果

对于 RPC 框架来说, 可以有如下负载均衡策略
1、随机选择策略 - 从服务列表中随机选择一个
2、轮询算法(Round Robin) - 依次调度不同的服务器，每次调度执行 i = (i + 1) mode n。
3、加权轮询(Weight Round Robin) - 在轮询算法的基础上，为每个服务实例设置一个权重，高性能的机器赋予更高的权重，
也可以根据服务实例的当前的负载情况做动态的调整，例如考虑最近5分钟部署服务器的 CPU、内存消耗情况。

4、哈希/一致性哈希策略 - 依据请求的某些特征，计算一个 hash 值，根据 hash 值将请求发送到对应的机器。
一致性 hash 还可以解决服务实例动态添加情况下，调度抖动的问题

*/

package xclient

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota //随机选择
	RoundRobinSelect                   //轮询算法
)

//服务发现最基本的接口
type Discovery interface {
	Refresh() error                 //从注册中心更新服务列表
	Update([]string) error          //手动更新服务列表
	Get(SelectMode) (string, error) //根据负载均衡策略，选择一个服务实例
	GetAll() ([]string, error)      //返回所有的服务实例
}

//实现一个不需要注册中心，服务列表由手工维护的服务发现的结构体
type MultiServersDiscovery struct {
	r  *rand.Rand
	mu sync.RWMutex //读多写少情况下，读写锁

	servers []string //服务列表
	index   int      //轮询算法，记录下标
}

func NewMultiServersDiscovery(servers ...string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	//避免每次从 0 开始，初始化时随机设定一个值
	d.index = d.r.Intn(math.MaxInt32 - 1)

	return d
}

//实现接口
var _ Discovery = (*MultiServersDiscovery)(nil)

//手工维护，不需要刷新
func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.servers = servers
	return nil
}

func (d *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	size := len(d.servers)
	if size == 0 {
		return "", errors.New("rpc discovery: no availabe servers")
	}

	index := 0

	switch mode {
	case RandomSelect:
		index = d.r.Intn(size)
	case RoundRobinSelect:
		index = d.index % size
		d.index = (index + 1) % size
	default:
		return "", fmt.Errorf("rpc discovery: not supported select mode %v", mode)
	}

	return d.servers[index], nil
}

func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	tmp := make([]string, len(d.servers), len(d.servers))
	copy(tmp, d.servers)
	return tmp, nil
}
