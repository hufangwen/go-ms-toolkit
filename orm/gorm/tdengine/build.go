/*
@Time : 25/1/2021 公元 10:20
@Author : philiphu
@File : build
@Software: GoLand
*/
package tdengine


import db_config "git.forms.io/universe/rapm/orm/tdorm/go-ms-toolkit/db-config"

func MakeDBUtil(dbConfig *db_config.DbConfig) DBUtil {
	return newGormTDengine(dbConfig, true)
}

func MakeDB(dbConfig *db_config.DbConfig) DB {
	return newGormTDengine(dbConfig, false)
}
