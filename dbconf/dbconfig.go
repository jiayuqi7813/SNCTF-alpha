package dbconf

import (

_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//定义全局数据库
var (
	Dbctf *gorm.DB
)
//定义全局数据库连接
func Inimysql()(err error){
	dsn := "root:root@tcp(127.0.0.1:3306)/snctf?charset=utf8mb4&parseTime=True&loc=Local"
	Dbctf ,err = gorm.Open("mysql",dsn)
	err = Dbctf.DB().Ping()
	return
}
func DBLink() {
	err := Inimysql()
	if err != nil{
		panic(err)
	}
}