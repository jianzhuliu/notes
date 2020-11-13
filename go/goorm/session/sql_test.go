package session

import (
	"testing"
)

func TestSql(t *testing.T) {
	s := TNewSession()
	_, _ = s.Sql("drop table if exists user").Exec()
	_, err := s.Sql("create table if not exists user(id int(10) unsigned not null primary key auto_increment,name varchar(30) not null default '')engine=innodb default charset=utf8mb4 ").Exec()
	if err != nil {
		t.Fatal(err)
	}

	result, err := s.Sql("insert into user(name) values(?),(?)", "name1", "name2").Exec()
	if err != nil {
		t.Fatal(err)
	}

	if affectedNum, err := result.RowsAffected(); err != nil || affectedNum != 2 {
		t.Fatalf("expect 2 but %d got", affectedNum)
	}

	row := s.Sql("select count(1) as c from user limit 1").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil {
		t.Fatal(err)
	}
	//t.Log(count)
	if count != 2 {
		t.Fatalf("expect 2 but %d got", count)
	}

	rows, err := s.Sql("select id,name from user limit 1").QueryRows()
	if err != nil {
		t.Fatal(err)
	}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			t.Fatal(err)
		}
		//t.Log(id,name)

		if id != 1 || name != "name1" {
			t.Fatalf("expect 1, name1 but got %v,%v", id, name)
		}
	}
}
