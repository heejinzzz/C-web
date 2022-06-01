package database

import (
	"github.com/heejinzzz/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var logFilename = "./webLog.log"

// 导入database包时自动执行init函数连接数据库
func init() {
	dsn := "root:we359889ig@tcp(127.0.0.1:3306)/C-web?charset=utf8mb3&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 连接数据库失败，记录日志
		errorMsg := "{C-webLog}\tFailed to connect database MySQL.  Error:" + err.Error()
		logger.RecordError(logFilename, errorMsg)
		panic(err)
	}
}

// User 一个用户对象
type User struct {
	Username string
	Password string
}

// TableName 将结构体与数据库中对应的表名绑定
func (User) TableName() string {
	return "user"
}

// Admin 一个管理员对象
type Admin struct {
	Adminname string
	Password  string
}

// TableName 将结构体与数据库中对应的表名绑定
func (Admin) TableName() string {
	return "admin"
}

// Article 一篇文章对象
type Article struct {
	Id          int
	Articlename string
	Author      string
	Addtime     string
	Content     string
}

// TableName 将结构体与数据库中对应的表名绑定
func (Article) TableName() string {
	return "article"
}

// Movie 一个电影对象
type Movie struct {
	Id        int
	Moviename string
	Addtime   string
	Coverpath string
	Path      string
}

// TableName 将结构体与数据库中对应的表名绑定
func (Movie) TableName() string {
	return "movie"
}
