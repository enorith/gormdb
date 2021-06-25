package gormdb_test

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/enorith/gormdb"
	"github.com/enorith/supports/carbon"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var testDSN = "host=localhost user=root password=root dbname=test port=13306 sslmode=disable TimeZone=Asia/Shanghai"

type Props struct {
	Age int    `json:"age"`
	Tag string `json:"tag"`
}

func (p *Props) Scan(v interface{}) error {
	switch data := v.(type) {
	case []byte:
		return json.Unmarshal(data, p)
	case string:
		return json.Unmarshal([]byte(data), p)
	}

	return nil
}

func (p Props) Value() (driver.Value, error) {
	return json.Marshal(p)
}

type User struct {
	ID        int            `gorm:"column:id;primaryKey"`
	Name      string         `gorm:"column:name"`
	Props     Props          `gorm:"column:props"`
	CreatedAt *carbon.Carbon `gorm:"column:created_at"`
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

	db.Find(&users)
	for _, user := range users {
		fmt.Println(user.CreatedAt.GetDateTimeString(), user.Props, user.Name)
	}
}
