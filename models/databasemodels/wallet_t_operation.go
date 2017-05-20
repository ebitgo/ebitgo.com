package dbmodels

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ebitgo/ebitgo.com/models/logger"
	"github.com/go-xorm/xorm"
)

type WalletTableOperation struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

func (this *WalletTableOperation) Init(e *xorm.Engine) error {
	if e == nil {
		err := errors.New("Can not get database engine!")
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : Init", "%v", err))
		return err
	}
	this.locker = &sync.Mutex{}
	this.engine = e
	return nil
}

func (this *WalletTableOperation) GetKey() string {
	return DB_WALLET_INFO_OPERATION
}

func (this *WalletTableOperation) Quary(qtype int, v ...interface{}) (err error) {
	if v == nil || len(v) < 1 {
		err = errors.New("Quary get use id , must input WalletDbT struct pointer as parameter")
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : Quary", "%v", err))
		return
	}

	switch qtype {
	case QT_ADD_RECORD:
		rut := v[0].(*WalletDbT)
		this.locker.Lock()
		err = this.AddWalletInfo(rut)
		this.locker.Unlock()
		return
	case QT_CHECK_RECORD:
		rut := v[0].(*WalletDbT)
		this.locker.Lock()
		this.GetWalletId(rut)
		this.locker.Unlock()
		return
	case QT_GET_RECORD:
		rut := v[0].(*WalletDbT)
		this.locker.Lock()
		err = this.GetWalletInfo(rut)
		this.locker.Unlock()
		return
	case QT_UPDATE_RECORD:
		rut := v[0].(*WalletDbT)
		this.locker.Lock()
		err = this.UpdateWalletInfo(rut)
		this.locker.Unlock()
		return
	case QT_QUARY_ALL_RECORD:
		ruts := v[0].(map[string][]*WalletDbT)
		uid := v[1].(uint64)
		this.locker.Lock()
		ruts["RESULT"], err = this.QuaryAll(uid)
		this.locker.Unlock()
		return
	case QT_DELETE_RECORD:
		rut := v[0].(*WalletDbT)
		this.locker.Lock()
		err = this.DeleteWalletInfo(rut)
		this.locker.Unlock()
		return
	case QT_SEARCH_RECORD:
		ruts := v[0].(map[string][]*WalletDbT)
		addr := v[1].(string)
		nickname := v[2].(string)
		this.locker.Lock()
		ruts["RESULT"], err = this.SearchWalletInfo(addr, nickname)
		this.locker.Unlock()
		return
	}

	err = errors.New(fmt.Sprintf("Quary type is not defined(%d)", qtype))
	logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : Quary", "%v", err))
	return
}

func (this *WalletTableOperation) GetWalletInfo(wdbt *WalletDbT) error {
	if wdbt.Id == 0 && len(wdbt.NickName) == 0 {
		err := errors.New("Get wallet info Fail! (WalletDbT.Id = 0 && WalletDbT.NickName is empty)")
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : GetWalletInfo", "%v", err))
		return err
	}
	b, err := this.engine.Get(wdbt)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : GetWalletInfo", "%v", err))
		return err
	}
	if !b {
		err := errors.New(fmt.Sprintf("Get wallet info Fail! (WalletDbT.Id = %d && WalletDbT.NickName = %s is not exist)", wdbt.Id, wdbt.NickName))
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : GetWalletInfo", "%v", err))
		return err
	}
	return nil
}

func (this *WalletTableOperation) GetWalletId(wdbt *WalletDbT) uint64 {
	tmp := &WalletDbT{
		NickName: wdbt.NickName,
	}
	b, err := this.engine.Get(tmp)
	if err != nil {
		wdbt.Id = 0
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : GetWalletId", "%v", err))
		return 0
	}
	if !b {
		wdbt.Id = 0
		return 0
	}
	*wdbt = *tmp
	return tmp.Id
}

func (this *WalletTableOperation) AddWalletInfo(wdbt *WalletDbT) error {
	if wdbt.UserId == 0 {
		err := errors.New("Add wallet info Fail! (WalletDbT.UserId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : AddWalletInfo", "%v", err))
		return err
	}
	if len(wdbt.NickName) == 0 {
		err := errors.New("Add wallet info Fail! (WalletDbT.NickName is empty)")
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : AddWalletInfo", "%v", err))
		return err
	}

	counts, err := this.engine.Count(&WalletDbT{})
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : AddWalletInfo", "Quary WalletDbT count has error \r\n\t%v", err))
		return err
	}

	wdbt.Id = this.createWalletID(counts)
	wdbt.RegTime = time.Now()
	_, err = this.engine.Insert(wdbt)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : AddWalletInfo", "Insert wallet to WalletDbT has error \r\n\t%v", err))
	}
	return err
}

func (this *WalletTableOperation) UpdateWalletInfo(wdbt *WalletDbT) error {
	sql := "update `wallet_db_t` set `user_id`=?, `nick_name`=?, `use_ga`=?, " +
		"`public_addr`=?, `secret_key`=?, `reg_time`=? " + "where id=?"
	_, err := this.engine.Exec(sql, wdbt.UserId, wdbt.NickName, wdbt.UseGa, wdbt.PublicAddr, wdbt.SecretKey,
		wdbt.RegTime.Format("2006-01-02 15:04:05"), wdbt.Id)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : UpdateWalletInfo", "%v", err))
	}
	return err
}

func (this *WalletTableOperation) QuaryAll(userID uint64) ([]*WalletDbT, error) {
	tmp := make([]*WalletDbT, 0)
	err := this.engine.Where("user_id =?", userID).Asc("id").Find(&tmp)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : QuaryAll", "%v", err))
	}
	return tmp, err
}

func (this *WalletTableOperation) DeleteWalletInfo(wdbt *WalletDbT) error {
	_, err := this.engine.Delete(wdbt)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : DeleteWalletInfo", "%v", err))
	}
	return err
}

func (this *WalletTableOperation) SearchWalletInfo(addr, nickname string) ([]*WalletDbT, error) {
	tmp := make([]*WalletDbT, 0)
	sqlcmd := ""
	if len(addr) > 0 {
		sqlcmd = `public_addr LIKE '%` + addr + `%'`
	}

	if len(nickname) > 0 {
		if len(sqlcmd) > 0 {
			sqlcmd += " and "
		}
		sqlcmd += `nick_name LIKE '%` + nickname + `%'`
	}

	err := this.engine.Select("*").Where(sqlcmd).Find(&tmp)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : SearchWalletInfo", "%v", err))
	}
	return tmp, err
}

func (this *WalletTableOperation) createWalletID(src int64) (id uint64) {
	// 10000以下ID值保留，普通用户ID最小10001
	if src == 0 {
		id = uint64(10001)
	} else {
		tmp := &WalletDbT{}
		b, e := this.engine.Limit(1).Desc("id").Get(tmp)
		if e != nil {
			logmodels.Logger.Error(logmodels.SPrintError("WalletTableOperation : createWalletID", "%v", e))
			id = uint64(src + 10001)
			return
		}
		if b {
			id = tmp.Id + 1
		}
	}
	return
}
