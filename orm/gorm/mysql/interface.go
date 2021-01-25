/*
@Time : 25/1/2021 公元 09:43
@Author : philiphu
@File : interface
@Software: GoLand
*/
package mysql


import "github.com/jinzhu/gorm"

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