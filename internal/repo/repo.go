package repo

import (
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/repo/db"
)

type Repo struct {
	db *db.DB
}

func NewRepo(conf *conf.Conf) *Repo {
	return &Repo{
		db: db.NewDB(conf.DB),
	}
}
