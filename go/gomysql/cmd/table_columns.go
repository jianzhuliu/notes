/*
columns 表测试,生成日期 "2020-11-25 12:29:22"
*/
package main

import (
	"fmt"
	"strconv"

	"gomysql/db"
	"gomysql/log"
	"gomysql/models"
)

func main() {
	db, err := db.GetMysqlDb()

	if err != nil {
		log.ExitOnError("db.GetMysqlDb() | err=%v", err)
	}

	modelObj := models.NewTobj_columns(db)

	fmt.Println(modelObj.Informaton())
	fmt.Println("columns:", modelObj.Columns())
	fmt.Println("current time:", modelObj.CurrentTime())

	all(modelObj)
	deleteAll(modelObj)

	//插入测试数据
	var lastInsertId int64
	var values map[string]interface{}

	for i := 1; i < 10; i++ {
		values = map[string]interface{}{
			"Status": i % 2,
			"Name":   "name_insert_" + strconv.Itoa(i),
			"Phone":  "129833444" + strconv.Itoa(i),
			"Info":   "info_insert" + strconv.Itoa(i),
		}
		lastInsertId = insert(modelObj, values)
	}

	all(modelObj)

	if lastInsertId > 0 {
		values["Info"] = "update_info"
		update(modelObj, lastInsertId, values)
	}

	all(modelObj)

	oneData := one(modelObj)

	delete(modelObj, oneData.Id)
	_ = one(modelObj)

	all(modelObj)
}

//查询一条记录
func one(modelObj *models.Tobj_columns) models.T_columns {
	fmt.Println("=============one=================")
	one, err := modelObj.Where("status =?", 1).One()
	if err != nil {
		log.ExitOnError("modelObj.One() | err=%v", err)
	}

	if one == nil {
		log.Info("one|empty")
		return models.T_columns{}
	}

	oneData, ok := modelObj.Interface(one)
	fmt.Println(ok, oneData.Id, oneData)
	return oneData
}

//查看所有记录
func all(modelObj *models.Tobj_columns) {
	fmt.Println("=============all=================")

	all, err := modelObj.Where("status =?", 1).OrderBy("id desc").Limit(10).All()
	if err != nil {
		log.ExitOnError("modelObj.All() | err=%v", err)
	}

	if all == nil || len(all) == 0 {
		log.Info("all|empty")
		return
	}

	for _, data := range all {
		realData, ok := modelObj.Interface(data)
		fmt.Println(ok, realData.Id, realData)
	}
}

//删除所有
func deleteAll(modelObj *models.Tobj_columns) {
	fmt.Println("=============deleteAll=================")

	rowsAffected, err := modelObj.Delete()
	if err != nil {
		log.ExitOnError("modelObj.Delete() | err=%v", err)
	}

	fmt.Printf("deleteAll|rowsAffected= %d \n", rowsAffected)
}

//删除
func delete(modelObj *models.Tobj_columns, id uint) {
	fmt.Println("=============delete=================")

	if id == 0 {
		log.Info("delete| id == 0")
		return
	}

	rowsAffected, err := modelObj.Where("id=?", id).Delete()
	if err != nil {
		log.ExitOnError("modelObj.Delete() | err=%v", err)
	}

	fmt.Printf("delete|rowsAffected= %d \n", rowsAffected)
}

//插入
func insert(modelObj *models.Tobj_columns, values map[string]interface{}) int64 {
	fmt.Println("=============insert=================")
	if len(values) == 0 {
		log.Info("insert| values is empty")
		return 0
	}

	lastInsertId, err := modelObj.Insert(values)
	if err != nil {
		log.ExitOnError("modelObj.Insert() | err=%v", err)
	}

	fmt.Printf("insert|lastInsertId= %d \n", lastInsertId)

	return lastInsertId
}

//更新
func update(modelObj *models.Tobj_columns, id int64, values map[string]interface{}) {
	fmt.Println("=============update=================")

	if id == 0 {
		log.Info("update| id == 0")
		return
	}

	if len(values) == 0 {
		log.Info("update| values is empty")
		return
	}

	rowsAffected, err := modelObj.Where("id=?", id).Update(values)
	if err != nil {
		log.ExitOnError("modelObj.Update() | err=%v", err)
	}

	fmt.Printf("update|rowsAffected= %d \n", rowsAffected)
}
