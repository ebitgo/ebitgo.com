package models

import (
	"errors"
	"time"
	// "github.com/ebitgo/ebitgo.com/models/databasemodels"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/logger"
)

type WalletAddOper struct {
	WalletBaseInfo
}

func (this *WalletAddOper) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:         true,
			RT_M_USEREMAIL:       this.LoginName,
			RT_M_WALLET_NICKNAME: this.WalletInfo.NickName,
			RT_M_AUTH_UPDATE:     true,
			RT_M_USER_AUTH:       this.UserInfo.AuthStr,
		}
		if this.IsExist {
			ret.Data.(map[string]interface{})[RT_M_GA_STATE] = this.UserInfo.UseGa
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *WalletAddOper) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	err := this.WalletBaseInfo.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	if this.IsExistWallet {
		err = errors.New("Wallet nick name already exist!")
		logmodels.Logger.Error(logmodels.SPrintError("WalletAddOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	timenow := time.Now()

	this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
		this.LoginName, this.UserInfo.Password, "")

	this.UserInfo.LastLoginTime = timenow

	this.ErrorMsg = DBManagerInst.Update_UserInfo(this.UserInfo)
	if this.ErrorMsg == nil {
		this.ErrorMsg = DBManagerInst.Add_Wallet(this.WalletInfo)
	}

	return this.ErrorMsg
}

func (this *WalletAddOper) DecodeContext(vals url.Values) {
	this.WalletBaseInfo.DecodeContext(vals)
	if this.ErrorMsg != nil {
		return
	}

	this.WalletInfo.SecretKey = vals.Get(POST_MARK_WALLET_SKEY)
	if this.checkValue(this.WalletInfo.SecretKey, POST_MARK_WALLET_SKEY) != nil {
		return
	}

	this.WalletInfo.UseGa = vals.Get(POST_MARK_MODIFY_GA) != "0"
	this.checkValue(vals.Get(POST_MARK_MODIFY_GA), POST_MARK_MODIFY_GA)
}
