package gormdb_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/enorith/gormdb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = "root:root@tcp(127.0.0.1:3306)/test"

type TestUser struct {
	ID        int    `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	CreatedAt string `gorm:"column:created_at"`
}

func (tu TestUser) TableName() string {
	return "users"
}

func Test_Paginator(t *testing.T) {
	m := gormdb.NewManager()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		},
	)

	m.RegisterDefault(func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	})
	p := gormdb.NewPaginator(1, 15, m)

	tx, e := m.GetConnection()
	if e != nil {
		t.Fatal(e)
	}

	var us []TestUser
	r, e := p.Paginate(tx, &us)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(r["data"], r["meta"])
}
