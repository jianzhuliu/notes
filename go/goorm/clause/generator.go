package clause

import (
	"fmt"
	"strings"
)

//生成sql语句需要的前缀与参数
type generator func(...interface{}) (string, []interface{})

var generators = map[OpType]generator{}

func init() {
	generators[OpTypeInsert] = _insert
	generators[OpTypeSelect] = _select
	generators[OpTypeWhere] = _where
	generators[OpTypeOrderBy] = _orderby
	generators[OpTypeLimit] = _limit
	generators[OpTypeValues] = _values

}

//insert into table (?,?)
func _insert(values ...interface{}) (string, []interface{}) {
	tblname := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("insert into %s (%s)", tblname, fields), nil
}

//select * from ?
func _select(values ...interface{}) (string, []interface{}) {
	tblname := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %s from %s", fields, tblname), nil
}

//where
func _where(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("where %s", values[0]), values[1:]
}

//order by id desc
func _orderby(values ...interface{}) (string, []interface{}) {
	desc := []string{}
	for _, v := range values {
		desc = append(desc, v.(string))
	}

	fields := strings.Join(desc, ",")
	return fmt.Sprintf("order by %s", fields), nil
}

//limit 10
func _limit(values ...interface{}) (string, []interface{}) {
	bindStr := genBindVars(len(values))
	return fmt.Sprintf("limit %s", bindStr), values
}

//values(?,?)(?,?)
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var args []interface{}

	sql.WriteString("values")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 < len(values) {
			sql.WriteString(",")
		}

		args = append(args, v...)
	}

	return sql.String(), args
}

//根据值个数，生成需要的参数位置信息
func genBindVars(num int) string {
	args := make([]string, 0, num)
	for i := 0; i < num; i++ {
		args = append(args, "?")
	}
	return strings.Join(args, ",")
}
