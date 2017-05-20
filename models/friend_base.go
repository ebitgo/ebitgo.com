package models

import (
	// "errors"
	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	// "github.com/ebitgo/ebitgo.com/models/logger"
	"net/url"
	"strconv"
	// "time"
)

type FriendBaseOperation struct {
	WalletBaseInfo
	IsFriendExist bool
	FriendInfo    *dbmodels.FriendlyShipDbT
}

func (this *FriendBaseOperation) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS: true,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *FriendBaseOperation) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	err := this.WalletBaseInfo.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	this.FriendInfo.UserId = this.UserInfo.Id
	this.FriendInfo.NickName = this.WalletInfo.NickName
	this.FriendInfo.PublicAddr = this.WalletInfo.PublicAddr
	this.IsFriendExist = false

	ret, err := DBManagerInst.Quary_Get_Friend(this.FriendInfo.UserId, this.FriendInfo.WalletId)
	if err == nil {
		this.IsFriendExist = (ret.UserId != 0 && this.FriendInfo.WalletId != 0)
	}

	return this.ErrorMsg
}

func (this *FriendBaseOperation) DecodeContext(vals url.Values) {
	this.WalletBaseInfo.DecodeContext(vals)
	if this.ErrorMsg != nil {
		return
	}

	wid := vals.Get(POST_MARK_WALLET_ID)
	if this.checkValue(wid, POST_MARK_WALLET_ID) != nil {
		return
	}

	this.FriendInfo = &dbmodels.FriendlyShipDbT{}

	this.FriendInfo.WalletId, this.ErrorMsg = strconv.ParseUint(wid, 10, 64)
}
