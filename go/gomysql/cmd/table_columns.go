/*
columns 表测试,生成日期 "2020-11-30 12:34:33"
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"gomysql/db"
	"gomysql/log"
	"gomysql/models"
)

var modelObj *models.Tobj_columns
var subCommandFunc map[string]func()

func init() {
	db, err := db.GetMysqlDb()

	if err != nil {
		log.ExitOnError("db.GetMysqlDb() | err=%v", err)
	}

	modelObj = models.NewTobj_columns(db)
	flag.Usage = usage

	//定义各个子命令，对应处理函数
	subCommandFunc = make(map[string]func())
	subCommandFunc["info"] = _info
	subCommandFunc["truncate"] = _truncate
	subCommandFunc["insert"] = _insert
	subCommandFunc["update"] = _update
	subCommandFunc["delete"] = _delete
	subCommandFunc["deleteall"] = _deleteall
	subCommandFunc["all"] = _all
	subCommandFunc["one"] = _one
	subCommandFunc["count"] = _count
	subCommandFunc["transaction"] = _transaction
}

//帮助说明函数
func usage() {
	fmt.Println("*************************************")
	fmt.Println("Usage of table_columns:")
	fmt.Println("info		展示表基本信息")
	fmt.Println("truncate	重置表")
	fmt.Println("insert		插入测试数据,后面接插入记录数")
	fmt.Println("update		更新表记录，后接记录id")
	fmt.Println("delete		删除表记录，后接记录id")
	fmt.Println("deleteall	删除表所有记录")
	fmt.Println("all		查询所有数据,后面接查询记录数")
	fmt.Println("one		查询一条记录，后接记录id")
	fmt.Println("count		统计记录条数")
	fmt.Println("transaction		事务")

	//输出默认参数
	flag.PrintDefaults()

	os.Exit(0)
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	cmdName := flag.Arg(0)
	if fn, ok := subCommandFunc[cmdName]; ok {
		//计时
		defer func(beginTime time.Time) {
			fmt.Println("done -- spend -- %v", time.Since(beginTime))
		}(time.Now())

		fn()
	} else {
		fmt.Printf("command name is not exist, %s \n", cmdName)
		flag.Usage()
	}
}

//基本信息
func _info() {
	fmt.Println("==============info====================begin")
	fmt.Println("current time:", modelObj.CurrentTime())
	fmt.Println(modelObj.Informaton())
	fmt.Println("columns:", modelObj.Columns())

	createSql, err := modelObj.CreateTableSql()

	if err != nil {
		log.ExitOnError("info|modelObj.CreateTableSql() | err=%v", err)
	}

	fmt.Println("create table sql:")
	fmt.Println(createSql)

	fmt.Println("==============info====================end")
}

//重置表
func _truncate() {
	fmt.Println("==============truncate====================begin")

	err := modelObj.Truncate()
	if err != nil {
		log.ExitOnError("truncate|modelObj.Truncate() | err=%v", err)
	}

	fmt.Println("==============truncate====================end")
}

//删除表记录数据
func _delete() {
	fmt.Println("==============delete====================begin")

	if flag.NArg() < 2 {
		fmt.Println("please input the record id")
		return
	}

	id, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.ExitOnError("delete| strconv.Atoi() | err=%v", err)
	}

	fmt.Println("delete| id =", id)

	if id > 0 {
		rowsAffected, err := modelObj.Where("id=?", id).Delete()
		if err != nil {
			log.ExitOnError("delete|modelObj.Delete() | err=%v", err)
		}

		fmt.Printf("delete|rowsAffected= %d \n", rowsAffected)
	}

	fmt.Println("==============delete====================end")
}

//删除表所有记录
func _deleteall() {
	fmt.Println("==============deleteall====================begin")

	rowsAffected, err := modelObj.Delete()
	if err != nil {
		log.ExitOnError("deleteAll|modelObj.Delete() | err=%v", err)
	}

	fmt.Printf("deleteAll|rowsAffected= %d \n", rowsAffected)

	fmt.Println("==============deleteAll====================end")
}

//查询所有数据
func _all() {
	fmt.Println("==============all====================begin")

	maxNum := 10

	if flag.NArg() > 1 {
		tmpNum, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("all| strconv.Atoi() | err=%v", err)
		}

		maxNum = tmpNum
	}

	fmt.Println("insert| maxNum =", maxNum)

	all, err := modelObj.Limit(maxNum).All()
	if err != nil {
		log.ExitOnError("all|modelObj.All() | err=%v", err)
	}

	if all == nil || len(all) == 0 {
		log.Info("all|empty")
		return
	}

	for _, data := range all {
		realData, ok := modelObj.Interface(data)
		if ok {
			fmt.Println(realData)
		}
	}

	fmt.Println("==============all====================end")
}

//查询一条记录
func _one() {
	fmt.Println("==============one====================begin")

	id := 0
	if flag.NArg() > 1 {
		realId, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("one| strconv.Atoi() | err=%v", err)
		}

		id = realId
	}

	fmt.Println("one| id =", id)

	one, err := fetchOne(id)
	if err != nil {
		log.ExitOnError("one|modelObj.One() | err=%v", err)
	}

	if one == nil {
		log.Info("one|empty")
		return
	}

	oneData, _ := modelObj.Interface(one)
	fmt.Println(oneData)
	out, err := json.MarshalIndent(oneData, "", " ")
	fmt.Println(string(out))

	fmt.Println("==============one====================end")
}

//查询一条记录
func fetchOne(id int) (interface{}, error) {
	one, err := modelObj.Where("id >= ?", id).One()
	if err != nil {
		return nil, err
	}

	return one, nil
}

//统计记录条数
func _count() {
	fmt.Println("==============count====================begin")
	var num int64
	var err error

	if flag.NArg() > 1 {
		id, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("one| strconv.Atoi() | err=%v", err)
		}

		fmt.Println("count| id =", id)
		num, err = modelObj.Where("id=?", id).Count()
	} else {
		num, err = modelObj.Count()
	}

	if err != nil {
		log.ExitOnError("count|modelObj.Count() | err=%v", err)
	}

	fmt.Printf("count|num= %d \n", num)

	fmt.Println("==============count====================end")
}

//事务
func _transaction() {
	id := 0
	isCommit := false

	if flag.NArg() > 1 {
		realId, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("one| strconv.Atoi() | err=%v", err)
		}

		id = realId

		if flag.NArg() > 2 {
			isCommit, err = strconv.ParseBool(flag.Arg(2))
			if err != nil {
				log.ExitOnError("one| strconv.ParseBool() | err=%v", err)
			}
		}
	}

	fmt.Println("transaction| id =", id)

	one, err := fetchOne(id)
	if err != nil {
		log.ExitOnError("one|modelObj.One() | err=%v", err)
	}

	if one == nil {
		log.Info("one|empty")
		return
	}

	//开启事务
	err = modelObj.Begin()
	if err != nil {
		log.ExitOnError("transaction| modelObj.Begin() | err=%v", err)
	}

	_, _ = modelObj.Where("id=?", id).Delete()

	if isCommit {
		//提交事务
		err = modelObj.Commit()

		if err != nil {
			log.ExitOnError("transaction| modelObj.Commit() | err=%v", err)
		}

		//验证
		if num, err := modelObj.Where("id=?", id).Count(); num > 0 {
			fmt.Println("transaction commit fail", num, err)
		} else {
			fmt.Println("transaction commit success")
		}

		return
	}
	err = modelObj.Rollback()

	if err != nil {
		log.ExitOnError("transaction| modelObj.Rollback() | err=%v", err)
	}

	//验证
	if num, err := modelObj.Cancel().Where("id=?", id).Count(); num != 1 {
		fmt.Println("transaction rollback fail", num, err)
	} else {
		fmt.Println("transaction rollback success")
	}
}

//更新表记录数据
func _update() {
	fmt.Println("==============update====================begin")

	if flag.NArg() < 2 {
		fmt.Println("please input the record id")
		return
	}

	id, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.ExitOnError("update| strconv.Atoi() | err=%v", err)
	}

	fmt.Println("update| id =", id)

	if id > 0 {
		values := map[string]interface{}{
			"Name":  "update_name",
			"Phone": "99999999",
			"Info":  "update_info",
		}

		rowsAffected, err := modelObj.Where("id=?", id).Update(values)
		if err != nil {
			log.ExitOnError("update|modelObj.Update() | err=%v", err)
		}

		fmt.Printf("update|rowsAffected= %d \n", rowsAffected)

		one, err := fetchOne(id)
		if err != nil {
			log.ExitOnError("update|modelObj.One() | err=%v", err)
		}

		if one == nil {
			log.Info("one|empty")
			return
		}

		oneData, _ := modelObj.Interface(one)
		fmt.Println(oneData)
		out, err := json.MarshalIndent(oneData, "", " ")
		if err != nil {
			log.Error("json.MarshalIndent() | err =2020-11-30 12:34:33", err)
		}
		fmt.Println(string(out))

	}

	fmt.Println("==============update====================end")
}

//插入测试数据
func _insert() {
	fmt.Println("==============insert====================begin")

	if flag.NArg() < 2 {
		fmt.Println("please input the insert max num")
		return
	}

	maxNum, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.ExitOnError("insert| strconv.Atoi() | err=%v", err)
	}
	fmt.Println("insert| maxNum =", maxNum)

	for i := 1; i <= maxNum; i++ {
		values := map[string]interface{}{
			"Status":  i % 2,
			"Name":    "name_insert_" + strconv.Itoa(i),
			"Phone":   "129833444" + strconv.Itoa(i),
			"Info":    "info_insert" + strconv.Itoa(i),
			"Created": models.TimeNormal{time.Now()},
		}

		lastInsertId, err := modelObj.Insert(values)
		if err != nil {
			log.ExitOnError("insert|modelObj.Insert() | err=%v", err)
		}

		fmt.Printf("insert|lastInsertId= %d \n", lastInsertId)
	}

	fmt.Println("==============insert====================end")
}
