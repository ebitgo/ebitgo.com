package dbmodels

import (
	"github.com/go-xorm/xorm"
)

const (
	DB_USER_ID_OPERATION       = "User_Id_Operation"
	DB_USER_INFO_OPERATION     = "User_Info_Operation"
	DB_WALLET_INFO_OPERATION   = "Wallet_Info_Operation"
	DB_FRIENDLY_SHIP_OPERATION = "Friendly_ship_Operation"

	QT_GET_RECORD = iota + 1
	QT_CHECK_RECORD
	QT_QUARY_ALL_RECORD
	QT_ADD_RECORD
	QT_UPDATE_RECORD
	QT_DELETE_RECORD
	QT_SEARCH_RECORD
)

type DatabaseOperationFuncs struct {
}

type DBOperationInterface interface {
	Init(e *xorm.Engine) error
	GetKey() string
	Quary(qtype int, v ...interface{}) error
}
