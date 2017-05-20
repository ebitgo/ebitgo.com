package dbmodels

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ebitgo/ebitgo.com/models/logger"
	"github.com/go-xorm/xorm"
)

type FriendlyShipTOperation struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

func (this *FriendlyShipTOperation) Init(e *xorm.Engine) error {
	if e == nil {
		err := errors.New("Can not get database engine!")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : Init", "%v", err))
		return err
	}
	this.locker = &sync.Mutex{}
	this.engine = e
	return nil
}

func (this *FriendlyShipTOperation) GetKey() string {
	return DB_FRIENDLY_SHIP_OPERATION
}

func (this *FriendlyShipTOperation) Quary(qtype int, v ...interface{}) (err error) {
	if v == nil || len(v) < 1 {
		err = errors.New("Quary get use id , must input WalletDbT struct pointer as parameter")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : Quary", "%v", err))
		return
	}

	switch qtype {
	case QT_ADD_RECORD:
		rut := v[0].(*FriendlyShipDbT)
		this.locker.Lock()
		err = this.AddFriendShip(rut)
		this.locker.Unlock()
		return
	case QT_CHECK_RECORD:
		rut := v[0].(*FriendlyShipDbT)
		this.locker.Lock()
		err = this.GetFriend(rut)
		this.locker.Unlock()
		return
	// case QT_GET_RECORD:
	// 	rut := v[0].(*FriendlyShipDbT)
	// 	this.locker.Lock()
	// 	err = this.GetWalletInfo(rut)
	// 	this.locker.Unlock()
	// 	return
	// case QT_UPDATE_RECORD:
	// 	rut := v[0].(*FriendlyShipDbT)
	// 	this.locker.Lock()
	// 	err = this.UpdateWalletInfo(rut)
	// 	this.locker.Unlock()
	// 	return
	case QT_QUARY_ALL_RECORD:
		ruts := v[0].(map[string][]*FriendlyShipDbT)
		uid := v[1].(uint64)
		this.locker.Lock()
		ruts["RESULT"], err = this.QuaryAll(uid)
		this.locker.Unlock()
		return
	case QT_DELETE_RECORD:
		rut := v[0].(*FriendlyShipDbT)
		this.locker.Lock()
		err = this.DeleteFriendInfo(rut)
		this.locker.Unlock()
		return
		// case QT_SEARCH_RECORD:
		// 	ruts := v[0].(map[string][]*FriendlyShipDbT)
		// 	addr := v[1].(string)
		// 	nickname := v[2].(string)
		// 	this.locker.Lock()
		// 	ruts["RESULT"], err = this.SearchWalletInfo(addr, nickname)
		// 	this.locker.Unlock()
		// 	return
	}

	err = errors.New(fmt.Sprintf("Quary type is not defined(%d)", qtype))
	logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : Quary", "%v", err))
	return
}

func (this *FriendlyShipTOperation) AddFriendShip(info *FriendlyShipDbT) error {
	if info.UserId == 0 {
		err := errors.New("Add friend info Fail! (FriendlyShipDbT.UserId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : AddFriendShip", "%v", err))
		return err
	}
	if info.WalletId == 0 {
		err := errors.New("Add friend info Fail! (FriendlyShipDbT.WalletId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : AddFriendShip", "%v", err))
		return err
	}

	_, err := this.engine.Insert(info)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : AddFriendShip", "Insert friend to WalletDbT has error \r\n\t%v", err))
	}
	return err
}

func (this *FriendlyShipTOperation) GetFriend(info *FriendlyShipDbT) error {
	if info.UserId == 0 {
		err := errors.New("Add friend info Fail! (FriendlyShipDbT.UserId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : GetFriend", "%v", err))
		return err
	}
	if info.WalletId == 0 {
		err := errors.New("Add friend info Fail! (FriendlyShipDbT.WalletId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : GetFriend", "%v", err))
		return err
	}

	b, err := this.engine.Get(info)

	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : GetFriend", "%v", err))
	}
	if !b {
		info.UserId = 0
		info.WalletId = 0
	}
	return err
}

func (this *FriendlyShipTOperation) QuaryAll(uid uint64) ([]*FriendlyShipDbT, error) {
	tmp := make([]*FriendlyShipDbT, 0)
	err := this.engine.Where("user_id =?", uid).Asc("nick_name").Find(&tmp)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : QuaryAll", "%v", err))
	}
	return tmp, err
}

func (this *FriendlyShipTOperation) DeleteFriendInfo(info *FriendlyShipDbT) error {
	if info.UserId == 0 {
		err := errors.New("Delete friend info Fail! (FriendlyShipDbT.UserId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : DeleteFriendInfo", "%v", err))
		return err
	}
	if info.WalletId == 0 {
		err := errors.New("Delete friend info Fail! (FriendlyShipDbT.WalletId = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : DeleteFriendInfo", "%v", err))
		return err
	}

	_, err := this.engine.Delete(info)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("FriendlyShipTOperation : DeleteFriendInfo", "%v", err))
	}
	return err
}
