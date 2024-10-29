package qdb

import (
	"fmt"
)

type setting struct {
	OpenLog            byte `comment:"开关 0否 1是"`
	SkipDefTransaction byte
	DBConfig           string `comment:"数据库类型|参数\n sqlite|./db/data.db&OFF\n sqlserver|用户名:密码@地址?database=数据库&encrypt=disable\n mysql|用户名:密码@tcp(127.0.0.1:3306)/数据库?charset=utf8mb4&parseTime=True&loc=Local"`
}

func loadSetting(key string) setting {
	dbName := key
	if dbName == "" {
		dbName = "data"
	}
	def := setting{
		OpenLog:            0,
		SkipDefTransaction: 1,
		DBConfig:           fmt.Sprintf("sqlite|./db/%s.db&OFF", dbName),
	}
	return def
}
