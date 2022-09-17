package logic

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	// 随机数种子
	rand.Seed(time.Now().Unix())

	// 连接db
	var err error
	db, err = gorm.Open("mysql", "root:1023564552tbd@tcp(172.16.0.8:3306)/account?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	db.SingularTable(true)
}
