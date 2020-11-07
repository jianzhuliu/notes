/*
一致性哈希算法将 key 映射到 2^32 的空间中，将这个数字首尾相连，形成一个环
1、计算节点/机器(通常使用节点的名称、编号和 IP 地址)的哈希值，放置在环上
2、计算 key 的哈希值，放置在环上，顺时针寻找到的第一个节点，就是应选取的节点/机器

一致性哈希算法，在新增/删除节点时，只需要重新定位该节点附近的一小部分数据，而不需要重新定位所有的节点

数据倾斜问题
如果服务器的节点过少，容易引起 key 的倾斜
虚拟节点扩充了节点的数量，解决了节点较少的情况下数据容易倾斜的问题。
而且代价非常小，只需要增加一个字典(map)维护真实节点与虚拟节点的映射关系即可。

*/
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type HashObj struct {
	hashFunc   Hash           //哈希算法函数
	replicas   int            //虚拟节点的数量
	hashValues []int          //已经排序了的 hash key列表
	hashToReal map[int]string //虚拟节点与真实节点的对应关系
}

//hash 函数采用依赖注入的方式，方便测试与扩展
func New(replicas int, fn Hash) *HashObj {
	m := &HashObj{
		replicas:   replicas,
		hashFunc:   fn,
		hashToReal: make(map[int]string),
	}

	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}

	return m
}

//添加真实节点
func (obj *HashObj) Add(nodes ...string) {
	for _, nodename := range nodes {
		for i := 0; i < obj.replicas; i++ {
			hashValue := int(obj.hashFunc([]byte(strconv.Itoa(i) + nodename)))
			obj.hashValues = append(obj.hashValues, hashValue)
			obj.hashToReal[hashValue] = nodename
		}
	}

	//排序
	sort.Ints(obj.hashValues)
}

//获取 key 对应hash 环上真实节点名
func (obj *HashObj) Get(key string) string {
	//没有添加真实节点的情况
	if len(obj.hashValues) == 0 {
		return ""
	}

	keyHashValue := int(obj.hashFunc([]byte(key)))
	len := len(obj.hashValues)

	//二分法搜索最大的下标, 结果范围 [0, len]
	idx := sort.Search(len, func(i int) bool {
		return obj.hashValues[i] >= keyHashValue
	})

	//环形结构取模 , 如果 idx == len ，则需要取下标为 0
	realIdx := idx % len
	hashValue := obj.hashValues[realIdx]

	return obj.hashToReal[hashValue]
}
