package tdengine

import (
	"database/sql"
	"fmt"
	db_config "github.com/hufangwen/go-ms-toolkit/db-config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_gormTDengine_QueryMap(t *testing.T) {
     dbConn,err := sql.Open("taosSql", fmt.Sprintf("%s:%s@/tcp(%s:%s)/", "root", "taosdata", "10.10.10.159", "6030"))
     assert.NoError(t, err)
     err = dbConn.Ping()
	 assert.NoError(t, err)
	db := InitConnect(&db_config.DbConfig{
		Username: "root",
		Password:"taosdata",
		Host:"10.10.10.159",
		Port:"6030",
		MaxIdleConns:100,
		MaxOpenConns:10,
		LogMode:true,
	})
	queryMap,err := db.QueryMap("SELECT ip as node_ip,mem FROM  rsmhostmetric_org001.host  WHERE ts > now -10d AND ip ='10.10.10.160'")
	assert.NoError(t, err)
	fmt.Println(queryMap)
}
type TestStruct struct{
	NodeIp string `json:"node_ip"`
	Mem float64 `json:"mem"`
}

func Test_gormTDengine_QueryStruct(t *testing.T) {
	dbConn,err := sql.Open("taosSql", fmt.Sprintf("%s:%s@/tcp(%s:%s)/", "root", "taosdata", "10.10.10.159", "6030"))
	assert.NoError(t, err)
	err = dbConn.Ping()
	assert.NoError(t, err)
	db := InitConnect(&db_config.DbConfig{
		Username: "root",
		Password:"taosdata",
		Host:"10.10.10.159",
		Port:"6030",
		MaxIdleConns:100,
		MaxOpenConns:10,
		LogMode:true,
	})
	test := []TestStruct{}
	err = db.QueryStruct(&test,"SELECT mem as mem,ip as node_ip FROM  rsmhostmetric_org001.host  WHERE ts > now -10d AND ip ='10.10.10.160'")
	if  err != nil{
		fmt.Println(err)
	}
	fmt.Println(len(test))
}