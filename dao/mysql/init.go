package mysql

import (
	"fmt"
	"github.com/132982317/profstik/pkg/utils/viper"
	"github.com/132982317/profstik/pkg/utils/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	DB      *gorm.DB
	config  = viper.Init("mysql")
	zlogger = zap.InitLogger()
)

// todo 池化mysql 连接池，sync,Pool ?

func getInfo(Dbr string) string {
	username := config.Viper.GetString(fmt.Sprintf("%s.username", Dbr))
	password := config.Viper.GetString(fmt.Sprintf("%s.password", Dbr))
	host := config.Viper.GetString(fmt.Sprintf("%s.port", Dbr))
	port := config.Viper.GetInt(fmt.Sprintf("%s.port", Dbr))
	Dbname := config.Viper.GetString(fmt.Sprintf("%s.database", Dbr))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, Dbname)
	return dsn
}

func init() {
	zlogger.Info("mysql server connection start")
	dsn := getInfo("mysql.source")
	var err error
	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		ConnPool:               nil,
	}); err != nil {
		zlogger.Fatalln("mysql server connection failed")
		log.Fatal(err)
	}

	if err = DB.AutoMigrate(&User{}, &Video{}, &Comment{}, &FavoriteVideoRelation{}, &FollowRelation{}, &Message{}, &FavoriteCommentRelation{}); err != nil {
		zlogger.Fatalln(err.Error())
	}

	db, err := DB.DB()
	if err != nil {
		zlogger.Fatalln(err.Error())
	}

	db.SetMaxOpenConns(config.Viper.GetInt("config.MaxIdleCons"))
	db.SetMaxIdleConns(config.Viper.GetInt("config.MaxIdleCons"))
	db.SetConnMaxLifetime(config.Viper.GetDuration("config.ConnMaxLifetime"))
}

func GetDB() *gorm.DB {
	return DB
}
