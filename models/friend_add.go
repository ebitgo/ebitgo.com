package models

import (
	"errors"
	"net/url"

	// "github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
	// "strconv"
)

type FriendAddOperation struct {
	FriendBaseOperation
}

func (this *FriendAddOperation) GetResultData() *ReslutOperation {
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

func (this *FriendAddOperation) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	err := this.FriendBaseOperation.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	if !this.IsExistWallet {
		err = errors.New("Wallet " + this.WalletInfo.NickName + " not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("FriendAddOperation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	if this.IsFriendExist {
		err = errors.New("You have added " + this.WalletInfo.NickName + " as a friend!")
		logmodels.Logger.Error(logmodels.SPrintError("FriendAddOperation : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	this.ErrorMsg = DBManagerInst.Add_Friend(this.FriendInfo)

	return this.ErrorMsg
}

func (this *FriendAddOperation) DecodeContext(vals url.Values) {
	this.FriendBaseOperation.DecodeContext(vals)
}
