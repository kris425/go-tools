package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func OpenMysql(addr string, user string, password string, db string) (*gorm.DB, error) {
	mysql, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, addr, db))
	if err != nil {
		return nil, err
	}
	mysql.LogMode(true)
	mysql.DB().SetMaxIdleConns(10)
	mysql.DB().SetMaxOpenConns(50)
	return mysql, nil
}
