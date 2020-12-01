package models

//*
import (
	"strconv"
	"time"
)

//生成测试数据
func (t *Tobj_columns) GenTestForUpdate() map[string]interface{} {
	values := map[string]interface{}{
		"Name":  "update_name",
		"Phone": "99999999",
		"Info":  "update_info",
	}
	return values
}

func (t *Tobj_columns) GenTestForInsert(num int) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, num)

	for i := 1; i <= num; i++ {
		values := map[string]interface{}{
			"Status":  i % 2,
			"Name":    "name_insert_" + strconv.Itoa(i),
			"Phone":   "129833444" + strconv.Itoa(i),
			"Info":    "info_insert" + strconv.Itoa(i),
			"Created": TimeNormal{time.Now()},
		}

		result = append(result, values)
	}

	return result
}

//*/
