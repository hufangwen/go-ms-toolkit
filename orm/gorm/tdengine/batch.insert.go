/*
@Time : 25/1/2021 公元 10:24
@Author : philiphu
@File : batch.insert
@Software: GoLand
*/
package tdengine

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	_ "github.com/taosdata/driver-go/taosSql"
	"github.com/jinzhu/gorm"
)

var flowIgnoreFields = []string{"Wall", "Ext", "Loc", "wall", "ext", "loc", "ID", "Id", "DeletedAt"}

// data must be slice
func DoBatchInsert(data interface{}, db *gorm.DB) error {
	return db.Exec(data.(string)).Error
}


// TODO 批量插入组装sql  想想 时间戳必须为第一位 可以尝试拼接
type BatchInsertSql struct {
	TableName string
	// 保存了属性信息，可以在组装sql时根据属性做不同的操作
	Fields    []reflect.StructField
	InsertSql string

	createdAt string
}


func (b *BatchInsertSql) ResultSql() string {
	return b.InsertSql
}

// 获取插入sql的字段部分
func getInsertFieldStr(rt reflect.Type, ignoreFs []string) (fieldNames []reflect.StructField, fStr string) {
	fStr = "("
	EnumAnObjFieldNames(rt, func(f reflect.StructField) {
		if f.Tag.Get("sql") != "-" {
			tmpName := f.Name
			// 如果没有ignore则纳入要用的字段中
			if !StrSliceContains(ignoreFs, tmpName) {
				fieldNames = append(fieldNames, f)
				fStr += gorm.ToDBName(tmpName) + ","
				// 保证ID只被加一次
			} else if f.Name == "ID" && f.Type.Kind() == reflect.String {
				fieldNames = append(fieldNames, f)
				fStr += gorm.ToDBName(tmpName) + ","
			}
		}
	})
	if fStr != "" {
		fStr = fStr[:(len(fStr) - 1)]
	}
	fStr += ")"
	return
}

// 获取插入sql的值部分
func (b *BatchInsertSql) getObjValuesForSql(rv reflect.Value, fields []reflect.StructField) (result string) {
	result = "("
	for _, f := range fields {
		// logger.Debug(f.Type.Kind().String())
		// logger.Debug(f.Type.Name())
		// fmt.Println("!!!!!!!", f.Name, f.Type.String(), f.Type.Kind())
		// 尚未实现根据类型做适配，因此必须都是string
		if f.Type.Kind() == reflect.Struct && strings.Contains(f.Type.String(), "Time") {
			if f.Name == "CreatedAt" || f.Name == "UpdatedAt" {
				result += "'" + b.createdAt + "',"
			} else {
				result += "'" + rv.FieldByName(f.Name).Interface().(time.Time).Format("2006-01-02 15:04:05") + "',"
			}
			// todo, 优化这个判断
		} else if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct && strings.Contains(f.Type.String(), "Time") {
			if f.Name == "CreatedAt" || f.Name == "UpdatedAt" {
				result += "'" + b.createdAt + "',"
			} else {
				t := rv.FieldByName(f.Name).Interface().(*time.Time)
				if t == nil {
					result += "NULL,"
				} else {
					result += "'" + t.Format("2006-01-02 15:04:05") + "',"
				}
			}
		} else if f.Type.Kind() == reflect.String {
			result += "'" + ClearData4str(rv.FieldByName(f.Name).String()) + "',"
		} else if f.Type.Kind() == reflect.Bool {
			if rv.FieldByName(f.Name).Bool() {
				result += "'1',"
			} else {
				result += "'0',"
			}
		} else if f.Type.Kind() == reflect.Map || f.Type.Kind() == reflect.Array {
			panic("not support map or array in batch insert: " + f.Name)
		} else {
			if f.Tag.Get("sql") != "-" {
				result += fmt.Sprintf("'%v',", rv.FieldByName(f.Name).Interface())
			}
		}
	}
	result = result[:(len(result) - 1)]
	result += ")"
	return
}

// 清洗特殊字符串，目前有：
// 1. 单引号转义(批量插入时报错)
func ClearData4str(str string) string {
	if strings.Contains(str, "'") {
		return strings.Replace(str, "'", " ", -1)
	}
	return str
}

// 看一个数组中是否含有某个元素
func StrSliceContains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// 迭代一个对象的所有字段名
func EnumAnObjFieldNames(rv reflect.Type, cb func(f reflect.StructField)) {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if tmpType.Kind() == reflect.Struct && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			EnumAnObjFieldNames(tmpType, cb)
		} else {
			cb(tmpF)
		}

	}
}

