package models

import (
	// "errors"
	"errors"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type WalletBaseInfo struct {
	AccountBaseInfo
	IsExistWallet bool
	LogAuth       string
	WalletInfo    *dbmodels.WalletDbT
}

func (this *WalletBaseInfo) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_WALLET_NICKNAME: this.WalletInfo.NickName,
			RT_M_EXIST:           this.IsExistWallet,
		}
		if this.IsExist {
			ret.Data.(map[string]interface{})[RT_M_GA_STATE] = this.WalletInfo.UseGa
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *WalletBaseInfo) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 检查登录用户是否有效
	err := this.AccountBaseInfo.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	if !this.IsExist { // 用户不存在
		err = errors.New("User not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("WalletBaseInfo : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	} else if this.LogAuth != this.UserInfo.AuthStr {
		err = errors.New("Login user timeout or have already login in other places!")
		logmodels.Logger.Error(logmodels.SPrintError("WalletBaseInfo : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	this.WalletInfo.UserId = this.UserInfo.Id

	// 检查数据库
	this.IsExistWallet = false
	uid, e := DBManagerInst.Quary_Check_Wallet_Exist(this.WalletInfo.NickName)
	if e == nil {
		this.IsExistWallet = uid > 0

		if this.IsExistWallet {
			this.WalletInfo, e = DBManagerInst.Quary_WalletInfo(uid)
		}
	}
	this.ErrorMsg = e
	return e
}

func (this *WalletBaseInfo) DecodeContext(vals url.Values) {
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

	this.WalletInfo.PublicAddr = vals.Get(POST_MARK_WALLET_PUBADDR)
	this.checkValue(this.WalletInfo.PublicAddr, POST_MARK_WALLET_PUBADDR)
}
