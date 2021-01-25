/*
@Time : 25/1/2021 公元 10:24
@Author : philiphu
@File : interface
@Software: GoLand
*/
package tdengine

import (
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
}
