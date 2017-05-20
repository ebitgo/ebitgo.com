package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_MYSQL_DRIVER    = "mysql"
	_SQLITE_DRIVER   = "sqlite3"
	_POSTGRES_DRIVER = "postgres"
	_ALIASNAME       = "ledgercn"
)

type DatabaseManager struct {
	DbEngine   *xorm.Engine
	operations map[string]dbmodels.DBOperationInterface
}

var DBManagerInst *DatabaseManager

func (this *DatabaseManager) InitDB() {
	logmodels.Logger.Info(logmodels.SPrintInfo("DatabaseManager : InitDB", "Init database begin"))

	if DatabaseConfig == nil {
		return
	}

	this.initEngine(DatabaseConfig)

	// 注册Orm的数据库表
	dbmodels.OrmRegModels(this.DbEngine)

	this.initOperation()

	logmodels.Logger.Info(logmodels.SPrintInfo("DatabaseManager : InitDB", "Init database success"))
}

func (this *DatabaseManager) initEngine(config *DatabaseInfo) {
	// 读取配置文件中的数据库类型
	dbType := beego.AppConfig.String("dbtype")

	this.DbEngine = nil
	switch dbType {
	case _MYSQL_DRIVER:
		this.DbEngine = this.getMySqlEngine(DatabaseConfig)
	case _SQLITE_DRIVER:
		this.DbEngine = this.getSqliteEngine(DatabaseConfig)
	case _POSTGRES_DRIVER:
		this.DbEngine = this.getPostgresEngine(DatabaseConfig)
	}
	if this.DbEngine == nil {
		logmodels.Logger.Error(logmodels.SPrintError("DatabaseManager : initEngine", "Undefined db type = [%s]", dbType))
		panic(1)
	}

	isDebug := beego.AppConfig.String("runmode")
	if isDebug == "dev" {
		this.DbEngine.ShowDebug = true
		this.DbEngine.ShowInfo = true
		this.DbEngine.ShowSQL = true
	}
	this.DbEngine.ShowErr = true
	this.DbEngine.ShowWarn = true
}

func (this *DatabaseManager) getMySqlEngine(config *DatabaseInfo) *xorm.Engine {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", //Asia%2FShanghai
		config.UserName, config.Password, config.Host, config.Port, _ALIASNAME)
	ret, err := xorm.NewEngine(_MYSQL_DRIVER, dataSourceName)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("DatabaseManager : getMySqlEngine", "Create MySql has error! \r\n\t %v", err))
		return nil
	}
	err = ret.Ping()
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("DatabaseManager : getMySqlEngine", "Create MySql Ping error! \r\n\t %v", err))
		return nil
	}
	return ret
}

func (this *DatabaseManager) getSqliteEngine(config *DatabaseInfo) *xorm.Engine {
	dataSourceName := fmt.Sprintf("./%s.db?charset=utf8&loc=Asia%2FShanghai", _ALIASNAME)
	ret, err := xorm.NewEngine(_SQLITE_DRIVER, dataSourceName)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("DatabaseManager : getSqliteEngine", "Create Sqlite has error! \r\n\t %v", err))
		return nil
	}
	return ret
}

func (this *DatabaseManager) getPostgresEngine(config *DatabaseInfo) *xorm.Engine {
	dataSourceName := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=verify-full",
		_ALIASNAME, config.UserName, config.Password, config.Host, config.Port)
	ret, err := xorm.NewEngine(_POSTGRES_DRIVER, dataSourceName)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("DatabaseManager : getPostgresEngine", "Create Postgres has error! \r\n\t %v", err))
		return nil
	}
	return ret
}

func (this *DatabaseManager) initOperation() {
	if this.operations == nil {
		this.operations = make(map[string]dbmodels.DBOperationInterface)
	}

	// 注册 UserIdT 操作
	uido := &dbmodels.UserIdTableOperation{}
	uido.Init(this.DbEngine)
	this.operations[uido.GetKey()] = uido

	// 注册 UserDbT 操作
	uto := &dbmodels.UserTableOperation{}
	uto.Init(this.DbEngine)
	this.operations[uto.GetKey()] = uto

	// 注册 WalletDbT 操作
	wto := &dbmodels.WalletTableOperation{}
	wto.Init(this.DbEngine)
	this.operations[wto.GetKey()] = wto

	// 注册 FriendlyShipDbT 操作
	fso := &dbmodels.FriendlyShipTOperation{}
	fso.Init(this.DbEngine)
	this.operations[fso.GetKey()] = fso
}

func (this *DatabaseManager) Quary_Check_Email_Exist(email string) (uint64, error) {
	uidtable := &dbmodels.UserIdT{
		Email: email,
	}
	uido, _ := this.operations[dbmodels.DB_USER_ID_OPERATION]
	err := uido.Quary(dbmodels.QT_CHECK_RECORD, uidtable)
	if err != nil {
		return 0, err
	}
	return uidtable.Id, nil
}

