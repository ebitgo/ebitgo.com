package models

import (
	"errors"
	"time"
	// "github.com/ebitgo/ebitgo.com/models/databasemodels"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/logger"
)

type WalletDeleteOper struct {
	WalletBaseInfo
	LoginPasswordorGAcode string
}

func (this *WalletDeleteOper) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:         true,
			RT_M_USEREMAIL:       this.LoginName,
			RT_M_WALLET_NICKNAME: this.WalletInfo.NickName,
			RT_M_AUTH_UPDATE:     true,
			RT_M_USER_AUTH:       this.UserInfo.AuthStr,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *WalletDeleteOper) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	err := this.WalletBaseInfo.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	if this.WalletInfo.UseGa {
		// 如果使用GA
		if !(verifyGoogleAuth(this.UserInfo.GaKey, this.LoginPasswordorGAcode)) {
			err = errors.New("Google auth code is incorrect!")
			logmodels.Logger.Error(logmodels.SPrintError("UserLoginOperation : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}
	} else {
		if this.LoginPasswordorGAcode != this.UserInfo.Password {
			err = errors.New("User login password is incorrect!")
			logmodels.Logger.Error(logmodels.SPrintError("WalletDeleteOper : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}
	}

	if !this.IsExistWallet {
		err = errors.New("Wallet nick name [" + this.WalletInfo.NickName + "] not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("WalletDeleteOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	timenow := time.Now()

	this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
		this.LoginName, this.UserInfo.Password, "")

	this.UserInfo.LastLoginTime = timenow

	this.ErrorMsg = DBManagerInst.Update_UserInfo(this.UserInfo)
	if this.ErrorMsg == nil {
		this.ErrorMsg = DBManagerInst.Delete_Wallet(this.WalletInfo)
	}

	return this.ErrorMsg
}

func (this *WalletDeleteOper) DecodeContext(vals url.Values) {
	this.WalletBaseInfo.DecodeContext(vals)
	if this.ErrorMsg != nil {
		return
	}

	this.WalletInfo.SecretKey = vals.Get(POST_MARK_WALLET_SKEY)
	if this.checkValue(this.WalletInfo.SecretKey, POST_MARK_WALLET_SKEY) != nil {
		return
	}

	this.LoginPasswordorGAcode = vals.Get(POST_MARK_PASS_WORD)
	if this.checkValue(this.LoginPasswordorGAcode, POST_MARK_PASS_WORD) != nil {
		return
	}

	this.WalletInfo.UseGa = vals.Get(POST_MARK_MODIFY_GA) != "0"
	this.checkValue(vals.Get(POST_MARK_MODIFY_GA), POST_MARK_MODIFY_GA)
}
