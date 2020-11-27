/*
restfull api
*/
package web

import (
	"encoding/json"
	"net/http"

	"gomysql/conf"
	"gomysql/db"
	"gomysql/log"
	"gomysql/utils"
)

var Idb db.Idb

type H map[string]interface{}

var apiList = map[string]string{}

//回包错误码 code 定义
type ReturnCode int

const (
	CodeSucc ReturnCode = iota
	CodeFail
)

//对外入口，根据传入地址，开启web服务
func HandleWeb(addr string) {
	var ok bool
	Idb, ok = db.GetDb(conf.V_db_driver)
	if !ok {
		log.ExitOnError("the db driver=%s has not registered", conf.V_db_driver)
	}

	log.Info("web is start at %s", addr)

	//首页
	http.HandleFunc("/", handleIndex)

	//添加路由
	registerRoute("version", "查看数据库版本", handleVersion)
	registerRoute("databases", "查看所有数据库列表", handleDatabases)
	registerRoute("tables", "查看所有表", handleTables)
	registerRoute("createsql", "查看表创建sql,带get参数 tblname", handleTableCreateSql)
	registerRoute("count", "查看表总记录数,带get参数 tblname", handleTableCount)
	registerRoute("all", "列出表数据,带get参数 tblname,perpage,pagenum", handleTableAll)

	//自动开启
	go utils.OpenBrowser(addr)

	//开启 web
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error("web|http.ListenAndServe()|%s|fail|err=%v", addr, err)
	}
}

//注册路由
func registerRoute(name, desc string, handler func(http.ResponseWriter, *http.Request)) {
	pattern := "/" + name
	apiList[pattern] = desc
	http.HandleFunc(pattern, handler)
}

//统一响应方法，状态码，错误信息，数据
func response(wr http.ResponseWriter, code ReturnCode, err_msg string, data interface{}) {
	wr.WriteHeader(http.StatusOK)
	wr.Header().Set("Content-Type", "application/json;charset=utf8")

	obj := H{
		"code":    code,
		"err_msg": err_msg,
		"data":    data,
	}

	//json格式回包
	encoder := json.NewEncoder(wr)
	if err := encoder.Encode(obj); err != nil {
		log.Error("web|encoder.Encode()|fail|err=%v", err)
		wr.Write([]byte(`{"code":500,"err_msg":"server internal error"}`))
	}
}

//出错时调用
func fail(wr http.ResponseWriter, code ReturnCode, err_msg string) {
	response(wr, code, err_msg, nil)
}

//成功时调用
func succ(wr http.ResponseWriter, data interface{}) {
	response(wr, CodeSucc, "", data)
}

//获取参数值
func param(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}
