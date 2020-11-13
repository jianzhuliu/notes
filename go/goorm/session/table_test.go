package session

import (
	"testing"
)

func TestTableModel(t *testing.T) {
	s := TNewSession().Model(&User{})
	table := s.table

	if table == nil || table.Name != "user" {
		t.Fatal("model fail")
	}

	s.Model(&Session{})
	//t.Log(s.table.Name)
	if s.table == nil || s.table.Name != "Session" {
		t.Fatal("fail to change model")
	}
}

func TestTableCreate(t *testing.T) {
	s := TNewSession().Model(&User{})

	table := s.table

	if table == nil || table.Name != "user" {
		t.Fatal("model fail")
	}

	err := s.DropTable()
	if err != nil {
		t.Fatal("drop table fail", err)
	}

	if s.HasTable() {
		t.Fatal("fail to drop table")
	}

	err = s.CreateTable()
	if err != nil {
		t.Fatal("create table fail", err)
	}

	if !s.HasTable() {
		t.Fatal("fail to create table")
	}
}
