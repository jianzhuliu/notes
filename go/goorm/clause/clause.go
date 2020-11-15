package clause

import (
	"strings"
)

//定义sql语句支持的类型
type OpType int

const (
	OpTypeBegin OpType = iota
	OpTypeInsert
	OpTypeSelect
	OpTypeUpdate
	OpTypeDelete
	OpTypeCount
	OpTypeWhere
	OpTypeOrderBy
	OpTypeLimit
	OpTypeValues
	OpTypeEnd
)

//不同操作类型，需要的sql语句与对应参数
type Clause struct {
	sql     map[OpType]string
	sqlArgs map[OpType][]interface{}
}

//设置不同操作类型，不同的数据
func (c *Clause) Set(opType OpType, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[OpType]string)
		c.sqlArgs = make(map[OpType][]interface{})
	}

	sql, args := generators[opType](values...)
	c.sql[opType] = sql
	c.sqlArgs[opType] = args
}

//构建 sql语句
func (c *Clause) Build() (string, []interface{}) {
	var sqls []string
	var args []interface{}

	for opType := OpTypeBegin + 1; opType < OpTypeEnd; opType++ {
		//	for _, opType := range opTypes {
		if sql, ok := c.sql[opType]; ok {
			sqls = append(sqls, sql)
			args = append(args, c.sqlArgs[opType]...)
		}
	}

	return strings.Join(sqls, " "), args
}
