package database

import "gorm.io/gorm"

type DB struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

func (d *DB) Init() {

}
