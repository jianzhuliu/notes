/*
公用部分,生成日期 "2020-11-24 18:23:36"
*/
package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	C_time_format_layout = "2006-01-02 15:04:05"
)

//所有表处理对象需实现的接口
type Isub interface {
	TableName() string           //指定表名
	Columns() []string           //字段列表，字符串切片类型
	Fields() string              //对应字段列表
	Informaton() string          //表信息描述
	CurrentTime() string         //获取当前时间，主要为兼容 time.Time 类型，需要导入 time 包
	One() (interface{}, error)   //查询单个记录
	All() ([]interface{}, error) //查询所有

	Where(desc string, args ...interface{}) Isub //链式操作 where 查询条件
	OrderBy(values ...string) Isub               //链式操作 where 查询条件
	Limit(values ...interface{}) Isub            //链式操作 where 查询条件
	Build() (string, []interface{})              //构建sql语句及参数
	Reset() Isub                                 //重置sql语句及参数，以便复用
	Log(string, ...interface{})                  //记录日志
}

//基础对象
type Tbase struct {
	//db 对象
	db       *sql.DB
	conds    []string
	condArgs []interface{}

	sub Isub //子对象，模板模式，用于子对象可访问 base 对象方法
}

//依赖注入 db 对象
func NewTbase(db *sql.DB) *Tbase {
	return &Tbase{db: db}
}

//获取 db 对象
func (t *Tbase) GetDb() *sql.DB {
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

	return result[0], nil
}

//记录日志
func (t *Tbase) Log(format string, args ...interface{}) {
	logger := log.New(os.Stdout, "[log]", log.LstdFlags|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, args...))
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
