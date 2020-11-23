package db

import (
	"bytes"
	"strings"
	"sort"
	"time"
	
	"text/template"
	"gomysql/conf"
)

/////////////////////////////////表结构转换为 struct 

const tmplTableToStruct = `

//database {{.tblname}} to struct created at {{.created}}
type T_{{.tblname|title}} Struct{
{{- range $column := .tableColumns}}
	{{$column.ColumnName|title}}		{{$column.KindStr}}		//DbOrder:{{$column.DbOrder}},KindSize:{{$column.KindSize}},DataType:{{$column.DataType}},ColumnType:{{$column.ColumnType}}
{{- end }}
}
`

var tableToStructT = template.Must(template.New("table_to_struct").
	Funcs(template.FuncMap{
		"title":strings.Title,
	}).
	Parse(tmplTableToStruct))

func ToStruct(tblname string, tableColumns TableColumnSice) (string, error){
	var buf bytes.Buffer
	
	sort.Sort(tableColumns)
	data := map[string]interface{}{
		"tblname": tblname,
		"tableColumns":tableColumns,
		"created":time.Now().Format(conf.C_time_layout),
	}
	err := tableToStructT.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
