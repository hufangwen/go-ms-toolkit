/*
@Time : 25/1/2021 公元 09:43
@Author : philiphu
@File : builder_test
@Software: GoLand
*/
package mysql



import (
	db_config "github.com/hufangwen/go-ms-toolkit/db-config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDBUtil(t *testing.T) {
	conf := db_config.NewDbConfig()
	conf.DbName = "hahaha_test"

	assert.NotNil(t, MakeDBUtil(conf))
}

func TestMakeDB(t *testing.T) {
	conf := db_config.NewDbConfig()
	conf.DbName = "hahaha_test"

	utilDB := MakeDBUtil(conf)
	assert.NotNil(t, utilDB)

	utilDB.CreateDB()

	db := MakeDB(conf)
	assert.NotNil(t, db)

	utilDB.DropDB()
}