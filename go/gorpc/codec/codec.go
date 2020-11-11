/*
客户端与服务器公用编解码部分
*/
package codec

import "io"

type Header struct {
	ServiceMethod string //客户端参数，请求服务名及方法，对应服务器是结构体及方法，比如 "Arith.Multiply"
	Seq           uint64 //客户端参数，请求序列号，用于区分不同的请求
	Error         string //服务器参数，出错时的错误信息
}

//解编码接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error         //读取头部
	ReadBody(interface{}) error       //读取body
	Write(*Header, interface{}) error //写入数据
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

//定义了2种数据编解码方式
const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
	NewCodecFuncMap[JsonType] = NewJsonCodec
}
