/*
@Time : 25/1/2021 公元 10:20
@Author : philiphu
@File : build
@Software: GoLand
*/
package tdengine


import (
	db_config "github.com/hufangwen/go-ms-toolkit/db-config"
	_ "github.com/taosdata/driver-go/taosSql"
)

func MakeTDb(dbConfig *db_config.DbConfig) TdengineDb {
	return InitConnect(dbConfig)
}
