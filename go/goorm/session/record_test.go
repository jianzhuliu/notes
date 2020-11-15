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

	/*
		for _, user := range users {
			t.Log(user)
		}
		//*/

	expect := []User{user1, user2}
	if !reflect.DeepEqual(expect, users) {
		t.Fatal("fail to find")
	}
}

func TestRecordCount(t *testing.T) {
	s := initRecord(t)
	affectedNum, err := s.Count()
	if err != nil {
		t.Fatal("insert fail")
	}

	if affectedNum != 2 {
		t.Fatalf("insert has row affected expect 1 but %d got", affectedNum)
	}
}

func TestRecordDelete(t *testing.T) {
	s := initRecord(t)
	affectedNum, err := s.Delete()
	if err != nil {
		t.Fatal("insert fail")
	}

	if affectedNum != 2 {
		t.Fatalf("insert has row affected expect 1 but %d got", affectedNum)
	}
}

func TestRecordUpdate(t *testing.T) {
	s := initRecord(t)
	kv := map[string]interface{}{
		"Name": "name999",
	}
	affectedNum, err := s.Update(kv)
	if err != nil {
		t.Fatal("insert fail")
	}

	if affectedNum != 2 {
		t.Fatalf("insert has row affected expect 1 but %d got", affectedNum)
	}

	var users []User
	err = s.Find(&users)
	if err != nil {
		t.Fatal("find fail", err)
	}

	if len(users) != 2 {
		t.Fatalf("find users expect 2 but %d got", len(users))
	}

	for _, user := range users {
		if user.Name != "name999" {
			t.Fatalf("record Id=%d should get name name999, but %s got", user.Id, user.Name)
		}
	}

}

func TestRecordChain(t *testing.T) {
	s := TNewSession()
	s.OrderBy("Id asc", "Name desc").Limit(0, 5).Where("Id=?", 1)
	sql, sqlArgs := s.clause.Build()

	expectSql := "where Id=? order by Id asc,Name desc limit ?,?"
	if sql != expectSql {
		t.Fatalf("expect %q, but %q got", expectSql, sql)
	}

	expectArgs := []interface{}{1, 0, 5}
	if !reflect.DeepEqual(expectArgs, sqlArgs) {
		t.Fatalf("expect args %v, but %v got", expectArgs, sqlArgs)
	}
}

func TestRecordFirst(t *testing.T) {
	s := initRecord(t)
	user := &User{}
	err := s.First(user)
	if err != nil {
		t.Fatal("first fail", err)
	}

	//t.Log(user)
	if user.Id != 1 || user.Name != "name1" {
		t.Fatalf("expect Id=%d, Name=%s, but get Id=%d, Name=%s", 1, "name1", user.Id, user.Name)
	}

	err = s.Where("Id=?", 2).First(user)
	if err != nil {
		t.Fatal("first fail", err)
	}
	//t.Log(user)
	if user.Name != "name2" {
		t.Fatalf("expect Name=name2 but %s got", user.Name)
	}
}
