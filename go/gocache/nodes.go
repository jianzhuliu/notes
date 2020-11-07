/*

 */
package gocache

type NodePicker interface {
	//根据key， 查找对应的节点
	PickNode(string) (NodeGetter, bool)
}

type NodeGetter interface {
	//节点，根据 group, 查询 对应 key 的数据
	Get(string, string) ([]byte, error)
}
