package models

import (
	"errors"
	"net/url"
	"time"

	"github.com/dgryski/dgoogauth"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type AccountUserGAConfig struct {
	AccountBaseInfo
	Type         string
	UserPassword string
	GAKey        string
	GACode       string
	GAURI        string
	Auth         string
}

func (this *AccountUserGAConfig) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_USEREMAIL: this.LoginName,
			RT_M_GA_STATE:  this.UserInfo.UseGa,
			RT_M_USER_LVL:  this.UserInfo.UserLevel,
		}

		if this.Type == GA_PT_MODIFY_GET {
			ret.Data.(map[string]interface{})[RT_M_GA_KEY] = this.UserInfo.GaKey
			ret.Data.(map[string]interface{})[RT_M_GA_URI] = this.GAURI
		} else if this.Type == GA_PT_MODIFY_NEW {
			ret.Data.(map[string]interface{})[RT_M_SUCCESS] = true
			ret.Data.(map[string]interface{})[RT_M_AUTH_UPDATE] = true
			ret.Data.(map[string]interface{})[RT_M_USER_AUTH] = this.UserInfo.AuthStr
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *AccountUserGAConfig) QuaryExecute() error {
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
		logmodels.Logger.Error(logmodels.SPrintError("AccountUserGAConfig : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	if this.Type == GA_PT_MODIFY_DELETE {
		// 检查密码是否正确
		if this.UserPassword != this.UserInfo.Password {
			err = errors.New("User password is incorrect!")
			logmodels.Logger.Error(logmodels.SPrintError("AccountUserGAConfig : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}
		// 检查GAcode
		if b := verifyGoogleAuth(this.UserInfo.GaKey, this.GACode); !b {
			err = errors.New("Google auth code is incorrect!")
			logmodels.Logger.Error(logmodels.SPrintError("AccountUserGAConfig : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}
		this.UserInfo.UseGa = false
		this.UserInfo.UserLevel -= 1
		this.UserInfo.GaKey = ""

		err = DBManagerInst.Update_UserInfo(this.UserInfo)

	} else if this.Type == GA_PT_MODIFY_GET {
		this.UserInfo.GaKey = randomBase32String()

		var totp dgoogauth.OTPConfig
		totp.Secret = this.UserInfo.GaKey
		this.GAURI = totp.ProvisionURIWithIssuer(this.LoginName, "eBitGo.com")
	} else {
		if this.UserInfo.AuthStr != this.Auth {
			err = errors.New("Login user timeout or have already login in other places!")
			logmodels.Logger.Error(logmodels.SPrintError("AccountUserGAConfig : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}

		if !(verifyGoogleAuth(this.GAKey, this.GACode)) {
			err = errors.New("Google auth code is incorrect!")
			logmodels.Logger.Error(logmodels.SPrintError("AccountUserGAConfig : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}

		timenow := time.Now()

		this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
			this.LoginName, this.UserInfo.Password, "")

		// fmt.Println("LastLoginTime = ", uInfo.LastLoginTime.Format("2006-01-02 15:04:05"))
		this.UserInfo.LastLoginTime = timenow

		this.UserInfo.GaKey = this.GAKey

		this.UserInfo.UseGa = true
		this.UserInfo.UserLevel += 1

		err = DBManagerInst.Update_UserInfo(this.UserInfo)
	}

	return err
}

func (this *AccountUserGAConfig) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.Type = vals.Get(POST_MARK_MODIFY_GA)
	if this.checkValue(this.Type, POST_MARK_MODIFY_GA) != nil {
		return
	}

	if this.Type == GA_PT_MODIFY_DELETE {

		this.GACode = vals.Get(POST_MARK_GACODE)
		if this.match_checkValue(this.GACode, POST_MARK_GACODE) != nil {
			return
		}

		this.UserPassword = vals.Get(POST_MARK_PASS_WORD)
		this.checkValue(this.UserPassword, POST_MARK_PASS_WORD)

	} else if this.Type == GA_PT_MODIFY_NEW {

		this.GAKey = vals.Get(POST_MARK_GAKEY)
		if this.match_checkValue(this.GAKey, POST_MARK_GAKEY) != nil {
			return
		}

		this.Auth = vals.Get(POST_MARK_AUTHCODE)
		if this.checkValue(this.Auth, POST_MARK_AUTHCODE) != nil {
			return
		}

		this.GACode = vals.Get(POST_MARK_GACODE)
		this.match_checkValue(this.GACode, POST_MARK_GACODE)
	}
}
