package models

import (
	// "errors"
	"errors"
	// "fmt"
	"net/url"
	"time"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type AccountModifyPassword struct {
	AccountBaseInfo
	UserAuths    string
	UserFid      string
	OrigPassword string
	NewPassword  string
}

func (this *AccountModifyPassword) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:     true,
			RT_M_USEREMAIL:   this.LoginName,
			RT_M_USER_AUTH:   this.UserInfo.AuthStr,
			RT_M_AUTH_UPDATE: true,
			RT_M_USER_LVL:    this.UserInfo.UserLevel,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *AccountModifyPassword) QuaryExecute() error {
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
		logmodels.Logger.Error(logmodels.SPrintError("AccountModifyPassword : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 检查密码是否正确
	if this.OrigPassword != this.UserInfo.Password {
		err = errors.New("User orig password is incorrect!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountModifyPassword : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 检查FID是否正确
	if this.UserFid != this.UserInfo.ActiveFlag {
		err = errors.New("User fid is incorrect!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountModifyPassword : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 检查Auth是否正确
	if this.UserAuths != this.UserInfo.AuthStr {
		err = errors.New("User not logged in, please login at first!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountModifyPassword : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 未激活账户，不支持修改密码操作
	if this.UserInfo.UserLevel == 0 {
		err = errors.New("User email not activation, please active your email at first!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountModifyPassword : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	timenow := time.Now()

	this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
		this.LoginName, this.NewPassword, "")

	this.UserInfo.LastLoginTime = timenow
	this.UserInfo.Password = this.NewPassword

	err = DBManagerInst.Update_UserInfo(this.UserInfo)

	return err
}

func (this *AccountModifyPassword) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.UserInfo = new(dbmodels.UserDbT)

	this.UserFid = vals.Get(POST_MARK_FID)
	if this.match_checkValue(this.UserFid, POST_MARK_FID) != nil {
		return
	}

	this.OrigPassword = vals.Get(POST_MARK_PASS_WORD)
	if this.checkValue(this.OrigPassword, POST_MARK_PASS_WORD) != nil {
		return
	}

	this.UserAuths = vals.Get(POST_MARK_AUTHCODE)
	if this.checkValue(this.UserAuths, POST_MARK_AUTHCODE) != nil {
		return
	}

	this.NewPassword = vals.Get(POST_MARK_NEW_PASS_WORD)
	this.checkValue(this.NewPassword, POST_MARK_NEW_PASS_WORD)
}
