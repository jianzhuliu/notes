/*
公用部分,生成日期 "2020-12-09 17:07:01"
*/
package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	C_time_format_layout = "2006-01-02 15:04:05"
	C_primary_key        = "id" //主键标识
)

//各个表对应创建表操作对象的方法
var TableToObjCreateFunc = make(map[string]func(*sql.DB) Isub)

//db 对象的封装
type IcommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

//所有表处理对象需实现的接口
type Isub interface {
	Db() IcommonDB
	TableName() string                                            //指定表名
	ColumnList() []string                                         //字段列表，字符串切片类型
	Columns() string                                              //对应字段列表
	FieldToColumn() map[string]string                             //结构体字段与表字段对应关系
	Informaton() string                                           //表信息描述
	CurrentTime() string                                          //获取当前时间，主要为兼容 time.Time 类型，需要导入 time 包
	One() (interface{}, error)                                    //查询单个记录
	All() ([]interface{}, error)                                  //查询所有
	Insert(map[string]interface{}, ...interface{}) (int64, error) //插入数据,返回插入后id,后支持布尔型参数，是否不忽略主键字段，默认忽略
	Update(map[string]interface{}) (int64, error)                 //更新数据,返回更新记录数
	Delete() (int64, error)                                       //删除数据,返回受影响行数及删除的行数
	Truncate() error                                              //重置表
	CreateTableSql() (string, error)                              //查看创建表的 sql 语句
	Count() (int64, error)                                        //查询记录数

	Where(string, ...interface{}) Isub //链式操作 where 查询条件
	OrderBy(values ...string) Isub     //链式操作 where 查询条件
	Limit(values ...interface{}) Isub  //链式操作 where 查询条件
	Build() (string, []interface{})    //构建sql语句及参数
	Reset() Isub                       //重置sql语句及参数，以便复用
	Log(string, ...interface{})        //记录日志

	//事务
	Begin() error    //开启事务
	Commit() error   //提交事务
	Rollback() error //事务回滚
	Cancel() Isub    //取消事务模式

	//生成测试数据
	GenTestForUpdate() map[string]interface{}
	GenTestForInsert(int) []map[string]interface{}
}

//time.Time 特殊处理
type TimeNormal struct {
	time.Time
}

//转换为数据库中存储值
//实现 driver.Valuer 接口
func (t TimeNormal) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

//扫描数据库字段
//实现 sql.Scanner 接口
func (t *TimeNormal) Scan(src interface{}) error {
	if value, ok := src.(time.Time); ok {
		*t = TimeNormal{Time: value}
		return nil
	}

	return fmt.Errorf("cann't convert %!v(MISSING) to time.Time", src)
}

//json 格式化
//实现 json.Marshaler 接口
func (t TimeNormal) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(C_time_format_layout)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, C_time_format_layout)
	b = append(b, '"')
	return b, nil
}

//基础对象
type Tbase struct {
	//db 对象
	db       *sql.DB
	conds    []string
	condArgs []interface{}

	sub Isub //子对象，模板模式，用于子对象可访问 base 对象方法

	tx *sql.Tx //开启事务
}

//依赖注入 db 对象
func NewTbase(db *sql.DB) *Tbase {
	return &Tbase{db: db}
}

//获取 db 对象
func (t *Tbase) Db() IcommonDB {
	if t.tx != nil {
		return t.tx
	}
	return t.db
}

//where 查询条件配置
func (t *Tbase) Where(desc string, args ...interface{}) Isub {
	if len(desc) > 0 {
		t.conds = append(t.conds, "where", desc)
		t.condArgs = append(t.condArgs, args...)
	}

	return t.sub
}

//order by  like order by id desc, name asc
func (t *Tbase) OrderBy(values ...string) Isub {
	if len(values) > 0 {
		t.conds = append(t.conds, "order by", strings.Join(values, ","))
	}

	return t.sub
}

//limit 支持 limit 2 or  limit 5,10
func (t *Tbase) Limit(values ...interface{}) Isub {
	if len(values) > 0 {
		t.conds = append(t.conds, "limit ", genBindVars(len(values)))
		t.condArgs = append(t.condArgs, values...)
	}

	return t.sub
}

//重置sql语句及参数，以便复用
func (t *Tbase) Reset() Isub {
	t.conds = nil
	t.condArgs = nil

	return t.sub
}

//构建sql语句及参数
func (t *Tbase) Build() (string, []interface{}) {
	if len(t.conds) > 0 {
		return strings.Join(t.conds, " "), t.condArgs
	}

	return "", nil
}

//获取一条记录
func (t *Tbase) One() (interface{}, error) {
	result, err := t.Limit(1).All()
	if err != nil {
		return nil, err
	}

	//没有数据
	if len(result) == 0 {
		return nil, nil
	}

	return result[0], nil
}

//查看生成表的 sql 语句
func (t *Tbase) CreateTableSql() (string, error) {
	var sql = fmt.Sprintf("show create table %s", t.sub.TableName())
	t.Log("CreateTableSql|sql:%s", sql)
	row := t.Db().QueryRow(sql)
	var db_tblname, createSql string
	if err := row.Scan(&db_tblname, &createSql); err != nil {
		return "", err
	}

	_ = db_tblname

	return createSql, nil
}

