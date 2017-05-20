package models

import (
	"errors"
	// "fmt"
	"net/url"
	"time"

	"github.com/ebitgo/ebitgo.com/models/logger"
)

type UserLoginOperation struct {
	AccountBaseInfo
	UserPassword string
	GaCode       string
	VerifInfo    *VerifiFormular
}

func (this *UserLoginOperation) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:   true,
			RT_M_USEREMAIL: this.LoginName,
			RT_M_USER_LVL:  this.UserInfo.UserLevel,
			RT_M_USER_AUTH: this.UserInfo.AuthStr,
			RT_M_GA_STATE:  this.UserInfo.UseGa,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *UserLoginOperation) QuaryExecute() error {
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
		logmodels.Logger.Error(logmodels.SPrintError("UserLoginOperation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 检查密码是否正确
	if this.UserPassword != this.UserInfo.Password {
		err = errors.New("User password is incorrect!")
		logmodels.Logger.Error(logmodels.SPrintError("UserLoginOperation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	if this.UserInfo.UseGa {
		// 如果使用GA
		if !(verifyGoogleAuth(this.UserInfo.GaKey, this.GaCode)) {
			err = errors.New("Google auth code is incorrect!")
			logmodels.Logger.Error(logmodels.SPrintError("UserLoginOperation : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}
	} else if this.VerifInfo.Verify() == false {
		err = errors.New("ViCode is not valid!")
		logmodels.Logger.Info(logmodels.SPrintInfo("UserLoginOperation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	timenow := time.Now()

	this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
		this.LoginName, this.UserPassword, "")

	// fmt.Println("LastLoginTime = ", uInfo.LastLoginTime.Format("2006-01-02 15:04:05"))
	this.UserInfo.LastLoginTime = timenow
	err = DBManagerInst.Update_UserInfo(this.UserInfo)
	this.ErrorMsg = err
	return err
}

func (this *UserLoginOperation) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.UserPassword = vals.Get(POST_MARK_PASS_WORD)
	if this.checkValue(this.UserPassword, POST_MARK_PASS_WORD) != nil {
		return
	}

	this.GaCode = vals.Get(POST_MARK_GACODE)
	if len(this.GaCode) != 0 {
		if this.match_checkValue(this.UserPassword, POST_MARK_PASS_WORD) != nil {
			return
		}
	}

	if this.VerifInfo == nil {
		this.VerifInfo = &VerifiFormular{}
	}

	this.VerifInfo.ViCode = vals.Get(POST_MARK_REG_VICODE)

	this.VerifInfo.Value1 = vals.Get(POST_MARK_REG_VALUE1)

	this.VerifInfo.Value2 = vals.Get(POST_MARK_REG_VALUE2)

	this.VerifInfo.Value3 = vals.Get(POST_MARK_REG_VALUE3)

	this.VerifInfo.Operator = vals.Get(POST_MARK_REG_OPERA)
}
