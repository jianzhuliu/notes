/*
mysql 相关变量
*/
package db

//mysql 字段类型对应 kind
//根据前4个字符，来绝对匹配
var TypeMysqlToKind = map[string]string{
	//日期类型， date,datetime,time,timestamp,year
	"date": "TimeNormal",
	"time": "TimeNormal",
	"year": "int16",

	//小数, float, double, decimal
	//浮点数
	"floa": "float32",
	"doub": "float64",
	"decc": "float64",

	//字符串类型， varchar,char,enum,text,blob,set,binary,varbinary
	//tinytext,mediumtext,longtext,tinyblob,mediumblob,longblob
	"char": "string",
	"varc": "string",
	"enum": "string",

	"text": "string",
	"blob": "[]byte",

	"bina": "[]byte",
	"varb": "[]byte",

	//json
	"json": "string",

	//整数类型， tinyint,smallint,mediumint,int,integer,bigint,bit
	"bigi": "int64",
	"inte": "int64",

	//特殊长度为3的类型
	"set": "string",
	"int": "int",
	"bit": "[]uint8",
}

//完整字段类型搜索
var WholeTypeMysqlToKind = map[string]string{
	"tinyblob":   "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",

	"tinytext":   "string",
	"mediumtext": "string",
	"longtext":   "string",

	"tinyint":   "int8",
	"smallint":  "int16",
	"mediumint": "int32",
}

//mysql 无符号类型
var UnsignedTypeMysqlToKind = map[string]string{
	//float, double
	"floa": "float32",
	"doub": "float64",

	//int,tinyint,smallint,mediumint,bigint
	"int":  "uint",
	"tiny": "uint8",
	"smal": "uint16",
	"medi": "uint32",
	"bigi": "uint64",

	"inte": "uint64",
}
