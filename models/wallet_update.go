package models

import (
	"errors"
	"net/url"
	"time"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type WalletUpdateOper struct {
	WalletBaseInfo
	MNickName string
	MGAUsed   bool
}

func (this *WalletUpdateOper) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:         true,
			RT_M_USEREMAIL:       this.LoginName,
			RT_M_AUTH_UPDATE:     true,
			RT_M_USER_AUTH:       this.UserInfo.AuthStr,
			RT_M_WALLET_NICKNAME: this.WalletInfo.NickName,
			RT_M_GA_STATE:        this.WalletInfo.UseGa,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *WalletUpdateOper) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 检查登录用户是否有效
	this.ErrorMsg = this.WalletBaseInfo.QuaryExecute()
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	if !this.IsExistWallet {
		err := errors.New("Wallet nick name not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("WalletUpdateOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	if this.MNickName != this.WalletInfo.NickName {

		// 检查新的名字是否被占用
		uid, e := DBManagerInst.Quary_Check_Wallet_Exist(this.MNickName)
		if e == nil {
			if uid > 0 {
				err := errors.New("Wallet nick name [ " + this.MNickName + " ] already exist!")
				logmodels.Logger.Error(logmodels.SPrintError("WalletUpdateOper : QuaryExecute", "%v", err))
				this.ErrorMsg = err
				return err
			}
		} else {
			logmodels.Logger.Error(logmodels.SPrintError("WalletUpdateOper : QuaryExecute", "%v", e))
			this.ErrorMsg = e
			return e
		}

	}

	timenow := time.Now()

	this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
		this.LoginName, this.UserInfo.Password, "")

	this.UserInfo.LastLoginTime = timenow

	this.ErrorMsg = DBManagerInst.Update_UserInfo(this.UserInfo)
	if this.ErrorMsg == nil {
		this.WalletInfo.NickName = this.MNickName
		this.WalletInfo.UseGa = this.MGAUsed
		this.ErrorMsg = DBManagerInst.Update_WalletInfo(this.WalletInfo)
	}

	return this.ErrorMsg
}

func (this *WalletUpdateOper) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.LogAuth = vals.Get(POST_MARK_AUTHCODE)
	if this.checkValue(this.LogAuth, POST_MARK_AUTHCODE) != nil {
		return
	}

	this.WalletInfo = &dbmodels.WalletDbT{}

	this.WalletInfo.NickName = vals.Get(POST_MARK_WALLET_NICKNAME)
	if this.match_checkValue(this.WalletInfo.NickName, POST_MARK_WALLET_NICKNAME) != nil {
		return
	}

	this.MNickName = vals.Get(POST_MARK_MODIFY_NICKNAME)
	if this.match_checkValue(this.MNickName, POST_MARK_MODIFY_NICKNAME) != nil {
		return
	}

	this.MGAUsed = vals.Get(POST_MARK_MODIFY_GA) != "0"
	this.checkValue(vals.Get(POST_MARK_MODIFY_GA), POST_MARK_MODIFY_GA)
}
