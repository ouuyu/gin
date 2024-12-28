package model

import (
	"main/common"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func createRootAccountIfNeed() error {
	var user User

	if err := DB.First(&user).Error; err != nil {
		common.SysLog("no user exists, create a root user for you: username is root, password is 123456")
		hashedPassword, err := common.Password2Hash("123456")
		if err != nil {
			return err
		}
		rootUser := User{
			Username: "root",
			Password: hashedPassword,
			Role:     common.RoleRootUser,
			Status:   common.UserStatusEnabled,
		}
		DB.Create(&rootUser)
	}
	return nil
}

func CountTable(tableName string) (num int64) {
	DB.Table(tableName).Count(&num)
	return
}

func InitDB() (err error) {
	var db *gorm.DB
	if os.Getenv("SQL_DSN") != "" {
		// Use MySQL
		db, err = gorm.Open(mysql.Open(os.Getenv("SQL_DSN")), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	} else {
		// Use SQLite
		db, err = gorm.Open(sqlite.Open(common.SQLitePath), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
		common.SysLog("SQL_DSN not set, using SQLite as database")
	}
	if err == nil {
		DB = db
		err = db.AutoMigrate(&User{}, &Config{})
		if err != nil {
			return err
		}
		err = createRootAccountIfNeed()
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Group{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&Order{})
		if err != nil {
			return err
		}
	} else {
		common.FatalLog(err)
	}
	return err
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
