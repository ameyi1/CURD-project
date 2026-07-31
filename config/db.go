package config

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	// 注意：导入你自己的 model 包
	"gin-study/model"
)

// 大写DB：全局数据库实例，外面文件能调用
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {

	// 连接sqlite
	dbName := GlobalConfig.Mysql.Dbname
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		panic("数据库连接失败")
	}

	// 自动根据模型建表
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Article{})
	db.Exec("CREATE INDEX idx_articles_user_id ON articles(user_id);")
	db.Exec("CREATE INDEX idx_articles_title ON articles(title);")

	// 赋值给全局DB，给其他地方用
	DB = db
}
