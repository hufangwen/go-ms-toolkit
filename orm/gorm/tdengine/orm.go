/*
@Time : 25/1/2021 公元 10:21
@Author : philiphu
@File : orm
@Software: GoLand
*/
package tdengine

import (
	"fmt"
	"github.com/hufangwen/go-ms-toolkit/log"
	"github.com/hufangwen/go-ms-toolkit/qyenv"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/taosdata/driver-go/taosSql"
	"go.uber.org/zap"

	db_config "github.com/hufangwen/go-ms-toolkit/db-config"
)

type gormTDengine struct {
	dbConfig *db_config.DbConfig
	db       *gorm.DB
	utilDB   *gorm.DB
}

func (gm *gormTDengine) CreateDB() {
	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + gm.dbConfig.DbName + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"

	err := gm.utilDB.Exec(createDbSQL).Error
	if err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库创建成功")
}

func (gm *gormTDengine) DropDB() {
	dropDbSQL := "DROP DATABASE IF EXISTS " + gm.dbConfig.DbName + ";"

	err := gm.utilDB.Exec(dropDbSQL).Error
	if err != nil {
		fmt.Println("删除失败：" + err.Error() + " sql:" + dropDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库删除成功")
}

func (gm *gormTDengine) GetDB() *gorm.DB {
	return gm.db
}

func (gm *gormTDengine) GetUtilDB() *gorm.DB {
	log.QyLogger.Info("init db connection: ", zap.String("db_host", gm.dbConfig.Host),
		zap.String("db_name", gm.dbConfig.DbName), zap.String("user", gm.dbConfig.Username))

	openedDb, err := gorm.Open("taosSql", fmt.Sprintf("%s:%s@/tcp(%s:%s)/%s?interpolateParams=true", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName))
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	openedDb.DB().SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	openedDb.DB().SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	openedDb.DB().SetConnMaxLifetime(time.Hour * 1)
	// 如果不是生产数据库则打开详细日志
	// if !strings.Contains(dbConfig.DbName, "prod") {
	if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
		openedDb.LogMode(true)
	}

	return openedDb
}

func (gm *gormTDengine) ClearAllData() {
	if qyenv.IsUnitTestEnv() && strings.Contains(gm.dbConfig.DbName, "test") {
		tmpDb := gm.db
		if tmpDb == nil {
			panic("尚未初始化数据库, 清空数据库失败")
		}
		if rs, err := tmpDb.Raw("show tables;").Rows(); err == nil {
			var tName string
			for rs.Next() {
				if err := rs.Scan(&tName); err != nil || tName == "" {
					fmt.Println("表名获取失败", err, tName)
					panic("表名获取失败")
				}
				if err := tmpDb.Exec(fmt.Sprintf("delete from %s", tName)).Error; err != nil {
					panic("清空表数据失败:" + err.Error())
				}
			}
		} else {
			panic("表名列表获取失败：" + err.Error())
		}
	} else {
		panic("非法操作！在非测试环境下调用了清空所有数据的方法")
	}
}

// TODO 支持多种并发写入
func (gm *gormTDengine) Create(value interface{}) error {
	return gm.GetDB().Exec(value.(string)).Error
}

func (gm *gormTDengine) BulkInsert(value interface{}) error {
	//先判断是不是数组 如果不是
	if kind := reflect.TypeOf(value).Kind(); kind != reflect.Slice || kind != reflect.Array {
		return gm.GetDB().Create(value).Error
	}
	return nil
}

func newGormTDengine(dbConfig *db_config.DbConfig, forUtil bool) *gormTDengine {
	gm := &gormTDengine{dbConfig: dbConfig}

	if forUtil {
		gm.initCdDb()
		return gm
	}

	// init db
	gm.initGormDB()

	return gm
}
// 该连接并没有指定特定的db
func InitConnect(dbConfig *db_config.DbConfig,) *gormTDengine {
	gm := &gormTDengine{dbConfig: dbConfig}
	gm.tdEngineConnect()
	return gm
}

func (gm *gormTDengine) tdEngineConnect()  {
	log.QyLogger.Info("init db connection: ", zap.String("db_host", gm.dbConfig.Host),
		zap.String("db_name", gm.dbConfig.DbName), zap.String("user", gm.dbConfig.Username))

	openedDb, err := gorm.Open("taosSql", fmt.Sprintf("%s:%s@/tcp(%s:%s)?interpolateParams=true", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port))
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	openedDb.DB().SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	openedDb.DB().SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	openedDb.DB().SetConnMaxLifetime(time.Hour * 1)
		openedDb.LogMode(gm.dbConfig.LogMode)

	gm.db = openedDb
}




// 这里找时间优化一下下
func (gm *gormTDengine) initGormDB() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	log.QyLogger.Info("init db connection: ", zap.String("db_host", gm.dbConfig.Host),
		zap.String("db_name", gm.dbConfig.DbName), zap.String("user", gm.dbConfig.Username))

	openedDb, err := gorm.Open("taosSql", fmt.Sprintf("%s:%s@/tcp(%s:%s)/%s?interpolateParams=true", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName))
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	openedDb.DB().SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	openedDb.DB().SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	openedDb.DB().SetConnMaxLifetime(time.Hour * 1)
	// 如果不是生产数据库则打开详细日志
	// if !strings.Contains(dbConfig.DbName, "prod") {
	if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
		openedDb.LogMode(true)
	}

	gm.db = openedDb
}

func (gm *gormTDengine) initCdDb() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	cStr := fmt.Sprintf("%s:%s@/tcp(%s:%s)/%s?interpolateParams=true", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName)
	openedDb, err := gorm.Open("taosSql", cStr)
	if err != nil {
		fmt.Println(cStr)
		panic("连接数据库出错:" + err.Error())
	}

	gm.utilDB = openedDb
}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
