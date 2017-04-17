// author by @xiaoyusilen

package tile38_sample

import (
	"testing"

	r "github.com/GoRethink/gorethink.git"
	log "github.com/Sirupsen/logrus"
	"github.com/geo-go/common/service"
	"github.com/geo-go/config"
	"github.com/stretchr/testify/assert"
)

func TestNewRethinkdb(t *testing.T) {
	// 读取配置文件
	cfg := config.ParseFromFlags()

	session := service.NewRethinkdb(cfg.RethinkAddress)

	// 创建数据库表测试
	response, err := r.DB("test").TableCreate("test").RunWrite(session)
	if err != nil {
		t.Errorf("Error creating table: %s", err)
	}

	log.Infof("%d table created", response.TablesCreated)

	// 插入数据测试
	response, err = r.DB("test").Table("test").Insert(map[string]string{
		"name": "test",
	}).RunWrite(session)

	if err != nil {
		t.Errorf("Error insert data: %s", err)
	}

	log.Infof("%d row inserted", response.Inserted)

	// 查询数据测试
	cur, err := r.DB("test").Table("test").Map(map[string]interface{}{
		"id":   r.Row.Field("id"),
		"name": r.Row.Field("name"),
	}).Run(session)
	if err != nil {
		t.Errorf("Error query: %s", err)
	}

	var res map[string]interface{}
	err = cur.One(&res)
	if err != nil {
		t.Errorf("Error query: %s", err)
	}

	assert := assert.New(t)

	assert.Equal(res["name"], "test", "they should be equal")

	log.Infof("Query success")

	// 更新数据测试
	response, err = r.DB("test").Table("test").Filter(map[string]string{
		"name": "test",
	}).Update(map[string]interface{}{
		"name": "test_change",
	}).RunWrite(session)

	if err != nil {
		t.Errorf("Error update data: %s", err)
	}

	log.Infof("%d row updated", response.Replaced)

	// 删除数据测试
	response, err = r.DB("test").Table("test").Filter(map[string]string{
		"name": "test_change",
	}).Delete().RunWrite(session)

	if err != nil {
		t.Errorf("Error delete data: %s", err)
	}

	log.Infof("%d row deleted", response.Deleted)

	// 删除数据库表测试
	response, err = r.DB("test").TableDrop("test").RunWrite(session)
	if err != nil {
		t.Errorf("Error drop table: %s", err)
	}

	log.Infof("%d table droped", response.TablesDropped)

	//Sample output
	//INFO[0000] 1 table created
	//INFO[0000] 1 row inserted
	//INFO[0000] Query success
	//INFO[0000] 1 row updated
	//INFO[0000] 1 row deleted
	//INFO[0000] 1 table droped
	//PASS
	//ok      github.com/geo-go/tile38_sample 0.708s

}
