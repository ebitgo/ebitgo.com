package models

import (
	"errors"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
	// "time"
)

type WalletUpdateSKEYOper struct {
	WalletBaseInfo
	USSKEY string
}

func (this *WalletUpdateSKEYOper) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:         true,
			RT_M_USEREMAIL:       this.LoginName,
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

func (this *WalletUpdateSKEYOper) QuaryExecute() error {
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
		logmodels.Logger.Error(logmodels.SPrintError("WalletUpdateSKEYOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}
	this.WalletInfo.SecretKey = this.USSKEY
	this.ErrorMsg = DBManagerInst.Update_WalletInfo(this.WalletInfo)

	return this.ErrorMsg
}

func (this *WalletUpdateSKEYOper) DecodeContext(vals url.Values) {
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

	this.USSKEY = vals.Get(POST_MARK_WALLET_SKEY)
	this.checkValue(this.USSKEY, POST_MARK_WALLET_SKEY)
}
