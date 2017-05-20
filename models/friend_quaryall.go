package models

import (
	"errors"
	// "time"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type FriendQuaryAllOper struct {
	FriendBaseOperation
	Friends map[string][]*dbmodels.FriendlyShipDbT
}

func (this *FriendQuaryAllOper) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:   true,
			RT_M_USEREMAIL: this.LoginName,
		}
		results := this.Friends["RESULT"]
		if results == nil {
			ret.Data.(map[string]interface{})[RT_M_RESULTS] = make([]interface{}, 0)
		} else {
			rlen := len(results)
			resultArray := make([]interface{}, rlen)
			for i := 0; i < rlen; i++ {
				resultArray[i] = map[string]interface{}{
					RT_M_WALLET_ID:         results[i].WalletId,
					RT_M_WALLET_NICKNAME:   results[i].NickName,
					RT_M_WALLET_PUBLICADDR: results[i].PublicAddr,
				}
			}
			ret.Data.(map[string]interface{})[RT_M_RESULTS] = resultArray
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *FriendQuaryAllOper) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 检查登录用户是否有效
	this.ErrorMsg = this.AccountBaseInfo.QuaryExecute()
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	if !this.IsExist { // 用户不存在
		err := errors.New("User not exist!")
		logmodels.Logger.Error(logmodels.SPrintError("FriendQuaryAllOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	} else if this.LogAuth != this.UserInfo.AuthStr {
		err := errors.New("Login user timeout or have already login in other places!")
		logmodels.Logger.Error(logmodels.SPrintError("FriendQuaryAllOper : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	this.Friends, this.ErrorMsg = DBManagerInst.Quary_ALL_Friends(this.UserInfo.Id)

	return this.ErrorMsg
}

func (this *FriendQuaryAllOper) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
		return
	}

	this.LogAuth = vals.Get(POST_MARK_AUTHCODE)
	this.checkValue(this.LogAuth, POST_MARK_AUTHCODE)
}
