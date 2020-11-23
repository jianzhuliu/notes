/*
columns 表结构生成时间 "2020-11-23 17:13:19"
请勿修改，如需新增方法，请另外同包同目录下创建文件处理
*/
package models

import (
	"strings"
	"time"
)

type T_columns struct {
	Status  int8      //状态 1:enable, 0:disable, -1:deleted
	Id      uint32    //
	Name    string    //用户名
	Phone   string    //手机号
	Gender  string    //员工性别
	Info    string    //描述信息
	Created time.Time //创建时间
}

//创建实例
func NewTcolumns() *T_columns {
	return &T_columns{}
}

//所有表字段,按数据库表字段顺序排列
func (t *T_columns) Columns() []string {
	return []string{
		"id", "name", "phone", "gender", "status", "info", "created",
	}
}

//打印表信息
func (t *T_columns) String() string {
	tmp := []string{
		"type T_columns struct{ ",
		"\tStatus\tint8\t//DbOrder:5,KindSize:1,DataType:tinyint,ColumnType:tinyint(1),ColumnComment:状态 1:enable, 0:disable, -1:deleted",
		"\tId\tuint32\t//DbOrder:1,KindSize:4,DataType:int,ColumnType:int unsigned,ColumnComment:",
		"\tName\tstring\t//DbOrder:2,KindSize:16,DataType:varchar,ColumnType:varchar(30),ColumnComment:用户名",
		"\tPhone\tstring\t//DbOrder:3,KindSize:16,DataType:char,ColumnType:char(11),ColumnComment:手机号",
		"\tGender\tstring\t//DbOrder:4,KindSize:16,DataType:enum,ColumnType:enum('male','female','unknow'),ColumnComment:员工性别",
		"\tInfo\tstring\t//DbOrder:6,KindSize:16,DataType:text,ColumnType:text,ColumnComment:描述信息",
		"\tCreated\ttime.Time\t//DbOrder:7,KindSize:24,DataType:timestamp,ColumnType:timestamp,ColumnComment:创建时间",
		"}",
	}

	return strings.Join(tmp, "\r\n")
}

func (t *T_columns) CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
