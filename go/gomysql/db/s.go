/*
结构体定义
*/
package db

//表字段
type TableColumn struct {
	DbOrder uint8  //数据库中字段排名
	KindStr string //对应 kind 类型

	ColumnName string //字段名
	ColumnType string //字段类型
	DataType   string //字段数据类型
}
