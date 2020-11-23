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
	"time"
	"strings"
)

type T_{{.tblname}} struct{
{{- range $column := .tableColumns}}
	{{$column.ColumnName|title}}		{{$column.KindStr}}		//{{$column.ColumnComment}}
{{- end }}
}

//创建实例
func NewT{{.tblname}}() *T_{{.tblname}} {
	return &T_{{.tblname}}{}
}

//所有表字段,按数据库表字段顺序排列
func (t *T_{{.tblname}}) Columns() []string{
	return []string{
	{{.fieldsStr}},
	}
}

//打印表信息
func (t *T_{{.tblname}}) String() string{
	tmp := []string{
	"type T_{{.tblname}} struct{ ",
	{{- range $column := .tableColumns}}
	"\t{{$column.ColumnName|title}}\t{{$column.KindStr}}\t//DbOrder:{{$column.DbOrder}},KindSize:{{$column.KindSize}},DataType:{{$column.DataType}},ColumnType:{{$column.ColumnType}},ColumnComment:{{$column.ColumnComment}}",
	{{- end}}
	"}",
	}

	return strings.Join(tmp, "\r\n")
}

func (t *T_{{.tblname}}) CurrentTime() string{
	return time.Now().Format("{{.timeFormatLayout}}")
}

`

var tableToStructT = template.Must(template.New("table_to_struct").
	Funcs(template.FuncMap{
		"title": strings.Title,
	}).
	Parse(tmplTableToStruct))

func ToStruct(tblname string, tableColumns TableColumnSice) (string, error) {
	var buf bytes.Buffer

	//fields := make([]string, 0, len(tableColumns))
	fieldsStr := ""
	for i, column := range tableColumns {
		if i > 0 {
			fieldsStr += ","
		}
		fieldsStr += "\"" + column.ColumnName + "\""
		//fields = append(fields, column.ColumnName)
	}

	sort.Sort(tableColumns)
	data := map[string]interface{}{
		"tblname":          tblname,
		"tableColumns":     tableColumns,
		"fieldsStr":        fieldsStr,
		"created":          time.Now().Format(conf.C_time_layout),
		"timeFormatLayout": conf.C_time_layout,
	}
	err := tableToStructT.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
