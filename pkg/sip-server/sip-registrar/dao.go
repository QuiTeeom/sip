package sip_registrar

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

func NewDao() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(ContactInfo{})

	return db
}

type ContactInfo struct {
	gorm.Model
	Address     string
	CallId      string
	ExpireDate  time.Time
	Cseq        int
	Contact     string
	GlobalRoute string
}
