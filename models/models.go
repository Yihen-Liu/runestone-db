package models

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"galactic-monitor/config"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

// Setup initializes the database instance
func Setup() {
	var err error
	if config.AppConf.Database.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config.AppConf.Database.User,
			config.AppConf.Database.Password,
			config.AppConf.Database.Host,
			config.AppConf.Database.Db)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			},
		})
		InitTables()
	}

	if config.AppConf.Database.Type == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s port=%s  user=%s password=%s dbname=%s sslmode=%s",
			config.AppConf.Database.Host, config.AppConf.Database.Port, config.AppConf.Database.User, config.AppConf.Database.Password, config.AppConf.Database.Db, "disable")

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			},
			Logger: logger.Default.LogMode(logger.Error),
		})

		if err != nil {
			panic("db connect error" + err.Error())
		}
		InitTables()
	}

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return setting.DatabaseSetting.TablePrefix + defaultTableName
	// }

	// db.SingularTable(true)
	// db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	// db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	// db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	sqlDB, err := db.DB()
	if err != nil {
		panic("db connect error")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	// defer db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
// func updateTimeStampForCreateCallback(scope *gorm.Scope) {
// 	if !scope.HasError() {
// 		nowTime := time.Now().Unix()
// 		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
// 			if createTimeField.IsBlank {
// 				createTimeField.Set(nowTime)
// 			}
// 		}

// 		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
// 			if modifyTimeField.IsBlank {
// 				modifyTimeField.Set(nowTime)
// 			}
// 		}
// 	}
// }

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
// func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
// 	if _, ok := scope.Get("gorm:update_column"); !ok {
// 		scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	}
// }

// deleteCallback will set `DeletedOn` where deleting
// func deleteCallback(scope *gorm.Scope) {
// 	if !scope.HasError() {
// 		var extraOption string
// 		if str, ok := scope.Get("gorm:delete_option"); ok {
// 			extraOption = fmt.Sprint(str)
// 		}

// 		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

// 		if !scope.Search.Unscoped && hasDeletedOnField {
// 			scope.Raw(fmt.Sprintf(
// 				"UPDATE %v SET %v=%v%v%v",
// 				scope.QuotedTableName(),
// 				scope.Quote(deletedOnField.DBName),
// 				scope.AddToVars(time.Now().Unix()),
// 				addExtraSpaceIfExist(scope.CombinedConditionSql()),
// 				addExtraSpaceIfExist(extraOption),
// 			)).Exec()
// 		} else {
// 			scope.Raw(fmt.Sprintf(
// 				"DELETE FROM %v%v%v",
// 				scope.QuotedTableName(),
// 				addExtraSpaceIfExist(scope.CombinedConditionSql()),
// 				addExtraSpaceIfExist(extraOption),
// 			)).Exec()
// 		}
// 	}
// }

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
