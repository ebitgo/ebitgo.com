package models

import (
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type DatabaseInfo struct {
	Host     string
	Port     string
	UserName string
	Password string
}

var DatabaseConfig = &DatabaseInfo{}

func MainInit(dbInfo *DatabaseInfo) {
	logmodels.Logger.Info(logmodels.SPrintInfo("MainInit", "================== Start =================="))
	DBManagerInst = new(DatabaseManager)
	DBManagerInst.InitDB()
}
