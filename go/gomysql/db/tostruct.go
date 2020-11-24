package db

import (
	"bytes"
	"sort"
	"strings"
	"time"

	"gomysql/conf"
	"text/template"
)

/////////////////////////////////表结构转换为 struct

const tmplTableToStruct = `/*
{{.tblname}} 表结构生成时间 "{{.created}}"
请勿修改，如需新增方法，请另外同包同目录下创建文件处理
*/
package models

import (
	"strings"
	"time"
	"database/sql"
	"fmt"
)

//表结构体
type T_{{.tblname}} struct{
{{- range $column := .tableColumns}}
	{{$column.ColumnName|Title}}		{{$column.KindStr}}		//{{$column.ColumnComment}}
{{- end }}
}

//打印结构体数据
func (t T_{{.tblname}}) String() string{
	formats := []string{}
	args := []interface{}{}
	{{range $column := .tableColumns}}
	formats = append(formats,"{{$column.ColumnName|Title}}:%v")
	
	{{- if eq $column.KindStr "time.Time" }}
	args = append(args, t.{{$column.ColumnName|Title}}.Format(C_time_format_layout))
	{{- else}}
	args = append(args, t.{{$column.ColumnName|Title}})
	{{- end}}
	{{end }}
	return fmt.Sprintf(strings.Join(formats, ","), args...)
}

//表操作对象
type Tobj_{{.tblname}} struct {
	*Tbase
}

//创建实例
func NewTobj_{{.tblname}}(db *sql.DB) *Tobj_{{.tblname}} {
	tbase := NewTbase(db)
	t := &Tobj_columns{
		Tbase: tbase,
	}
	
	tbase.sub = t
	return t
}

//对象转换
func (t *Tobj_{{.tblname}}) Interface(value interface{}) (data T_{{.tblname}} , ok bool){
	data, ok = value.(T_{{.tblname}})
	return
}

//设置表名
func (t *Tobj_{{.tblname}}) TableName() string {
	return "{{.tblname}}"
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_{{.tblname}}) Columns() []string{
	return []string{
	{{- range $field := .fields}}
	"{{$field}}",
	{{- end }}
	}
}

//所有表字段,按数据库表字段顺序排列
func (t *Tobj_{{.tblname}}) Fields() string {
	return "{{Join .fields "," }}"
}

//打印表信息
func (t *Tobj_{{.tblname}}) Informaton() string{
	tmp := []string{
	"type T_{{.tblname}} struct{ ",
	{{- range $column := .tableColumns}}
	"\t{{$column.ColumnName|Title}}\t{{$column.KindStr}}\t//DbOrder:{{$column.DbOrder}},KindSize:{{$column.KindSize}},DataType:{{$column.DataType}},ColumnType:{{$column.ColumnType}},ColumnComment:{{$column.ColumnComment}}",
	{{- end}}
	"}",
	}

	return strings.Join(tmp, "\r\n")
}

//获取当前时间
func (t *Tobj_{{.tblname}}) CurrentTime() string{
	return time.Now().Format(C_time_format_layout)
}

//获取满足条件下，所有记录
func (t *Tobj_{{.tblname}}) All() ([]interface{},error) {
	condStr,args := t.Build()
	var sql = fmt.Sprintf("select %s from %s %s", t.Fields(), t.TableName(), condStr)
	t.Log("sql:%s,args:%v",sql, args)
	defer t.Reset()
	rows, err := t.db.Query(sql, args...)
	
	if err != nil {
		return nil, err
	}
	
	data := []interface{}{}
	
	for rows.Next() {
		{{- range $column := .tableColumns}}
		var {{$column.ColumnName}} {{$column.KindStr}}
		{{- end }}

		if err := rows.Scan(
		{{- range $field := .fields}}
		&{{$field}},
		{{- end }}
		); err != nil {
			return nil, err
		}
		
		result := T_columns{}
		{{- range $column := .tableColumns}}
		result.{{$column.ColumnName|Title}} = {{$column.ColumnName}}
		{{- end }}
		
		data = append(data, result)
	}
	
	return data, nil
}

`

var tableToStructT = template.Must(template.New("table_to_struct").
	Funcs(template.FuncMap{
		"Title": strings.Title,
		"Join":  strings.Join,
	}).
	Parse(tmplTableToStruct))

func ToStruct(tblname string, tableColumns TableColumnSice) (string, error) {
	var buf bytes.Buffer

	fields := make([]string, 0, len(tableColumns))
	for _, column := range tableColumns {
		fields = append(fields, column.ColumnName)
	}

	sort.Sort(tableColumns)
	data := map[string]interface{}{
		"tblname":          tblname,
		"tableColumns":     tableColumns,
		"fields":           fields,
		"created":          time.Now().Format(conf.C_time_layout),
		"timeFormatLayout": conf.C_time_layout,
	}
	err := tableToStructT.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
