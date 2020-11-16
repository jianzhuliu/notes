package goorm

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"goorm/conf"
	"goorm/session"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := NewEngine(conf.GetDsnByDriver(conf.DriverMysql))

	if engine == nil || err != nil {
		t.Fatal("数据库连接失败", err)
	}

	defer engine.Close()
}

//初始化引擎
func initEngine(t *testing.T) *Engine {
	t.Helper()

	engine, err := NewEngine(conf.GetDsnByDriver(conf.DriverMysql))
	if err != nil {
		t.Fatal("fail to open engine", err)
	}

	return engine
}

//回滚测试
func TestTransactionRollBack(t *testing.T) {
	e := initEngine(t)
	defer e.Close()

	s := e.NewSession().Model(&conf.User{})
	_ = s.DropTable()
	_ = s.CreateTable()

	if !s.HasTable() {
		t.Fatal("create table fail")
	}

	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		s.Model(&conf.User{})
		_, _ = s.Insert(conf.User1, conf.User2)
		err = errors.New("rollback test")
		return
	})

	if err == nil {
		t.Fatal("fail to transaction")
	}

	affectedNum, err := s.Count()
	if err != nil {
		t.Fatal(err)
	}

	if affectedNum != 0 {
		t.Fatalf("expect affected num 0 ,but %d got", affectedNum)
	}
}

//提交测试
func TestTransactionCommit(t *testing.T) {
	e := initEngine(t)
	defer e.Close()

	s := e.NewSession().Model(&conf.User{})
	_ = s.DropTable()
	_ = s.CreateTable()

	if !s.HasTable() {
		t.Fatal("create table fail")
	}

	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		s.Model(&conf.User{})

		_, err = s.Insert(conf.User1, conf.User2)
		if err != nil {
			return
		}

		return
	})

	if err != nil {
		t.Fatal(err)
	}

	affectedNum, err := s.Count()
	if err != nil {
		t.Fatal(err)
	}

	if affectedNum != 2 {
		t.Fatalf("expect affected num 2 ,but %d got", affectedNum)
	}
}
