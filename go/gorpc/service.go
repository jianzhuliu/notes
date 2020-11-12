package gorpc

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64 //被调用次数
}

//获取函数调用次数
func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

//构造参数的对应类型的实例
func (m *methodType) newArgv() reflect.Value {
	var value reflect.Value
	if m.ArgType.Kind() == reflect.Ptr {
		value = reflect.New(m.ArgType.Elem())
	} else {
		value = reflect.New(m.ArgType).Elem()
	}

	return value
}

//构造返回值对应类型的实例
// Reply 一定是个指针类型，只有指针类型，才可以修改其映射对应的值
func (m *methodType) newReplyv() reflect.Value {
	value := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		value.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		value.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}

	return value
}

type service struct {
	name   string                 //映射的结构体的名称
	typ    reflect.Type           //结构体的类型
	targetObj   reflect.Value          //结构体的实例本身，调用时作为第0个参数
	method map[string]*methodType //所有符合条件的方法列表
}

func newService(targetObj interface{}) *service {
	s := new(service)
	s.targetObj = reflect.ValueOf(targetObj)
	s.name = reflect.Indirect(s.targetObj).Type().Name()
	s.typ = reflect.TypeOf(targetObj)

	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	s.registerMethods()
	return s
}

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mTyp := method.Type
		if mTyp.NumIn() != 3 || mTyp.NumOut() != 1 {
			continue
		}

		if mTyp.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		argType, replyType := mTyp.In(1), mTyp.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}

		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}

		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

//是否为可导出类型，或者内置类型
func isExportedOrBuiltinType(typ reflect.Type) bool {
	return ast.IsExported(typ.Name()) || typ.PkgPath() == ""
}

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.targetObj, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
