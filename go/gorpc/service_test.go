package gorpc

import (
	"reflect"
	"testing"
)

type Foo int

type Args struct {
	Num1 int
	Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func TestNewService(t *testing.T) {
	var foo Foo
	s := newService(&foo)

	if s.name != "Foo" {
		t.Fatalf("s.name is not right ,expect foo, but %s got", s.name)
	}

	if len(s.method) != 1 {
		t.Fatalf("wrong service Method, expect 1, but %d got", len(s.method))
	}

	if mType := s.method["Sum"]; mType == nil {
		t.Fatal("wrong Metho, Sum shouldn't nil")
	}
}

func TestMethodType_Call(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	mType := s.method["Sum"]

	argv := mType.newArgv()
	replyv := mType.newReplyv()
	argv.Set(reflect.ValueOf(Args{Num1: 1, Num2: 3}))

	err := s.call(mType, argv, replyv)

	if err != nil {
		t.Fatalf("s.call should be nil , but %v got", err)
	}

	if *replyv.Interface().(*int) != 4 {
		t.Fatalf("the result not right ,expect 4, but %d got", *replyv.Interface().(*int))
	}

	if mType.NumCalls() != 1 {
		t.Fatalf("the num calls not right ,expect 1, but %d got", mType.NumCalls())
	}
}
