//the sub command "tostruct", created at "2020-11-22 17:37:21"
package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gomysql/conf"
	"gomysql/db"
	"gomysql/utils"
)

var (
	all bool //是否处理所有表
	out bool //生成文件
)

func init() {
	//新建子命令
	subCommand := NewSubCommand("tostruct", "format db table to struct")

	//子命令配置执行函数
	subCommand.SetRun(RunTostruct)

	//设置解析参数前处理
	subCommand.SetBeforeParse(BeforeParseTostruct)

	//添加子命令
	AddCommand(subCommand)
}

//执行之前的处理，比如重新设置参数默认值
func BeforeParseTostruct(sub *SubCommand) error {
	//添加自定义参数
	sub.BoolVar(&all, "all", false, "format all table to struct")
	sub.BoolVar(&out, "out", false, "gen to file, the default folder is models")

	/*
		//取消验证数据库名
		sub.SetFlagValue("check_database", "false")
		//*/

	//*
	//取消验证表名
	sub.SetFlagValue("check_table", "false")
	//*/

	return nil
}

//查看数据库版本号
func RunTostruct() error {
	//参数校验
	//*
	if !all && len(conf.V_db_table) == 0 {
		return fmt.Errorf("please set the table name, -table or gen all table to struct, -all")
	}
	//*/

	Idb, ok := db.GetDb(conf.V_db_driver)
	if !ok {
		return fmt.Errorf("the db driver=%s has not registered", conf.V_db_driver)
	}

	if len(conf.V_db_table) > 0 {
		//处理单个表
		tableColumns, err := Idb.Fields(conf.V_db_database, conf.V_db_table)

		if err != nil {
			return err
		}

		/*
			fmt.Printf("%s of %s to kind:\n", conf.V_db_table, conf.V_db_database)
			for _, column := range tableColumns {
				fmt.Printf("%-20s ========> %-20s \n", column.ColumnName, column.KindStr)
			}
			//*/

		str, err := db.ToStruct(conf.V_db_table, tableColumns)
		if err != nil {
			return err
		}

		if out {
			//生成文件
			err = genModelFile(conf.V_db_database, conf.V_db_table, str)
			if err != nil {
				return err
			}

			err = genTbaseFile()
			if err != nil {
				return err
			}

			err = genCmdFile(conf.V_db_table)
			if err != nil {
				return err
			}

			fmt.Printf("table %s has finished \n", conf.V_db_table)
		} else {
			fmt.Println("====================\n", str)
		}

		return nil
	} else {
		//处理所有表
		tables, err := Idb.Tables()
		if err != nil {
			return err
		}

		if len(tables) == 0 {
			return fmt.Errorf("this is no table yet from database %s", conf.V_db_database)
		}

		allTables := make(map[string][]db.TableColumn, len(tables))

		for _, tblname := range tables {
			tableColumns, err := Idb.Fields(conf.V_db_database, tblname)
			if err != nil {
				return err
			}

			allTables[tblname] = tableColumns
		}

		fmt.Printf("all table to kind of %s :\n", conf.V_db_database)
		for tblname, tableColumns := range allTables {
			/*
				fmt.Printf("\t%-15s--------------------------------------\n", tblname)
				for _, column := range tableColumns {
					fmt.Printf("\t\t%-20s ========> %-20s \n", column.ColumnName, column.KindStr)
				}
				//*/

			str, err := db.ToStruct(tblname, tableColumns)
			if err != nil {
				return err
			}

			if out {
				//生成文件
				err = genModelFile(conf.V_db_database, tblname, str)
				if err != nil {
					return err
				}

				err = genCmdFile(tblname)
				if err != nil {
					return err
				}

				fmt.Printf("table %s has finished \n", tblname)
			} else {
				fmt.Println(str)
			}
		}

		if out {
			err = genTbaseFile()
			if err != nil {
				return err
			}
		}

		return nil
	}
}

//格式化
func gofmt(file string) {
	//上下文信息
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//执行格式化命令
	_ = exec.CommandContext(ctx, "go", "fmt", file).Run()
	return
}

//生成文件
func genModelFile(database, tblname, str string) error {
	modelsPath, err := utils.GenFolder("models")
	if err != nil {
		return err
	}

	modelsFileName := fmt.Sprintf("%s_%s.go", database, tblname)
	modelsFile := filepath.Join(modelsPath, modelsFileName)
	err = ioutil.WriteFile(modelsFile, []byte(str), os.ModePerm)
	if err != nil {
		return err
	}

	gofmt(modelsFile)

	return nil
}

