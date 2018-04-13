package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "gitlab.ucloudadmin.com/wu/logrus"
)

var dbx *sqlx.DB

func InitSqlxClient(connStr string) {
	var err error
	dbx, err = sqlx.Connect("mysql", connStr)
	if err != nil {
		log.WithError(err).Fatalf("[initSqlxClient] connect db:%s failed", connStr)
	}
}
