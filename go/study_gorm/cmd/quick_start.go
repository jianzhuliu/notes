package main

import (
	"flag"
	"fmt"

	"gitee.com/jianzhuliu/common/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//需要首先创建数据库
//create database gorm default character set utf8mb4 collate utf8mb4_0900_ai_ci;
var database string = "gorm"

func init() {
	//初始化 db 命令行参数
	conf.FlagInit(conf.Fdb)
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	//命令解析
	flag.Parse()

	//指定数据库名
	conf.V_db_dbname = database

	dsn := conf.FlagDbDsn()
	fmt.Println(dsn)
	//return

	//建立链接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整形主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	fmt.Println("%+v", product)

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	db.First(&product, 1)
	fmt.Println("%+v", product)

	// Delete - 删除 product
	//db.Delete(&product, 1)

}
