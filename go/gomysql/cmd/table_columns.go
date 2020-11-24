/*
columns 表测试,生成日期 "2020-11-24 18:23:36"
*/
package main

import (
	"fmt"
	"gomysql/db"
	"gomysql/log"
	"gomysql/models"
)

func main() {
	db, err := db.GetMysqlDb()

	if err != nil {
		log.ExitOnError("db.GetMysqlDb() | err=%!v(MISSING)", err)
	}

	modelObj := models.NewTobj_columns(db)

	fmt.Println(modelObj.Informaton())
	fmt.Println("columns:", modelObj.Columns())
	fmt.Println("current time:", modelObj.CurrentTime())

	fmt.Println("=============one=================")
	one, err := modelObj.Where("status =?", 1).One()
	if err != nil {
		log.ExitOnError("modelObj.One() | err=2020-11-24 18:23:36", err)
	}

	oneData, ok := modelObj.Interface(one)
	fmt.Println(ok, oneData.Id, oneData)

	all, err := modelObj.Where("status =?", 1).OrderBy("id desc").Limit(2).All()
	if err != nil {
		log.ExitOnError("modelObj.All() | err=%!v(MISSING)", err)
	}

	fmt.Println("=============all=================")
	for _, data := range all {
		realData, ok := modelObj.Interface(data)
		fmt.Println(ok, realData.Id, realData)
	}
}
