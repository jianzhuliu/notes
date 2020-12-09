/*
columns 表结构生成时间 "2020-12-09 17:07:01"
请勿修改，如需新增方法，请另外同包同目录下创建文件处理
*/
package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

//注册表对应创建操作对象的方法
func init() {
	TableToObjCreateFunc["columns"] = func(db *sql.DB) Isub {
		return NewTobj_columns(db)
	}
}

//表结构体
type T_columns struct {
	Created TimeNormal `json:"created,omitempty"` //创建时间
	Status  int8       `json:"status,omitempty"`  //状态 1:enable, 0:disable, -1:deleted
	Id      uint       `json:"-"`                 //
	Gender  string     `json:"gender,omitempty"`  //员工性别
	Info    string     `json:"info,omitempty"`    //描述信息
	Name    string     `json:"name,omitempty"`    //用户名
	Phone   string     `json:"phone,omitempty"`   //手机号
}

//打印结构体数据
func (t T_columns) String() string {
	formats := []string{}
	args := []interface{}{}

	formats = append(formats, "Created:%v")
	args = append(args, t.Created)

	formats = append(formats, "Status:%v")
	args = append(args, t.Status)

	formats = append(formats, "Id:%v")
	args = append(args, t.Id)

	formats = append(formats, "Gender:%v")
	args = append(args, t.Gender)

	formats = append(formats, "Info:%v")
	args = append(args, t.Info)

	formats = append(formats, "Name:%v")
	args = append(args, t.Name)

	formats = append(formats, "Phone:%v")
	args = append(args, t.Phone)

	return fmt.Sprintf(strings.Join(formats, ","), args...)
}

//生成键值对
func (t T_columns) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["Created"] = t.Created
	result["Status"] = t.Status
	result["Id"] = t.Id
	result["Gender"] = t.Gender
	result["Info"] = t.Info
	result["Name"] = t.Name
	result["Phone"] = t.Phone

	return result
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
		"created",
		"gender",
		"id",
		"info",
		"name",
		"phone",
		"status",
	}
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_columns) Columns() string {
	return "created,gender,id,info,name,phone,status"
}

//结构体字段与表字段对应关系
func (t *Tobj_columns) FieldToColumn() map[string]string {
	return map[string]string{
		"Created": "created",
		"Status":  "status",
		"Id":      "id",
		"Gender":  "gender",
		"Info":    "info",
		"Name":    "name",
		"Phone":   "phone",
	}
}

//打印表信息
func (t *Tobj_columns) Informaton() string {
	tmp := []string{
		"type T_columns struct{ ",
		"\tCreated\tTimeNormal\t//DbOrder:1,KindSize:0,DataType:timestamp,ColumnType:timestamp,ColumnComment:创建时间",
		"\tStatus\tint8\t//DbOrder:7,KindSize:1,DataType:tinyint,ColumnType:tinyint(1),ColumnComment:状态 1:enable, 0:disable, -1:deleted",
		"\tId\tuint\t//DbOrder:3,KindSize:8,DataType:int,ColumnType:int unsigned,ColumnComment:",
		"\tGender\tstring\t//DbOrder:2,KindSize:16,DataType:enum,ColumnType:enum('male','female','unknow'),ColumnComment:员工性别",
		"\tInfo\tstring\t//DbOrder:4,KindSize:16,DataType:text,ColumnType:text,ColumnComment:描述信息",
		"\tName\tstring\t//DbOrder:5,KindSize:16,DataType:varchar,ColumnType:varchar(30),ColumnComment:用户名",
		"\tPhone\tstring\t//DbOrder:6,KindSize:16,DataType:char,ColumnType:char(11),ColumnComment:手机号",
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
		var created TimeNormal
		var status int8
		var id uint
		var gender string
		var info string
		var name string
		var phone string

		if err := rows.Scan(
			&created,
			&gender,
			&id,
			&info,
			&name,
			&phone,
			&status,
		); err != nil {
			return nil, err
		}

		result := T_columns{}
		result.Created = created
		result.Status = status
		result.Id = id
		result.Gender = gender
		result.Info = info
		result.Name = name
		result.Phone = phone

		data = append(data, result)
	}

	return data, nil
}
