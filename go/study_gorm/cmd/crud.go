package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"study_gorm/config"

	"gitee.com/jianzhuliu/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	//初始化 db 命令行参数
	conf.FlagInit(conf.Fdb)
}

type User struct {
	gorm.Model
	Name string
	Age  uint8 `gorm:"default:10"`
}

//定义钩子
func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	return
}

func (u *User) BeforeUpdate(db *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

func main() {
	//命令解析
	flag.Parse()

	//指定数据库名
	conf.V_db_dbname = config.C_database

	dsn := conf.FlagDbDsn()
	fmt.Println(dsn)
	//return

	//自定义日志处理对象
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond, //慢查询
		LogLevel:      logger.Info,
		Colorful:      false, // 禁用彩色打印
	})

	//建立链接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//全局配置
		DryRun: true,
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	//数据库连接池配置
	sqlDB, err := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 迁移 schema
	//db.AutoMigrate(&User{})

	fmt.Println("========================================insert begin====================")
	_insert(db)

	fmt.Println("========================================select begin====================")
	_select(db)

	fmt.Println("========================================update begin====================")
	_update(db)

	fmt.Println("========================================delete begin====================")
	_delete(db)
}

func _insert(db *gorm.DB) {
	//创建记录
	user := User{Name: "Jinzhu", Age: 18}

	result := db.Create(&user) // 通过数据的指针来创建
	fmt.Printf("ID:%d, RowsAffected:%d, Error:%v \n", user.ID, result.RowsAffected, result.Error)

	db.Select("Name", "Age", "CreatedAt").Create(&user)

	//批量插入
	var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.Select("Name", "CreatedAt").Create(&users)

	for _, user := range users {
		fmt.Println("user.ID: ", user.ID) // 1,2,3
	}

	//通过 map 批量创建
	db.Model(&User{}).Create(map[string]interface{}{
		"Name": "jinzhu", "Age": 18,
	})

	// batch insert from `[]map[string]interface{}{}`
	db.Model(&User{}).Create([]map[string]interface{}{
		{"Name": "jinzhu_1", "Age": 18},
		{"Name": "jinzhu_2", "Age": 20},
	})
}

//查询
func _select(db *gorm.DB) {
	var user User
	var users []User

	// 获取第一条记录（主键升序）
	db.First(&user)

	// 获取一条记录，没有指定排序字段
	db.Take(&user)

	// 获取最后一条记录（主键降序）
	db.Last(&user)

	result := db.First(&user)
	//result.RowsAffected // 返回找到的记录数
	//result.Error        // returns error
	fmt.Printf("ID:%d, RowsAffected:%d, Error:%v \n", user.ID, result.RowsAffected, result.Error)

	// 检查 ErrRecordNotFound 错误
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("not found")
	}

	//根据主键查询
	db.First(&user, 10)
	db.First(&user, "10")
	db.Find(&users, []int{1, 2, 3})

	// 获取全部记录
	_ = db.Find(&users)
	// SELECT * FROM users;
}

//更新
func _update(db *gorm.DB) {
	var user User

	//db.First(&user)
	user.ID = 1
	user.Name = "jinzhu 2"
	user.Age = 100
	db.Save(&user)

	// 条件更新
	db.Model(&User{}).Where("age = ?", 10).Update("name", "hello")

	// User 的 ID 是 `111`
	db.Model(&user).Update("name", "hello")

	// 根据 `map` 更新属性
	db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18})

	//更新选定列
	// Select 和 Map

	//批量更新
	// 根据 struct 更新
	db.Model(&User{}).Where("id = ?", "111").Updates(User{Name: "hello", Age: 18})

	// 根据 map 更新
	db.Table("users").Where("id IN ?", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})

	//全部更新
	db.Model(&User{}).Where("1 = 1").Update("name", "jinzhu")
	db.Exec("UPDATE users SET name = ?", "jinzhu")
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&User{}).Update("name", "jinzhu")
}

//删除
func _delete(db *gorm.DB) {
	var user User

	// 带额外条件的删除
	db.Where("name = ?", "jinzhu").Delete(&user)

	db.Delete(&User{}, 10)
	db.Delete(&User{}, "10")

	db.Delete(&user, []int{1, 2, 3})

	//全局删除
	db.Where("1 = 1").Delete(&User{})

	db.Exec("DELETE FROM user")

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&User{})

}
