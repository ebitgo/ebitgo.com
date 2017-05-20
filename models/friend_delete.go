package models

import (
	"errors"
	"time"
	// "github.com/ebitgo/ebitgo.com/models/databasemodels"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/logger"
)

type FriendDeleteOper struct {
	FriendBaseOperation
}

func (this *FriendDeleteOper) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:     true,
			RT_M_USEREMAIL:   this.LoginName,
			RT_M_AUTH_UPDATE: true,
			RT_M_USER_AUTH:   this.UserInfo.AuthStr,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *FriendDeleteOper) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	err := this.FriendBaseOperation.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	if !this.IsFriendExist {
		err = errors.New("Friend nick name [" + this.WalletInfo.NickName + "] not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("FriendDeleteOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	timenow := time.Now()

	this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
		this.LoginName, this.UserInfo.Password, "")

	this.UserInfo.LastLoginTime = timenow

	this.ErrorMsg = DBManagerInst.Update_UserInfo(this.UserInfo)
	if this.ErrorMsg == nil {
		this.ErrorMsg = DBManagerInst.Delete_Friend(this.FriendInfo)
	}

	return this.ErrorMsg
}

func (this *FriendDeleteOper) DecodeContext(vals url.Values) {
	this.FriendBaseOperation.DecodeContext(vals)
}
