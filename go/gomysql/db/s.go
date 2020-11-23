/*
结构体定义
*/
package db

//表字段
type TableColumn struct {
	DbOrder  uint8  //数据库中字段排名
	KindSize uint8  //对应 kind size大小
	KindStr  string //对应 kind 类型

	ColumnName string //字段名
	ColumnType string //字段类型
	DataType   string //字段数据类型
}

type TableColumnSice []TableColumn

//排序
func (d TableColumnSice) Len() int      { return len(d) }
func (d TableColumnSice) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d TableColumnSice) Less(i, j int) bool {
	if d[i].KindSize == d[j].KindSize {
		return d[i].DbOrder < d[j].DbOrder
	}
	return d[i].KindSize < d[j].KindSize
}
