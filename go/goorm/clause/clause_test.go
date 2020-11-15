package clause

import (
	"reflect"
	"testing"
)

func TestClauseSet(t *testing.T) {
	var clause Clause
	clause.Set(OpTypeInsert, "user", []string{"Id", "Name"})
	sql := clause.sql[OpTypeInsert]
	sqlArgs := clause.sqlArgs[OpTypeInsert]

	//t.Log(sql)
	//t.Log(sqlArgs)

	if sql != "insert into user (Id,Name)" || len(sqlArgs) != 0 {
		t.Fatal("fail to clause set")
	}

	values := make([]interface{}, 2)
	values[0] = []interface{}{1, "name1"}
	values[1] = []interface{}{2, "name2"}
	//t.Log(values)
	clause.Set(OpTypeValues, values...)
	sql = clause.sql[OpTypeValues]
	sqlArgs = clause.sqlArgs[OpTypeValues]
	//t.Log(sql)
	//t.Log(sqlArgs)

	expectSql := "values(?,?),(?,?)"
	if expectSql != sql {
		t.Fatalf("expect %q, but %q got", expectSql, sql)
	}

	expectArgs := []interface{}{1, "name1", 2, "name2"}
	if !reflect.DeepEqual(expectArgs, sqlArgs) {
		t.Fatalf("expect args %v, but %v got", expectArgs, sqlArgs)
	}
}

func TestClauseBuild(t *testing.T) {
	var clause Clause
	clause.Set(OpTypeSelect, "user", []string{"Id", "Name"})
	//t.Log(clause.sql[OpTypeSelect])
	clause.Set(OpTypeWhere, "Id=? and Name=?", 1, "goorm")
	clause.Set(OpTypeOrderBy, "Id asc", "Name desc")
	clause.Set(OpTypeLimit, 1, 9)
	//t.Log(clause.sql[OpTypeLimit])

	sql, sqlArgs := clause.Build()

	expectSql := "select Id,Name from user where Id=? and Name=? order by Id asc,Name desc limit ?,?"
	//t.Log(sql)
	//t.Log(sqlArgs)
	if sql != expectSql {
		t.Fatalf("expect %q, but %q got", expectSql, sql)
	}

	expectArgs := []interface{}{1, "goorm", 1, 9}
	if !reflect.DeepEqual(expectArgs, sqlArgs) {
		t.Fatalf("expect args %v, but %v got", expectArgs, sqlArgs)
	}
}

func TestClauseInsert(t *testing.T) {
	var clause Clause
	clause.Set(OpTypeInsert, "user", []string{"Id", "Name"})

	values := make([]interface{}, 2)
	values[0] = []interface{}{1, "name1"}
	values[1] = []interface{}{2, "name2"}
	clause.Set(OpTypeValues, values...)

	sql, sqlArgs := clause.Build()

	expectSql := "insert into user (Id,Name) values(?,?),(?,?)"
	//t.Log(sql)
	//t.Log(sqlArgs)
	if sql != expectSql {
		t.Fatalf("expect %q, but %q got", expectSql, sql)
	}

	expectArgs := []interface{}{1, "name1", 2, "name2"}
	if !reflect.DeepEqual(expectArgs, sqlArgs) {
		t.Fatalf("expect args %v, but %v got", expectArgs, sqlArgs)
	}
}

func TestClauseCount(t *testing.T) {
	var clause Clause
	clause.Set(OpTypeCount, "user", "Id")
	sql, sqlArgs := clause.Build()
	//t.Log(sql)
	//t.Log(sqlArgs)

	expectSql := "select count(Id) from user"
	if sql != expectSql {
		t.Fatalf("expect sql %s but %s got", expectSql, sql)
	}

	if len(sqlArgs) != 0 {
		t.Fatalf("expect args num 0 but %d got (%v)", len(sqlArgs), sqlArgs)
	}
}

func TestClauseDelete(t *testing.T) {
	var clause Clause
	clause.Set(OpTypeDelete, "user")
	sql, sqlArgs := clause.Build()
	//t.Log(sql)
	//t.Log(sqlArgs)

	expectSql := "delete from user"
	if sql != expectSql {
		t.Fatalf("expect sql %s but %s got", expectSql, sql)
	}

	if len(sqlArgs) != 0 {
		t.Fatalf("expect args num 0 but %d got (%v)", len(sqlArgs), sqlArgs)
	}
}

func TestClauseUpdate(t *testing.T) {
	var clause Clause
	clause.Set(OpTypeUpdate, "user", map[string]interface{}{"Name": "name999", "Id": 123})
	sql, sqlArgs := clause.Build()
	t.Log(sql)
	t.Log(sqlArgs)

	expectSql := "update user set Name=?,Id=?"
	if sql != expectSql {
		t.Fatalf("expect sql %s but %s got", expectSql, sql)
	}

	if len(sqlArgs) != 2 {
		t.Fatalf("expect args num 2 but %d got (%v)", len(sqlArgs), sqlArgs)
	}

	expectArgs := []interface{}{"name999", 123}
	if !reflect.DeepEqual(expectArgs, sqlArgs) {
		t.Fatalf("expect args %v, but %v got", expectArgs, sqlArgs)
	}
}
