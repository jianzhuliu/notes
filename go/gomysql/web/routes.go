/*
注册路由
*/
package web

import (
	"fmt"
	"gomysql/log"
	"gomysql/models"
	"net/http"
	"strconv"
)

//首页
func handleIndex(wr http.ResponseWriter, r *http.Request) {
	succ(wr, H{"api list": apiList})
}

//查看版本号信息
func handleVersion(wr http.ResponseWriter, r *http.Request) {
	version, err := Idb.Version()
	if err != nil {
		log.Error("web|Idb.Version()|fail|err=%v", err)
		fail(wr, CodeFail, "cannot get the db version")
		return
	}

	succ(wr, H{"version": version})
}

//查看数据库列表
func handleDatabases(wr http.ResponseWriter, r *http.Request) {
	databases, err := Idb.Databases()
	if err != nil {
		log.Error("web|Idb.Databases()|fail|err=%v", err)
		fail(wr, CodeFail, "cannot get the db databases")
		return
	}

	succ(wr, H{"databases": databases})
}

//查看所有表
func handleTables(wr http.ResponseWriter, r *http.Request) {
	tables, err := Idb.Tables()
	if err != nil {
		log.Error("web|Idb.Tables()|fail|err=%v", err)
		fail(wr, CodeFail, "cannot get the db tables")
		return
	}

	succ(wr, H{"tables": tables})
}

//查看表创建 sql
func handleTableCreateSql(wr http.ResponseWriter, r *http.Request) {
	tblname := param(r, "tblname")
	if tblname == "" {
		fail(wr, CodeFail, "param tblname is empty")
		return
	}

	create_table_sql, err := Idb.ShowCreateTableSql(tblname)
	if err != nil {
		log.Error("web|Idb.ShowCreateTableSql()|fail|err=%v", err)
		fail(wr, CodeFail, fmt.Sprintf("cannot get the table(%s) create sql", tblname))
		return
	}

	succ(wr, H{"create_table_sql": create_table_sql, "tblname": tblname})
}

//查看表总记录数
func handleTableCount(wr http.ResponseWriter, r *http.Request) {
	tblname := param(r, "tblname")
	if tblname == "" {
		fail(wr, CodeFail, "param tblname is empty")
		return
	}

	//获取创建表操作对象的函数
	fn, ok := models.TableToObjCreateFunc[tblname]
	if !ok {
		fail(wr, CodeFail, fmt.Sprintf("table(%[1]s) has not create models, please do command go run main.go tostruct -table=%[1]s -out", tblname))
		return
	}

	db := Idb.Db()
	modelsObj := fn(db)

	num, err := modelsObj.Count()
	if err != nil {
		log.Error("web|modelObj.Count()|fail|err=%v", err)
		fail(wr, CodeFail, fmt.Sprintf("cannot get the table(%s) total record", tblname))
		return
	}

	succ(wr, H{"count": num, "tblname": tblname})
}

//列出表数据
func handleTableAll(wr http.ResponseWriter, r *http.Request) {
	tblname := param(r, "tblname")
	if tblname == "" {
		fail(wr, CodeFail, "param tblname is empty")
		return
	}

	perpage := 10
	pagenum := 1

	paramPerPage := param(r, "perpage") //每页显示
	paramPageNum := param(r, "pagenum") //第几页
	log.Info("perpage = %s, pagenum= %s", paramPerPage, paramPageNum)

	if paramPerPage != "" {
		realPerPage, err := strconv.Atoi(paramPerPage)
		if err == nil {
			perpage = realPerPage
		}
	}

	if paramPageNum != "" {
		realPageNum, err := strconv.Atoi(paramPageNum)
		if err == nil {
			pagenum = realPageNum
			if pagenum < 1 {
				pagenum = 1
			}
		}
	}

	//获取创建表操作对象的函数
	fn, ok := models.TableToObjCreateFunc[tblname]
	if !ok {
		fail(wr, CodeFail, fmt.Sprintf("table(%[1]s) has not create models, please do command go run main.go tostruct -table=%[1]s -out", tblname))
		return
	}

	db := Idb.Db()
	modelsObj := fn(db)

	//获取总记录数
	total, err := modelsObj.Count()

	if int(total) < pagenum*perpage {
		pagenum = 1
	}

	offset := (pagenum - 1) * perpage
	log.Info("tblname=%s, perpage=%d, pagenum=%d, offset=%d", tblname, perpage, pagenum, offset)

	//分页查询数据
	data, err := modelsObj.Limit(offset, perpage).All()
	if err != nil {
		log.Error("web|modelObj.Count()|fail|err=%v", err)
		fail(wr, CodeFail, fmt.Sprintf("cannot get the table(%s) total record", tblname))
		return
	}

	succ(wr, H{"data": data, "tblname": tblname})
}
