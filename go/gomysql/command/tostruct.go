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
	"flag"
	"os"
	"time"
	"encoding/json"
	
	"gomysql/models"
	"gomysql/db"
	"gomysql/log"
)

var modelObj *models.Tobj_%[1]s
var subCommandFunc map[string]func()

func init(){
	db, err := db.GetMysqlDb()
	
	if err != nil {
		log.ExitOnError("db.GetMysqlDb() | err=%%v", err)
	}
	
	modelObj = models.NewTobj_%[1]s(db)
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
func usage(){
	fmt.Println("*************************************")
	fmt.Println("Usage of table_%[1]s:")
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
	if fn,ok := subCommandFunc[cmdName]; ok {
		//计时
		defer func(beginTime time.Time){
			fmt.Println("done -- spend -- %%v", time.Since(beginTime))
		}(time.Now())
		
		fn()
	} else {
		fmt.Printf("command name is not exist, %%s \n", cmdName)
		flag.Usage()
	}
}

//基本信息
func _info(){
	fmt.Println("==============info====================begin")
	fmt.Println("current time:", modelObj.CurrentTime())
	fmt.Println(modelObj.Informaton())
	fmt.Println("columns:", modelObj.Columns())
	
	createSql, err := modelObj.CreateTableSql()
	
	if err != nil {
		log.ExitOnError("info|modelObj.CreateTableSql() | err=%%v", err)
	}
	
	fmt.Println("create table sql:")
	fmt.Println(createSql)
	
	fmt.Println("==============info====================end")
}

//重置表
func _truncate(){
	fmt.Println("==============truncate====================begin")
	
	err := modelObj.Truncate()
	if err != nil {
		log.ExitOnError("truncate|modelObj.Truncate() | err=%%v", err)
	}
	
	fmt.Println("==============truncate====================end")
}



//删除表记录数据
func _delete(){
	fmt.Println("==============delete====================begin")
	
	if flag.NArg() < 2 {
		fmt.Println("please input the record id")
		return
	}

	id, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.ExitOnError("delete| strconv.Atoi() | err=%%v", err)
	}
	
	fmt.Println("delete| id =", id)
	
	if id > 0 {
		rowsAffected, err := modelObj.Where("id=?", id).Delete()
		if err != nil {
			log.ExitOnError("delete|modelObj.Delete() | err=%%v", err)
		}
		
		fmt.Printf("delete|rowsAffected= %%d \n", rowsAffected)
	}
	
	fmt.Println("==============delete====================end")
}


//删除表所有记录
func _deleteall(){
	fmt.Println("==============deleteall====================begin")
	
	rowsAffected, err := modelObj.Delete()
	if err != nil {
		log.ExitOnError("deleteAll|modelObj.Delete() | err=%%v", err)
	}
	
	fmt.Printf("deleteAll|rowsAffected= %%d \n", rowsAffected)
	
	fmt.Println("==============deleteAll====================end")
}

//查询所有数据
func _all(){
	fmt.Println("==============all====================begin")
	
	maxNum := 10
	
	if flag.NArg() > 1 {
		tmpNum, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("all| strconv.Atoi() | err=%%v", err)
		}
		
		maxNum = tmpNum
	}
	
	fmt.Println("insert| maxNum =", maxNum)
	
	all, err := modelObj.Limit(maxNum).All()
	if err != nil {
		log.ExitOnError("all|modelObj.All() | err=%%v", err)
	}

	if all== nil || len(all) == 0 {
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
func _one(){
	fmt.Println("==============one====================begin")
	
	id := 0
	if flag.NArg() > 1 {
		realId, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("one| strconv.Atoi() | err=%%v", err)
		}
		
		id = realId
	}

	fmt.Println("one| id =", id)
	
	one, err := fetchOne(id)
	if err != nil {
		log.ExitOnError("one|modelObj.One() | err=%%v", err)
	}
	
	if one == nil {
		log.Info("one|empty")
		return
	}

	oneData,_ := modelObj.Interface(one)
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
func _count(){
	fmt.Println("==============count====================begin")
	var num int64 
	var err error
	
	if flag.NArg() > 1 {
		id, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("one| strconv.Atoi() | err=%%v", err)
		}
		
		fmt.Println("count| id =", id)
		num, err = modelObj.Where("id=?",id).Count()
	} else {
		num, err = modelObj.Count()
	}

	if err != nil {
		log.ExitOnError("count|modelObj.Count() | err=%%v", err)
	}
	
	fmt.Printf("count|num= %%d \n", num)
	
	fmt.Println("==============count====================end")
}

