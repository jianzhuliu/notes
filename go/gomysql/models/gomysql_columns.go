/*
columns 表结构生成时间 "2020-11-25 12:29:22"
请勿修改，如需新增方法，请另外同包同目录下创建文件处理
*/
package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

//表结构体
type T_columns struct {
	Status  int8      //状态 1:enable, 0:disable, -1:deleted
	Id      uint      //
	Name    string    //用户名
	Phone   string    //手机号
	Gender  string    //员工性别
	Info    string    //描述信息
	Created time.Time //创建时间
}

//打印结构体数据
func (t T_columns) String() string {
	formats := []string{}
	args := []interface{}{}

	formats = append(formats, "Status:%v")
	args = append(args, t.Status)

	formats = append(formats, "Id:%v")
	args = append(args, t.Id)

	formats = append(formats, "Name:%v")
	args = append(args, t.Name)

	formats = append(formats, "Phone:%v")
	args = append(args, t.Phone)

	formats = append(formats, "Gender:%v")
	args = append(args, t.Gender)

	formats = append(formats, "Info:%v")
	args = append(args, t.Info)

	formats = append(formats, "Created:%v")
	args = append(args, t.Created.Format(C_time_format_layout))

	return fmt.Sprintf(strings.Join(formats, ","), args...)
}

//表操作对象
type Tobj_columns struct {
	*Tbase
}

//创建实例
func NewTobj_columns(db *sql.DB) *Tobj_columns {
	tbase := NewTbase(db)
	t := &Tobj_columns{
		Tbase: tbase,
	}

	tbase.sub = t
	return t
}

//对象转换
func (t *Tobj_columns) Interface(value interface{}) (data T_columns, ok bool) {
	data, ok = value.(T_columns)
	return
}

//设置表名
func (t *Tobj_columns) TableName() string {
	return "columns"
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_columns) ColumnList() []string {
	return []string{
		"id",
		"name",
		"phone",
		"gender",
		"status",
		"info",
		"created",
	}
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_columns) Columns() string {
	return "id,name,phone,gender,status,info,created"
}

//结构体字段与表字段对应关系
func (t *Tobj_columns) FieldToColumn() map[string]string {
	return map[string]string{
		"Status":  "status",
		"Id":      "id",
		"Name":    "name",
		"Phone":   "phone",
		"Gender":  "gender",
		"Info":    "info",
		"Created": "created",
	}
}

//打印表信息
func (t *Tobj_columns) Informaton() string {
	tmp := []string{
		"type T_columns struct{ ",
		"\tStatus\tint8\t//DbOrder:5,KindSize:1,DataType:tinyint,ColumnType:tinyint(1),ColumnComment:状态 1:enable, 0:disable, -1:deleted",
		"\tId\tuint\t//DbOrder:1,KindSize:8,DataType:int,ColumnType:int unsigned,ColumnComment:",
		"\tName\tstring\t//DbOrder:2,KindSize:16,DataType:varchar,ColumnType:varchar(30),ColumnComment:用户名",
		"\tPhone\tstring\t//DbOrder:3,KindSize:16,DataType:char,ColumnType:char(11),ColumnComment:手机号",
		"\tGender\tstring\t//DbOrder:4,KindSize:16,DataType:enum,ColumnType:enum('male','female','unknow'),ColumnComment:员工性别",
		"\tInfo\tstring\t//DbOrder:6,KindSize:16,DataType:text,ColumnType:text,ColumnComment:描述信息",
		"\tCreated\ttime.Time\t//DbOrder:7,KindSize:24,DataType:timestamp,ColumnType:timestamp,ColumnComment:创建时间",
		"}",
	}

	return strings.Join(tmp, "\r\n")
}

//获取当前时间
func (t *Tobj_columns) CurrentTime() string {
	return time.Now().Format(C_time_format_layout)
}

//获取满足条件下，所有记录
func (t *Tobj_columns) All() ([]interface{}, error) {
	condStr, args := t.Build()
	var sql = fmt.Sprintf("select %s from %s %s", t.Columns(), t.TableName(), condStr)
	t.Log("All|sql:%s,args:%v", sql, args)
	defer t.Reset()
	rows, err := t.Db().Query(sql, args...)

	if err != nil {
		return nil, err
	}

	data := []interface{}{}

	for rows.Next() {
		var status int8
		var id uint
		var name string
		var phone string
		var gender string
		var info string
		var created time.Time

		if err := rows.Scan(
			&id,
			&name,
			&phone,
			&gender,
			&status,
			&info,
			&created,
		); err != nil {
			return nil, err
		}

		result := T_columns{}
		result.Status = status
		result.Id = id
		result.Name = name
		result.Phone = phone
		result.Gender = gender
		result.Info = info
		result.Created = created

		data = append(data, result)
	}

	return data, nil
}
