/*
@Time : 25/1/2021 公元 10:21
@Author : philiphu
@File : orm
@Software: GoLand
*/
package tdengine

import (
	"database/sql"
	"fmt"
	"github.com/hufangwen/go-ms-toolkit/log"
	"strings"
	"time"

	_ "github.com/taosdata/driver-go/taosSql"
	"go.uber.org/zap"

	db_config "github.com/hufangwen/go-ms-toolkit/db-config"
)

type gormTDengine struct {
	dbConfig *db_config.DbConfig
	 *sql.DB
	logModel bool
}

func (gm *gormTDengine) CreateDB() {
	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + gm.dbConfig.DbName + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"

	_,err := gm.Exec(createDbSQL)
	if err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库创建成功")
}

func (gm *gormTDengine) DropDB() {
	dropDbSQL := "DROP DATABASE IF EXISTS " + gm.dbConfig.DbName + ";"

	_,err := gm.Exec(dropDbSQL)
	if err != nil {
		fmt.Println("删除失败：" + err.Error() + " sql:" + dropDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库删除成功")
}

func (gm *gormTDengine) GetDB() *gormTDengine {
	return gm
}


// TODO 支持多种并发写入
func (gm *gormTDengine) Create(value interface{}) error {
	_,err:= gm.GetDB().Exec(value.(string))
	return err
}


// 该连接并没有指定特定的db
func InitConnect(dbConfig *db_config.DbConfig,) *gormTDengine {
	gm := &gormTDengine{dbConfig: dbConfig}
	gm.tdEngineConnect()
	return gm
}

func (gm *gormTDengine) tdEngineConnect() {
	log.QyLogger.Info("init db connection: ", zap.String("db_host", gm.dbConfig.Host),
		zap.String("db_name", gm.dbConfig.DbName), zap.String("user", gm.dbConfig.Username))

	openedDb, err := sql.Open("taosSql", fmt.Sprintf("%s:%s@/tcp(%s:%s)/", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port))
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	openedDb.SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	openedDb.SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	openedDb.SetConnMaxLifetime(time.Hour * 1)
	gm.logModel = gm.dbConfig.LogMode
	gm.DB = openedDb
}

func (gm *gormTDengine)Query(query string, args ...interface{}) (*sql.Rows, error) {
	if gm.logModel{
		fmt.Sprintf(strings.ReplaceAll(query,"?","%v"),args)
	}
	return gm.DB.Query(query,args)
}

func (gm *gormTDengine)QueryRow(query string, args ...interface{}) *sql.Row {
	if gm.logModel{
		fmt.Sprintf(strings.ReplaceAll(query,"?","%v"),args)
	}
	return gm.DB.QueryRow(query,args)
}




func (gm *gormTDengine)Exec(query string, args ...interface{}) (sql.Result, error) {
	if gm.logModel{
		fmt.Sprintf(strings.ReplaceAll(query,"?","%v"),args)
	}
	return gm.DB.Exec(query,args)
}

func (gm *gormTDengine)QueryMap(query string, args ...interface{})(map[string]interface{},error) {
	var reset map[string]interface{}
	if gm.logModel{
		fmt.Sprintf(strings.ReplaceAll(query,"?","%v"),args)
	}
	rows,err := gm.DB.Query(query,args)
	if err != nil{
		fmt.Sprintf("query map error = %v",err)
		return nil, err
	}
	unm,err := rows.Columns()
	for k,_ :=  range unm{
		fmt.Sprintf(k)
	}

}