//事务
func _transaction(){
	id := 0
	isCommit := false
	
	if flag.NArg() > 1 {
		realId, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.ExitOnError("one| strconv.Atoi() | err=%%v", err)
		}
		
		id = realId
		
		if flag.NArg() > 2 {
			isCommit, err = strconv.ParseBool(flag.Arg(2))
			if err != nil {
				log.ExitOnError("one| strconv.ParseBool() | err=%%v", err)
			}
		}
	}

	fmt.Println("transaction| id =", id)
	
	one, err := fetchOne(id)
	if err != nil {
		log.ExitOnError("one|modelObj.One() | err=%%v", err)
	}
	
	if one == nil {
		log.Info("one|empty")
		return
	}

	//开启事务
	err = modelObj.Begin()
	if err != nil {
		log.ExitOnError("transaction| modelObj.Begin() | err=%%v", err)
	}
	
	_, _ = modelObj.Where("id=?", id).Delete()
	
	if isCommit {
		//提交事务
		err = modelObj.Commit()
	
		if err != nil {
			log.ExitOnError("transaction| modelObj.Commit() | err=%%v", err)
		}
		
		//验证
		if num, err := modelObj.Where("id=?", id).Count();num > 0 {
			fmt.Println("transaction commit fail", num, err)
		} else {
			fmt.Println("transaction commit success")
		}
	
		return 
	}
	err = modelObj.Rollback()
	
	if err != nil {
		log.ExitOnError("transaction| modelObj.Rollback() | err=%%v", err)
	}
	
	//验证
	if num, err := modelObj.Cancel().Where("id=?", id).Count(); num != 1 {
		fmt.Println("transaction rollback fail", num, err)
	} else {
		fmt.Println("transaction rollback success")
	}
}

//更新表记录数据
func _update(){
	fmt.Println("==============update====================begin")
	
	if flag.NArg() < 2 {
		fmt.Println("please input the record id")
		return
	}

	id, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.ExitOnError("update| strconv.Atoi() | err=%%v", err)
	}
	
	fmt.Println("update| id =", id)
	
	if id > 0 {
		values := map[string]interface{}{
			"Name":"update_name" ,
			"Phone": "99999999",
			"Info":"update_info",
		}
		
		rowsAffected, err := modelObj.Where("id=?", id).Update(values)
		if err != nil {
			log.ExitOnError("update|modelObj.Update() | err=%%v", err)
		}
		
		fmt.Printf("update|rowsAffected= %%d \n", rowsAffected)
		
		one, err := fetchOne(id)
		if err != nil {
			log.ExitOnError("update|modelObj.One() | err=%%v", err)
		}
		
		if one == nil {
			log.Info("one|empty")
			return
		}

		oneData,_ := modelObj.Interface(one)
		fmt.Println(oneData)
		out, err := json.MarshalIndent(oneData, "", " ")
		fmt.Println(string(out))
	
	}
	
	fmt.Println("==============update====================end")
}

//插入测试数据
func _insert(){
	fmt.Println("==============insert====================begin")
	
	if flag.NArg() < 2 {
		fmt.Println("please input the insert max num")
		return
	}

	maxNum, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.ExitOnError("insert| strconv.Atoi() | err=%%v", err)
	}
	fmt.Println("insert| maxNum =", maxNum)
	
	for i:=1; i<=maxNum; i++ {
		values := map[string]interface{}{
			"Status":i%%2,
			"Name":"name_insert_" +  strconv.Itoa(i),
			"Phone": "129833444" + strconv.Itoa(i),
			"Info":"info_insert" +  strconv.Itoa(i),
		}
		
		lastInsertId, err := modelObj.Insert(values)
		if err != nil {
			log.ExitOnError("insert|modelObj.Insert() | err=%%v", err)
		}
		
		fmt.Printf("insert|lastInsertId= %%d \n", lastInsertId)
	}
	
	fmt.Println("==============insert====================end")
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