//删除记录
func (t *Tbase) Delete() (int64, error) {
	condStr, args := t.Build()
	var sql = fmt.Sprintf("delete from %s %s", t.sub.TableName(), condStr)
	t.Log("Delete|sql:%s,args:%v", sql, args)
	defer t.Reset()
	result, err := t.Db().Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//重置表
func (t *Tbase) Truncate() (err error) {
	var sql = fmt.Sprintf("truncate table %s", t.sub.TableName())
	t.Log("Truncate|sql:%s", sql)
	_, err = t.Db().Exec(sql)

	return
}

//插入，后支持布尔型参数，是否不忽略主键字段，默认忽略
func (t *Tbase) Insert(values map[string]interface{}, params ...interface{}) (int64, error) {
	//简单参数校验
	if len(values) == 0 {
		return 0, fmt.Errorf("values is empty")
	}

	_ = t.Reset()
	defer t.Reset()

	skipPK := true
	if len(params) > 0 {
		if p, ok := params[0].(bool); ok {
			skipPK = p
		}
	}

	//构造sql 及参数
	fieldToColumn := t.sub.FieldToColumn()
	for k, v := range values {
		column, ok := fieldToColumn[k]
		if !ok {
			return 0, fmt.Errorf("key %s is not table column", k)
		}

		if !skipPK || column != C_primary_key {
			t.conds = append(t.conds, column)
			t.condArgs = append(t.condArgs, v)
		}
	}

	if len(t.conds) == 0 {
		return 0, fmt.Errorf("invalid values")
	}

	condStr := strings.Join(t.conds, ",")
	args := t.condArgs
	var sql = fmt.Sprintf("insert into %s(%s) values (%s)", t.sub.TableName(), condStr, genBindVars(len(t.conds)))
	t.Log("Insert|sql:%s,args:%v", sql, args)
	result, err := t.Db().Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

//更新
func (t *Tbase) Update(values map[string]interface{}) (int64, error) {
	//简单参数校验
	if len(values) == 0 {
		return 0, fmt.Errorf("values is empty")
	}

	defer t.Reset()

	columnList := []string{}
	columnValues := []interface{}{}

	//构造sql 及参数
	fieldToColumn := t.sub.FieldToColumn()
	for k, v := range values {
		column, ok := fieldToColumn[k]
		if !ok {
			return 0, fmt.Errorf("key %s is not table column", k)
		}

		columnList = append(columnList, fmt.Sprintf("%s = ?", column))
		columnValues = append(columnValues, v)
	}

	condStr, condArgs := t.Build()
	columns := strings.Join(columnList, ",")

	args := append(columnValues, condArgs...)
	var sql = fmt.Sprintf("update %s set %s %s", t.sub.TableName(), columns, condStr)
	t.Log("Update|sql:%s,args:%v", sql, args)
	result, err := t.Db().Exec(sql, args...)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

//查询记录数
func (t *Tbase) Count() (int64, error) {
	condStr, args := t.Build()
	var sql = fmt.Sprintf("select count(1) as c from %s %s", t.sub.TableName(), condStr)
	t.Log("Count|sql:%s,args:%v", sql, args)
	defer t.Reset()
	row := t.Db().QueryRow(sql, args...)
	var db_count int64
	if err := row.Scan(&db_count); err != nil {
		return 0, err
	}

	return db_count, nil
}

//记录日志
func (t *Tbase) Log(format string, args ...interface{}) {
	logger := log.New(os.Stdout, "[log]", log.LstdFlags|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, args...))
}

/////事务
//开启事务
func (t *Tbase) Begin() (err error) {
	t.Log("transaction begin")

	if t.tx, err = t.db.Begin(); err != nil {
		t.Log("db.Begin()|fail|", err)
	}

	return
}

//提交事务
func (t *Tbase) Commit() (err error) {
	t.Log("transaction commit")

	if err = t.tx.Commit(); err != nil {
		t.Log("tx.Commit()|fail|", err)
	}

	return
}

//事务回滚
func (t *Tbase) Rollback() (err error) {
	t.Log("transaction rollback")

	if err = t.tx.Rollback(); err != nil {
		t.Log("tx.Rollback()|fail|", err)
	}

	return
}

//取消事务模式
func (t *Tbase) Cancel() Isub {
	t.tx = nil
	return t.sub
}

//根据值个数，生成需要的参数位置信息
func genBindVars(num int) string {
	args := make([]string, 0, num)
	for i := 0; i < num; i++ {
		args = append(args, "?")
	}
	return strings.Join(args, ",")
}

//获取调用者所在函数名及文件、行号
func getFuncNameAndFileLine(args ...int) (string, string, int) {
	skip := 1
	if len(args) > 0 {
		skip = args[0]
	}

	//skip, 0:表示本函数, 1:上层调用者
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", 0
	}

	//获取调用者函数名
	fn := runtime.FuncForPC(pc)
	funcName := fn.Name()

	short := file
	c := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			c++
			if c > 1 {
				break
			}
		}
	}
	file = short

	return funcName, file, line
}

//生成测试数据
func (t *Tbase) GenTestForUpdate() map[string]interface{} {
	return nil
}

func (t *Tbase) GenTestForInsert(int) []map[string]interface{} {
	return nil
}
