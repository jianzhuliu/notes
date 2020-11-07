/*
负责与外部交互，控制缓存存储和获取

当缓存不存在时，需要一个接口来获取来源数据，暴露给外部，由外部来实现

*/

package gocache

import (
	"fmt"
	"gocache/singleflight"
	"log"
	"sync"
)

//获取来源数据的接口
type Getter interface {
	Get(string) ([]byte, error)
}

//定义一个接口型函数,实现了接口方法，并调用自身
//好处，方便调用时，即可以传入函数作为参数，也可以传入实现了接口的结构体作为参数
type GetterFunc func(string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//缓存的命名空间
type Group struct {
	name      string //唯一标识名
	getter    Getter //获取缓存来源数据接口
	mainCache *cache //缓存策略，并发控制

	//分布式节点选择器
	nodePicker NodePicker

	//加载器，确保查询key，只有一次调用 db
	loader *singleflight.Entry
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

//创建一个实例， getter 不能为nil
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: NewCache(cacheBytes),

		loader: &singleflight.Entry{},
	}

	groups[name] = g
	return g
}

//根据标识名，获取对应的缓存组
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	//先从缓存中获取
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[gocache] hit ", key)
		return v, nil
	}

	//不存在，就通过 getter 来获取
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		//远程获取
		if g.nodePicker != nil {
			//根据节点选择器，选择一个节点
			if nodeGetter, ok := g.nodePicker.PickNode(key); ok {
				if value, err = g.getFromRemoteNodeGetter(nodeGetter, key); err == nil {
					return value, nil
				}

				log.Println("[gocache] failed to get from nodeGetter", err)
			}
		}

		//本地获取
		return g.loadLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}

	return
}

func (g *Group) loadLocally(key string) (ByteView, error) {
	//从 getter 获取数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	//value := ByteView{b: cloneBytes(bytes)}
	value := NewByteView(cloneBytes(bytes))

	//更新缓存
	g.updateCache(key, value)

	log.Println("[gocache] loadLocally ", key)

	return value, nil
}

//远程节点调用获取数据
func (g *Group) getFromRemoteNodeGetter(nodeGetter NodeGetter, key string) (ByteView, error) {
	bytes, err := nodeGetter.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return NewByteView(bytes), nil
}

//更新缓存
func (g *Group) updateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

//分布式节点选择器 注入
func (g *Group) RegisterNodePicker(nodePicker NodePicker) {
	if g.nodePicker != nil {
		panic("RegisterNodePicker called more than once")
	}

	g.nodePicker = nodePicker
}
