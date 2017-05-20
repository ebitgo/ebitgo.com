package models

import (
	"errors"
	// "fmt"
	"net/url"
	"time"

	"github.com/ebitgo/ebitgo.com/models/logger"
)

type AccountAuthVerify struct {
	AccountBaseInfo
	IsHeartBit bool
	AuthKey    string
	Success    bool
	UpdateAuth bool
}

func (this *AccountAuthVerify) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:     this.Success,
			RT_M_USEREMAIL:   this.LoginName,
			RT_M_AUTH_UPDATE: this.UpdateAuth,
			RT_M_USER_AUTH:   this.AuthKey,
			RT_M_USER_LVL:    this.UserInfo.UserLevel,
			RT_M_GA_STATE:    this.UserInfo.UseGa,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *AccountAuthVerify) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}
	err := this.AccountBaseInfo.QuaryExecute()
	if err != nil {
		return err
	}

	if !this.IsExist {
		err = errors.New("Account is not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountAuthVerify : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	authStr := getSHA256Key(this.UserInfo.LastLoginTime.Format("2006-01-02 15:04:05"), this.LoginName,
		this.UserInfo.Password, "")
	this.Success = this.AuthKey == authStr

	if this.IsHeartBit {
		this.UpdateAuth = this.Success
		if this.UpdateAuth {
			this.UserInfo.LastLoginTime = time.Now()
			this.UserInfo.AuthStr = getSHA256Key(this.UserInfo.LastLoginTime.Format("2006-01-02 15:04:05"), this.LoginName,
				this.UserInfo.Password, "")
			this.AuthKey = this.UserInfo.AuthStr
			err = DBManagerInst.Update_UserInfo(this.UserInfo)
		}
	} else {
		this.UpdateAuth = false
	}
	return err
}

func (this *AccountAuthVerify) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.IsHeartBit = len(vals.Get(POST_MARK_HEART)) > 0

	this.AuthKey = vals.Get(POST_MARK_AUTHCODE)
	this.match_checkValue(this.AuthKey, POST_MARK_AUTHCODE)
}