func (this *DatabaseManager) Quary_UserInfo(uid uint64) (*dbmodels.UserDbT, error) {
	info := &dbmodels.UserDbT{
		Id: uid,
	}
	uido, _ := this.operations[dbmodels.DB_USER_INFO_OPERATION]
	err := uido.Quary(dbmodels.QT_GET_RECORD, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (this *DatabaseManager) Add_RegUser(email string, addUser *dbmodels.UserDbT) error {
	uidtable := &dbmodels.UserIdT{
		Email: email,
	}
	uido, _ := this.operations[dbmodels.DB_USER_ID_OPERATION]
	err := uido.Quary(dbmodels.QT_ADD_RECORD, uidtable)
	if err != nil {
		return err
	}

	addUser.Id = uidtable.Id
	ruo, _ := this.operations[dbmodels.DB_USER_INFO_OPERATION]
	err = ruo.Quary(dbmodels.QT_ADD_RECORD, addUser)
	return err
}

func (this *DatabaseManager) Update_UserInfo(uinfo *dbmodels.UserDbT) error {
	ruo, _ := this.operations[dbmodels.DB_USER_INFO_OPERATION]
	updateInfo := dbmodels.UserDbT{}
	updateInfo = *uinfo
	err := ruo.Quary(dbmodels.QT_UPDATE_RECORD, &updateInfo)
	return err
}

func (this *DatabaseManager) Quary_Check_Wallet_Exist(nickName string) (uint64, error) {
	wdbt := &dbmodels.WalletDbT{
		NickName: nickName,
	}
	wido, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	err := wido.Quary(dbmodels.QT_CHECK_RECORD, wdbt)
	if err != nil {
		return 0, err
	}
	return wdbt.Id, nil
}

func (this *DatabaseManager) Quary_WalletInfo(uid uint64) (*dbmodels.WalletDbT, error) {
	info := &dbmodels.WalletDbT{
		Id: uid,
	}
	wdbt, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	err := wdbt.Quary(dbmodels.QT_GET_RECORD, info)
	return info, err
}

func (this *DatabaseManager) Add_Wallet(w *dbmodels.WalletDbT) error {
	wdbt, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	err := wdbt.Quary(dbmodels.QT_ADD_RECORD, w)
	return err
}

func (this *DatabaseManager) Quary_ALL_Wallet(uid uint64) (map[string][]*dbmodels.WalletDbT, error) {
	wdbt, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	ret := make(map[string][]*dbmodels.WalletDbT)
	err := wdbt.Quary(dbmodels.QT_QUARY_ALL_RECORD, ret, uid)
	return ret, err
}

func (this *DatabaseManager) Update_WalletInfo(w *dbmodels.WalletDbT) error {
	wdbt, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	updateInfo := dbmodels.WalletDbT{}
	updateInfo = *w
	err := wdbt.Quary(dbmodels.QT_UPDATE_RECORD, &updateInfo)
	return err
}

func (this *DatabaseManager) Delete_Wallet(w *dbmodels.WalletDbT) error {
	wdbt, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	err := wdbt.Quary(dbmodels.QT_DELETE_RECORD, w)
	return err
}

func (this *DatabaseManager) Search_WalletInfo(addr, nickname string) (map[string][]*dbmodels.WalletDbT, error) {
	ret := make(map[string][]*dbmodels.WalletDbT)
	wido, _ := this.operations[dbmodels.DB_WALLET_INFO_OPERATION]
	err := wido.Quary(dbmodels.QT_SEARCH_RECORD, ret, addr, nickname)
	return ret, err
}

func (this *DatabaseManager) Add_Friend(fst *dbmodels.FriendlyShipDbT) error {
	fso, _ := this.operations[dbmodels.DB_FRIENDLY_SHIP_OPERATION]
	err := fso.Quary(dbmodels.QT_ADD_RECORD, fst)
	return err
}

func (this *DatabaseManager) Quary_Get_Friend(userid, walletid uint64) (*dbmodels.FriendlyShipDbT, error) {
	fst := &dbmodels.FriendlyShipDbT{
		UserId:   userid,
		WalletId: walletid,
	}
	fso, _ := this.operations[dbmodels.DB_FRIENDLY_SHIP_OPERATION]
	err := fso.Quary(dbmodels.QT_CHECK_RECORD, fst)
	if err != nil {
		return nil, err
	}
	return fst, nil
}

func (this *DatabaseManager) Quary_ALL_Friends(uid uint64) (map[string][]*dbmodels.FriendlyShipDbT, error) {

	logmodels.SPrintDebug("DatabaseManager : Quary_ALL_Friends", "Userid = %d", uid)

	fso, _ := this.operations[dbmodels.DB_FRIENDLY_SHIP_OPERATION]
	ret := make(map[string][]*dbmodels.FriendlyShipDbT)
	err := fso.Quary(dbmodels.QT_QUARY_ALL_RECORD, ret, uid)

	// logmodels.SPrintDebug("DatabaseManager : Quary_ALL_Friends", "ret = %v", ret)
	return ret, err
}

func (this *DatabaseManager) Delete_Friend(w *dbmodels.FriendlyShipDbT) error {
	fso, _ := this.operations[dbmodels.DB_FRIENDLY_SHIP_OPERATION]
	err := fso.Quary(dbmodels.QT_DELETE_RECORD, w)
	return err
}
