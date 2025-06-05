package db

import (
	"database/sql"

	"github.com/Zkeai/DDPay/common/database"
)

type DB struct {
	db *database.DB
}

func NewDB(conf *database.Config) *DB {
	return &DB{
		db: database.NewDB(conf),
	}
}

// GetDB 返回底层数据库连接
func (d *DB) GetDB() *sql.DB {
	return d.db.DB
}
