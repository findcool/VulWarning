package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/virink/vulwarning/common"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

// InitConnect -
func InitConnect(conf common.Config, debug bool) (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=30s",
		conf.MySQL.User, conf.MySQL.Pass, conf.MySQL.Host, conf.MySQL.Name, conf.MySQL.Charset,
	)
	if db, err = gorm.Open("mysql", dsn); err != nil {
		return nil, err
	}
	db.LogMode(debug)
	db.DB().SetConnMaxLifetime(100 * time.Second) // 最大连接周期，超过时间的连接就close
	db.DB().SetMaxOpenConns(100)                  // 设置最大连接数
	db.DB().SetMaxIdleConns(16)                   // 设置闲置连接数

	// 禁用默认表名的复数形式
	db.SingularTable(true)

	// 表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.MySQL.Prefix + defaultTableName
	}

	return db, nil
}

// Resp - Web Server Response
type Resp struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

// InitTable -
func InitTable() error {
	// NOTE: Drop if exists ?
	return db.CreateTable(&Warning{}).Error
}
