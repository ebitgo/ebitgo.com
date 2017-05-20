package dbmodels

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ebitgo/ebitgo.com/models/logger"
	"github.com/go-xorm/xorm"
)

type UserTableOperation struct {
	engine *xorm.Engine
	locker *sync.Mutex
}

func (this *UserTableOperation) Init(e *xorm.Engine) error {
	if e == nil {
		err := errors.New("Can not get database engine!")
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : Init", "%v", err))
		return err
	}
	this.locker = &sync.Mutex{}
	this.engine = e
	return nil
}

func (this *UserTableOperation) GetKey() string {
	return DB_USER_INFO_OPERATION
}

func (this *UserTableOperation) Quary(qtype int, v ...interface{}) (err error) {
	if v == nil || len(v) < 1 {
		err = errors.New("Quary get use id , must input RegUserDbT struct pointer as parameter")
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : Quary", "%v", err))
		return
	}
	rut := v[0].(*UserDbT)

	switch qtype {
	case QT_ADD_RECORD:
		this.locker.Lock()
		err = this.AddUserInfo(rut)
		this.locker.Unlock()
		return
	case QT_GET_RECORD:
		this.locker.Lock()
		err = this.GetUserInfo(rut)
		this.locker.Unlock()
		return
	case QT_UPDATE_RECORD:
		this.locker.Lock()
		err = this.UpdateUserInfo(rut)
		this.locker.Unlock()
		return
	}

	err = errors.New(fmt.Sprintf("Quary type is not defined(%d)", qtype))
	logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : Quary", "%v", err))
	return
}

func (this *UserTableOperation) AddUserInfo(rut *UserDbT) error {
	if rut.Id == 0 {
		err := errors.New("Add user info Fail! (UserDbT.Id = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : AddUserInfo", "%v", err))
		return err
	}

	rut.RegTime = time.Now()
	rut.LastLoginTime = time.Now()
	_, err := this.engine.Insert(rut)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : AddUserInfo", "Insert user to UserDbT has error \r\n\t%v", err))
	}
	return err
}

func (this *UserTableOperation) GetUserInfo(rut *UserDbT) error {
	if rut.Id == 0 {
		err := errors.New("Get user info Fail! (UserDbT.Id = 0)")
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : GetUserInfo", "%v", err))
		return err
	}
	b, err := this.engine.Get(rut)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : GetUserInfo", "%v", err))
		return err
	}
	if !b {
		err := errors.New(fmt.Sprintf("Get user info Fail! (UserDbT.Id = %d is not exist)", rut.Id))
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : GetUserInfo", "%v", err))
		return err
	}
	return nil
}

func (this *UserTableOperation) UpdateUserInfo(rut *UserDbT) error {
	sql := "update `user_db_t` set `id`=?, `password`=?, `mobile_phone`=?, " +
		"`use_ga`=?, `ga_key`=?, `reg_time`=?, `last_login_time`=?, " +
		"`auth_str`=?, `user_level`=?, `active_flag`=? " +
		"where id=?"
	_, err := this.engine.Exec(sql, rut.Id, rut.Password, rut.MobilePhone, rut.UseGa, rut.GaKey,
		rut.RegTime.Format("2006-01-02 15:04:05"), rut.LastLoginTime.Format("2006-01-02 15:04:05"),
		rut.AuthStr, rut.UserLevel, rut.ActiveFlag, rut.Id)
	// _, err := this.engine.Update(rut)
	if err != nil {
		logmodels.Logger.Error(logmodels.SPrintError("UserTableOperation : UpdateUserInfo", "%v", err))
	}
	return err
}
