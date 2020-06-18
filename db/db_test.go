package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestOpenMysql(t *testing.T) {
	d, err := OpenMysql("localhost", "root", "123456", "test")
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	defer d.Close()

	d.AutoMigrate(&Product{})

	d.Create(&Product{Code: "111", Price: 1000})

	var product Product
	d.First(&product, 1)
	fmt.Println(product)

	d.Model(&product).Update("Price", 2000)
	d.First(&product, "code=?", "111")
	fmt.Println(product)

	d.Delete(&product)
}

func TestOpenMssql(t *testing.T) {
	d, err := OpenMssql("192.168.0.200:1433", "sa", "mssql@123", "admin_platform")
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}
	defer d.Close()
}
