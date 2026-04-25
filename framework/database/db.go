package database

import (
	"encoder/domain"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	Db          *gorm.DB
	Dsn         string
	DsnTest     string
	DbType      string
	DbTypeTest  string
	Debug       bool
	AutoMigrate bool
	Env         string
}

func NewDb() *Database {
	return &Database{}
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error
	if d.Env != "test" {
		d.Db, err = gorm.Open(d.DbType, d.Dsn)
	} else {
		d.Db, err = gorm.Open(d.DbTypeTest, d.DsnTest)
	}

	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	if d.Debug {
		d.Db.Debug()
	}

	if d.AutoMigrate {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		if d.Env != "test" {
			d.Db.Model(&domain.Job{}).AddForeignKey("video_id", "videos(id)", "cascade", "cascade")
		}
	}

	return d.Db, nil
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "test"
	dbInstance.DsnTest = ":memory:"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.Debug = true
	dbInstance.AutoMigrate = true

	conn, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	return conn

}
