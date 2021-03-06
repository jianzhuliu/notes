/*
本配置文件不属于框架结构内，仅方便测试
*/
package conf

//数据库引擎类型
type DriverType int

const (
	DriverMysql DriverType = iota + 1
	DriverSqlite
)

//根据数据库引擎类型，获取 data source name
func GetDsnByDriver(driverType DriverType) (driver, dsn string) {
	switch driverType {
	case DriverMysql:
		driver = "mysql"
		dsn = "root:@tcp(127.0.0.1:3306)/goorm?parseTime=true&loc=Local"
	case DriverSqlite:
		driver = "sqlite3"
		dsn = "goorm.db"
	}

	return
}


type User struct {
	Id   int `goorm:"primary key"`
	Name string
}

//自定义表名
func (u *User) TableName() string {
	return "user"
}

var User1 = &User{Id:1,Name:"name1"}
var User2 = &User{Id:2,Name:"name2"}