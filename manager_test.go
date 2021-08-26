package gormdb_test

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/enorith/gormdb"
	"github.com/enorith/supports/carbon"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDSN = "root:root@tcp(127.0.0.1:3306)/test"

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
		return gorm.Open(mysql.Open(testDSN))
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
