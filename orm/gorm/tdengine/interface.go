/*
@Time : 25/1/2021 公元 10:24
@Author : philiphu
@File : interface
@Software: GoLand
*/
package tdengine

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/taosdata/driver-go/taosSql"
)

type DBUtil interface {
	CreateDB()
	DropDB()
	GetUtilDB() *gorm.DB
}

type DB interface {
	GetDB() *gorm.DB
	ClearAllData()
	Create(value interface{}) error
	BulkInsert(value interface{}) error
}

type TdengineDb interface {
	CreateDB() error
	DropDB() error
	GetDB() *sql.DB
	Create(value interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryMap(query string, args ...interface{})([]map[string]interface{},error)
	QueryStruct(obj interface{},query string, args ...interface{}) error
}

