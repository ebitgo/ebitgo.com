package dbmodels

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ebitgo/ebitgo.com/models/logger"
	"github.com/go-xorm/xorm"
)

type UserIdTableOperation struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

func (this *UserIdTableOperation) Init(e *xorm.Engine) error {
	if e == nil {
		err := errors.New("Can not get database engine!")
		logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : Init", "%v", err))
		return err
	}
	this.locker = &sync.Mutex{}
	this.engine = e
	return nil
}

func (this *UserIdTableOperation) GetKey() string {
	return DB_USER_ID_OPERATION
}

func (this *UserIdTableOperation) Quary(qtype int, v ...interface{}) (err error) {
	if v == nil || len(v) < 1 {
		err = errors.New("Quary get use id , must input UserIdT struct pointer as parameter")
		logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : Quary", "%v", err))
		return
	}
	uidTab := v[0].(*UserIdT)

	switch qtype {
	case QT_CHECK_RECORD:
		this.locker.Lock()
		uid := this.GetUserId(uidTab.Email, uidTab)
		logmodels.Logger.Trace(logmodels.SPrintTrace("UserIdTableOperation : Quary", "Check user id = %d", uid))
		this.locker.Unlock()
		return
	case QT_ADD_RECORD:
		this.locker.Lock()
		err = this.AddUser(uidTab)
		this.locker.Unlock()
		return
	}

	err = errors.New(fmt.Sprintf("Quary type is not defined(%d)", qtype))
	logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : Quary", "%v", err))
	return
}

func (this *UserIdTableOperation) GetUserId(email string, uidTab *UserIdT) uint64 {
	if len(uidTab.Email) == 0 {
		uidTab.Email = email
	}
	b, e := this.engine.Get(uidTab)
	if e != nil {
		logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : GetUserId", "%v", e))
		return 0
	}
	if b {
		return uidTab.Id
	}
	return 0
}

func (this *UserIdTableOperation) AddUser(uidTab *UserIdT) error {
	counts, err := this.engine.Count(&UserIdT{})
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : AddUser", "Quary UserIdT count has error \r\n\t%v", err))
		return err
	}

	uidTab.Id = this.createUserID(counts)
	_, err = this.engine.Insert(uidTab)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : AddUser", "Insert userinfo to UserIdT has error \r\n\t%v", err))
	}
	return err
}

func (this *UserIdTableOperation) createUserID(src int64) (id uint64) {
	// 20000以下ID值保留，普通用户ID最小20001
	if src == 0 {
		id = uint64(20001)
	} else {
		tmp := &UserIdT{}
		b, e := this.engine.Limit(1).Desc("id").Get(tmp)
		if e != nil {
			logmodels.Logger.Error(logmodels.SPrintError("UserIdTableOperation : createUserID", "%v", e))
			id = uint64(src + 20001)
			return
		}
		if b {
			id = tmp.Id + 1
		}
	}
	return
}
