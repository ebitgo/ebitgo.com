package models

import (
	// "errors"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type AccountActivation struct {
	AccountBaseInfo
	UserPassword string
	UserFid      string
	UserTid      string
}

func (this *AccountActivation) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:   true,
			RT_M_USEREMAIL: this.LoginName,
			RT_M_USER_LVL:  this.UserInfo.UserLevel,
			RT_M_USER_AUTH: this.UserInfo.AuthStr,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *AccountActivation) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 先检查有没有当前登录用户
	err := this.AccountBaseInfo.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}
	// 用户不存在
	if !this.IsExist {
		err = errors.New("User not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountActivation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 检查密码是否正确
	if this.UserPassword != this.UserInfo.Password {
		err = errors.New("User password is incorrect!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountActivation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 检查FID是否正确
	if this.UserFid != this.UserInfo.ActiveFlag {
		err = errors.New("User fid is incorrect!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountActivation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	uid := fmt.Sprintf("%d", this.UserInfo.Id)
	// 检查TID是否正确
	if this.UserTid != uid {
		err = errors.New("User tid is incorrect!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountActivation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	if this.UserInfo.UserLevel == 0 {
		this.UserInfo.UserLevel += 1
	}

	timenow := time.Now()

	hashstr := getSHA256Key(timenow.Format("2006-01-02 15:04:05"), this.LoginName, this.UserPassword, "")
	this.UserInfo.AuthStr = hashstr

	this.UserInfo.LastLoginTime = timenow

	err = DBManagerInst.Update_UserInfo(this.UserInfo)

	return err
}

func (this *AccountActivation) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.UserFid = vals.Get(POST_MARK_FID)
	if this.match_checkValue(this.UserFid, POST_MARK_FID) != nil {
		return
	}

	this.UserPassword = vals.Get(POST_MARK_PASS_WORD)
	if this.checkValue(this.UserPassword, POST_MARK_PASS_WORD) != nil {
		return
	}

	this.UserInfo = new(dbmodels.UserDbT)

	this.UserTid = vals.Get(POST_MARK_TID)
	this.match_checkValue(this.UserTid, POST_MARK_TID)
}
