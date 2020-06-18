package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

func OpenMssql(addr string, user string, password string, db string) (*gorm.DB, error) {
	mysql, err := gorm.Open("mssql", fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", user, password, addr, db))
	if err != nil {
		return nil, err
	}
	mysql.LogMode(true)
	mysql.DB().SetMaxIdleConns(10)
	mysql.DB().SetMaxOpenConns(50)
	return mysql, nil
}
