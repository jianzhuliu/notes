/*
columns 表测试,生成日期 "2020-11-23 17:13:19"
*/
package main

import (
	"fmt"
	"gomysql/models"
)

func main() {
	model := models.NewTcolumns()

	fmt.Println("model:")
	fmt.Println(model)
	fmt.Println("columns:", model.Columns())
	fmt.Println("current time:", model.CurrentTime())
}