//生成公用文件
func genTbaseFile() error {
	cmdPath, err := utils.GenFolder("models")
	if err != nil {
		return err
	}

	file := filepath.Join(cmdPath, "tbase.go")

	content := db.GenTbase()
	err = ioutil.WriteFile(file, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}

	gofmt(file)

	return nil
}

//表测试入口
const tmplCmdModel = `/*
%[1]s 表测试,生成日期 "%[2]s"
*/
package main

import (
	"fmt"
	"strconv"
	
	"gomysql/models"
	"gomysql/db"
	"gomysql/log"
)

func main() {
	db, err := db.GetMysqlDb()
	
	if err != nil {
		log.ExitOnError("db.GetMysqlDb() | err=%%v", err)
	}
	
	modelObj := models.NewTobj_%[1]s(db)

	fmt.Println(modelObj.Informaton())
	fmt.Println("columns:", modelObj.Columns())
	fmt.Println("current time:", modelObj.CurrentTime())
	
	all(modelObj)
	deleteAll(modelObj)
	
	//插入测试数据
	var lastInsertId int64
	var values map[string]interface{}
	
	for i:=1; i<10; i++ {
		values = map[string]interface{}{
			"Status":i%%2,
			"Name":"name_insert_" +  strconv.Itoa(i),
			"Phone": "129833444" + strconv.Itoa(i),
			"Info":"info_insert" +  strconv.Itoa(i),
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
func one(modelObj *models.Tobj_%[1]s) models.T_%[1]s {
	fmt.Println("=============one=================")
	one, err := modelObj.Where("status =?", 1).One()
	if err != nil {
		log.ExitOnError("modelObj.One() | err=%%v", err)
	}
	
	if one == nil {
		log.Info("one|empty")
		return models.T_%[1]s{}
	}

	oneData, ok := modelObj.Interface(one)
	fmt.Println(ok, oneData.Id, oneData)
	return oneData
}

//查看所有记录
func all(modelObj *models.Tobj_%[1]s){
	fmt.Println("=============all=================")

	all, err := modelObj.Where("status =?", 1).OrderBy("id desc").Limit(10).All()
	if err != nil {
		log.ExitOnError("modelObj.All() | err=%%v", err)
	}

	if all== nil || len(all) == 0 {
		log.Info("all|empty")
		return 
	}
	
	for _, data := range all {
		realData, ok := modelObj.Interface(data)
		fmt.Println(ok, realData.Id, realData)
	}
}

//删除所有
func deleteAll(modelObj *models.Tobj_%[1]s) {
	fmt.Println("=============deleteAll=================")
	
	rowsAffected, err := modelObj.Delete()
	if err != nil {
		log.ExitOnError("modelObj.Delete() | err=%%v", err)
	}
	
	fmt.Printf("deleteAll|rowsAffected= %%d \n", rowsAffected)
}

//删除
func delete(modelObj *models.Tobj_%[1]s, id uint) {
	fmt.Println("=============delete=================")
	
	if id == 0 {
		log.Info("delete| id == 0")
		return
	}
	
	rowsAffected, err := modelObj.Where("id=?", id).Delete()
	if err != nil {
		log.ExitOnError("modelObj.Delete() | err=%%v", err)
	}
	
	fmt.Printf("delete|rowsAffected= %%d \n", rowsAffected)
}

//插入
func insert(modelObj *models.Tobj_%[1]s, values map[string]interface{}) int64 {
	fmt.Println("=============insert=================")
	if len(values) == 0 {
		log.Info("insert| values is empty")
		return 0
	}
	
	lastInsertId, err := modelObj.Insert(values)
	if err != nil {
		log.ExitOnError("modelObj.Insert() | err=%%v", err)
	}
	
	fmt.Printf("insert|lastInsertId= %%d \n", lastInsertId)
	
	return lastInsertId
}

//更新
func update(modelObj *models.Tobj_%[1]s, id int64,values map[string]interface{}) {
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
		log.ExitOnError("modelObj.Update() | err=%%v", err)
	}
	
	fmt.Printf("update|rowsAffected= %%d \n", rowsAffected)
}

`

//生成 cmd file
func genCmdFile(tblname string) error {
	cmdPath, err := utils.GenFolder("cmd")
	if err != nil {
		return err
	}

	cmdFileName := fmt.Sprintf("table_%s.go", tblname)
	cmdFile := filepath.Join(cmdPath, cmdFileName)

	created := time.Now().Format(conf.C_time_layout)
	content := fmt.Sprintf(tmplCmdModel, tblname, created)
	err = ioutil.WriteFile(cmdFile, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}

	gofmt(cmdFile)

	return nil

}
