/*
缓存值的抽象与封装

只读数据结构 ByteView 用于表示缓存值
可以支持任意的数据类型存储，如字符串，图片
*/
package gocache

type ByteView struct {
	b []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{b: b}
}

//首先实现满足缓存条件的 接口  lru.Value
func (v ByteView) Len() int {
	return len(v.b)
}

//只读结构，提供一个复制的方法，防止外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

//字节切片复制函数
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

//用于打印时调用
func (v ByteView) String() string {
	return string(v.b)
}
