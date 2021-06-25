package gormdb_test

import (
	"testing"

	"github.com/enorith/gormdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDSN = "host=localhost user=root password=root dbname=test port=13306 sslmode=disable TimeZone=Asia/Shanghai"

type Props struct {
	Age int    `json:"age"`
	Tag string `json:"tag"`
}

//func (p *Props) Scan(v interface{}) error {
//	switch data := v.(type) {
//	case []byte:
//		return json.Unmarshal(data, p)
//	case string:
//		return json.Unmarshal([]byte(data), p)
//	}
//
//	return nil
//}
//
//func (p Props) Value() (driver.Value, error) {
//	return json.Marshal(p)
//}

type User struct {
	ID    int    `gorm:"column:id;primaryKey"`
	Name  string `gorm:"column:name"`
	Props Props  `gorm:"column:props"`
}

func Test_ManagerRegister(t *testing.T) {
	m := gormdb.NewManager()

	m.RegisterDefault(func() (*gorm.DB, error) {
		return gorm.Open(postgres.Open(testDSN))
	})

	db, err := m.GetConnection()
	if err != nil {
		t.Error(err)
	}

	var users []User

	t.Error(db.Find(&users).Error)
	t.Log(users)
}
