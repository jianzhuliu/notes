package session

import (
	"goorm/log"
	"testing"
)

type Account struct {
	Id     int `goorm:"primary key"`
	Passwd string
}

func (a *Account) BeforeInsert(s *Session) {
	log.Info("before insert", a)
	a.Id += 100
}

func (a *Account) AfterQuery(s *Session) {
	log.Info("after query", a)
	a.Passwd = "*******"
}

func TestHooks(t *testing.T) {
	s := TNewSession().Model(&Account{})
	_ = s.DropTable()
	err := s.CreateTable()
	if err != nil {
		t.Fatal("create table fail", err)
	}

	account := &Account{Id: 1, Passwd: "asfasd"}
	affectedNum, err := s.Insert(account)
	if err != nil {
		t.Fatal("insert fail", err)
	}

	if affectedNum != 1 {
		t.Fatalf("expect affected num 1 ,but %d got", affectedNum)
	}

	result := &Account{}
	err = s.First(result)
	if err != nil {
		t.Fatal("query first fail", err)
	}

	t.Log(result)
	if result.Id != 101|| result.Passwd != "*******" {
		t.Fatal("fail to call hooks after query", result)
	}
}
