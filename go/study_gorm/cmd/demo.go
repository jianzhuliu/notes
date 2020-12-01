package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"study_gorm/config"

	"gitee.com/jianzhuliu/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var ctx context.Context

func init() {
	//初始化 db 命令行参数
	conf.FlagInit(conf.Fdb)

	ctx = context.Background()
}

type User struct {
	ID        uint `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Age       uint8 `gorm:"default:10"`
}

/*
//自定义表名
func (u User) TableName() string{
	return "user"
}
//*/

//定义钩子
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

func (u *User) BeforeUpdate(db *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

func initDb() {
	//指定数据库名
	conf.V_db_dbname = config.C_database
	dsn := conf.FlagDbDsn()

	//自定义日志处理对象
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond, //慢查询
		LogLevel:      logger.Info,
		Colorful:      false, // 禁用彩色打印
	})

	//命名规则
	nameingStrategy := schema.NamingStrategy{
		TablePrefix:   "t_",
		SingularTable: true,
	}

	newLogger.Info(ctx, "the database dsn:%s", dsn)
	var err error

	//建立链接
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//全局配置
		DryRun:         false,
		Logger:         newLogger,
		NamingStrategy: nameingStrategy,
	})

	if err != nil {
		newLogger.Error(ctx, "failed to connect database,err: %v", err)
		os.Exit(2)
	}

	//数据库连接池配置
	sqlDB, err := db.DB()

	if err != nil {
		newLogger.Error(ctx, "failed to do db.DB(),err: %v", err)
		os.Exit(2)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func main() {
	//命令解析
	flag.Parse()

	initDb()

	// 迁移 schema
	db.AutoMigrate(&User{})

	tblname := getTableName(&User{})
	db.Logger.Info(ctx, "the table name is %s", tblname)

	_truncate()

	_insert(10)
	_selectall()

	_deleteall()
	_selectall()

}

//获取对象对应的表名
func getTableName(value interface{}) string {
	err := db.Statement.Parse(value)
	if err == nil {
		return db.Statement.Table
	}

	db.Logger.Error(ctx, "fail to get table name, model:%v, err:%v", value, err)
	return ""
}

//清空表
func _truncate() {
	db.Exec(fmt.Sprintf("truncate table %s", getTableName(&User{})))
}

//删除所有记录
func _deleteall() {
	db.Unscoped().Where("1 = 1").Delete(&User{})
	//db.Exec(fmt.Sprintf("delete from %s", getTableName(&User{})))
}

//插入测试数据
func _insert(num int) {
	kv := make([]map[string]interface{}, 0, num)
	for i := 0; i < num; i++ {
		kv = append(kv, map[string]interface{}{
			"Name": "name_" + strconv.Itoa(i),
			"Age":  i + 10,
		})
	}

	db.Model(&User{}).Select("name", "age").Create(kv)
}

//查询所有
func _selectall() {
	var users []User
	result := db.Select("id", "name", "age", "created_at").Find(&users)
	db.Logger.Info(ctx, "RowsAffected:%d, err:%v", result.RowsAffected, result.Error)

	for _, user := range users {
		fmt.Printf("Id=%d\tAge=%d\tName=%s\t\n", user.ID, user.Age, user.Name)
	}
}
