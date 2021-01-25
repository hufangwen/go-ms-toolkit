/*
@Time : 25/1/2021 公元 09:42
@Author : philiphu
@File : builder
@Software: GoLand
*/
package mysql

import db_config "git.forms.io/universe/rapm/orm/tdorm/go-ms-toolkit/db-config"

func MakeDBUtil(dbConfig *db_config.DbConfig) DBUtil {
	return newGormMysql(dbConfig, true)
}

func MakeDB(dbConfig *db_config.DbConfig) DB {
	return newGormMysql(dbConfig, false)
}