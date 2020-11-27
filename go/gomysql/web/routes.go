/*
注册路由
*/
package web

import (
	"gomysql/log"
	"net/http"
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
