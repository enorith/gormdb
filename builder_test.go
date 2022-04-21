package gormdb_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/enorith/gormdb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_Builder(t *testing.T) {
	tx, e := gormdb.DefaultManager.GetConnection()

	if e != nil {
		t.Fatal(e)
	}

	b := gormdb.NewBuilder[User](tx)

	u, _ := b.First(12)
	fmt.Println(u)

	us, e := b.Query(func(d *gorm.DB) *gorm.DB {
		return d.Where("id > ?", 1)
	}).Get()

	if e != nil {
		t.Fatal(e)
	}

	for _, u := range us {
		fmt.Println(u)
	}
}

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		},
	)

	gormdb.DefaultManager.RegisterDefault(func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(testDSN), &gorm.Config{
			Logger: newLogger,
		})
	})

}
