/*
segments 表结构生成时间 "2020-12-09 17:06:37"
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
	TableToObjCreateFunc["segments"] = func(db *sql.DB) Isub {
		return NewTobj_segments(db)
	}
}

//表结构体
type T_segments struct {
	Create_time int64  `json:"create_time,omitempty"` //
	Max_id      int64  `json:"max_id,omitempty"`      //
	Update_time int64  `json:"update_time,omitempty"` //
	Step        int    `json:"step,omitempty"`        //
	Biz_tag     string `json:"biz_tag,omitempty"`     //
	Remark      string `json:"remark,omitempty"`      //
}

//打印结构体数据
func (t T_segments) String() string {
	formats := []string{}
	args := []interface{}{}

	formats = append(formats, "Create_time:%v")
	args = append(args, t.Create_time)

	formats = append(formats, "Max_id:%v")
	args = append(args, t.Max_id)

	formats = append(formats, "Update_time:%v")
	args = append(args, t.Update_time)

	formats = append(formats, "Step:%v")
	args = append(args, t.Step)

	formats = append(formats, "Biz_tag:%v")
	args = append(args, t.Biz_tag)

	formats = append(formats, "Remark:%v")
	args = append(args, t.Remark)

	return fmt.Sprintf(strings.Join(formats, ","), args...)
}

//生成键值对
func (t T_segments) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["Create_time"] = t.Create_time
	result["Max_id"] = t.Max_id
	result["Update_time"] = t.Update_time
	result["Step"] = t.Step
	result["Biz_tag"] = t.Biz_tag
	result["Remark"] = t.Remark

	return result
}

//表操作对象
type Tobj_segments struct {
	*Tbase
}

//创建实例
func NewTobj_segments(db *sql.DB) *Tobj_segments {
	tbase := NewTbase(db)
	t := &Tobj_segments{
		Tbase: tbase,
	}

	tbase.sub = t
	return t
}

//对象转换
func (t *Tobj_segments) Interface(value interface{}) (data T_segments, ok bool) {
	data, ok = value.(T_segments)
	return
}

//设置表名
func (t *Tobj_segments) TableName() string {
	return "segments"
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_segments) ColumnList() []string {
	return []string{
		"biz_tag",
		"create_time",
		"max_id",
		"remark",
		"step",
		"update_time",
	}
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_segments) Columns() string {
	return "biz_tag,create_time,max_id,remark,step,update_time"
}

//结构体字段与表字段对应关系
func (t *Tobj_segments) FieldToColumn() map[string]string {
	return map[string]string{
		"Create_time": "create_time",
		"Max_id":      "max_id",
		"Update_time": "update_time",
		"Step":        "step",
		"Biz_tag":     "biz_tag",
		"Remark":      "remark",
	}
}

//打印表信息
func (t *Tobj_segments) Informaton() string {
	tmp := []string{
		"type T_segments struct{ ",
		"\tCreate_time\tint64\t//DbOrder:2,KindSize:0,DataType:bigint,ColumnType:bigint,ColumnComment:",
		"\tMax_id\tint64\t//DbOrder:3,KindSize:0,DataType:bigint,ColumnType:bigint,ColumnComment:",
		"\tUpdate_time\tint64\t//DbOrder:6,KindSize:0,DataType:bigint,ColumnType:bigint,ColumnComment:",
		"\tStep\tint\t//DbOrder:5,KindSize:8,DataType:int,ColumnType:int,ColumnComment:",
		"\tBiz_tag\tstring\t//DbOrder:1,KindSize:16,DataType:varchar,ColumnType:varchar(128),ColumnComment:",
		"\tRemark\tstring\t//DbOrder:4,KindSize:16,DataType:varchar,ColumnType:varchar(200),ColumnComment:",
		"}",
	}

	return strings.Join(tmp, "\r\n")
}

//获取当前时间
func (t *Tobj_segments) CurrentTime() string {
	return time.Now().Format(C_time_format_layout)
}

//获取满足条件下，所有记录
func (t *Tobj_segments) All() ([]interface{}, error) {
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
		var create_time int64
		var max_id int64
		var update_time int64
		var step int
		var biz_tag string
		var remark string

		if err := rows.Scan(
			&biz_tag,
			&create_time,
			&max_id,
			&remark,
			&step,
			&update_time,
		); err != nil {
			return nil, err
		}

		result := T_segments{}
		result.Create_time = create_time
		result.Max_id = max_id
		result.Update_time = update_time
		result.Step = step
		result.Biz_tag = biz_tag
		result.Remark = remark

		data = append(data, result)
	}

	return data, nil
}
