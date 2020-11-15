package session

import (
	"reflect"
	"testing"
)

var (
	user1 = User{Id: 1, Name: "name1"}
	user2 = User{Id: 2, Name: "name2"}
)

func initRecord(t *testing.T) *Session {
	t.Helper() //如果出错，显示调用函数的信息

	s := TNewSession().Model(&User{})
	_ = s.DropTable()
	err := s.CreateTable()
	if err != nil {
		t.Fatal("create table fail", err)
	}

	if !s.HasTable() {
		t.Fatal("table is not exist")
	}

	insertNum, err := s.Insert(user1, user2)
	if err != nil {
		t.Fatal("insert fail", err)
	}

	if insertNum != 2 {
		t.Fatalf("insert has row affected expect 2 but %d got", insertNum)
	}

	return s
}

func TestRecordInsert(t *testing.T) {
	s := initRecord(t)
	user3 := &User{Id: 3, Name: "name3"}
	insertNum, err := s.Insert(user3)
	if err != nil {
		t.Fatal("insert fail")
	}

	if insertNum != 1 {
		t.Fatalf("insert has row affected expect 1 but %d got", insertNum)
	}

}

func TestRecordFind(t *testing.T) {
	s := initRecord(t)
	var users []User
	err := s.Find(&users)
	if err != nil {
		t.Fatal("find fail", err)
	}

	if len(users) != 2 {
		t.Fatalf("find users expect 2 but %d got", len(users))
	}

	for _, user := range users {
		t.Log(user)
	}

	expect := []User{user1, user2}
	if !reflect.DeepEqual(expect, users) {
		t.Fatal("fail to find")
	}
}
