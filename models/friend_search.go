package models

import (
	"errors"
	"net/url"
	"time"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type FriendSearchOperation struct {
	WalletBaseInfo
	SearchResult map[string][]*dbmodels.WalletDbT
}

func (this *FriendSearchOperation) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:     true,
			RT_M_AUTH_UPDATE: false,
		}
		if len(this.LoginName) > 0 {
			ret.Data.(map[string]interface{})[RT_M_AUTH_UPDATE] = true
			ret.Data.(map[string]interface{})[RT_M_USER_AUTH] = this.UserInfo.AuthStr
		}
		results := this.SearchResult["RESULT"]
		if results == nil {
			ret.Data.(map[string]interface{})[RT_M_RESULTS] = make([]interface{}, 0)
		} else {
			rlen := len(results)
			resultArray := make([]interface{}, rlen)
			for i := 0; i < rlen; i++ {
				resultArray[i] = map[string]interface{}{
					RT_M_WALLET_ID:         results[i].Id,
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

func (this *FriendSearchOperation) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 检查登录用户是否有效
	if len(this.LoginName) > 0 {
		err := this.AccountBaseInfo.QuaryExecute()
		if err != nil {
			this.ErrorMsg = err
			return err
		}

		if !this.IsExist { // 用户不存在
			err = errors.New("User not exist!")
			logmodels.Logger.Error(logmodels.SPrintError("FriendSearchOperation : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		} else if this.LogAuth != this.UserInfo.AuthStr {
			err = errors.New("Login user timeout or have already login in other places!")
			logmodels.Logger.Error(logmodels.SPrintError("FriendSearchOperation : QuaryExecute", "%v", err))
			this.ErrorMsg = err
			return err
		}

		this.WalletInfo.UserId = this.UserInfo.Id
	}

	this.SearchResult, this.ErrorMsg = DBManagerInst.Search_WalletInfo(this.WalletInfo.PublicAddr, this.WalletInfo.NickName)

	if this.ErrorMsg == nil && len(this.LoginName) > 0 {
		timenow := time.Now()

		this.UserInfo.AuthStr = getSHA256Key(timenow.Format("2006-01-02 15:04:05"),
			this.LoginName, this.UserInfo.Password, "")

		// fmt.Println("LastLoginTime = ", uInfo.LastLoginTime.Format("2006-01-02 15:04:05"))
		this.UserInfo.LastLoginTime = timenow
		this.ErrorMsg = DBManagerInst.Update_UserInfo(this.UserInfo)
	}
	return this.ErrorMsg
}

func (this *FriendSearchOperation) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_USER_NAME)
	if len(this.LoginName) > 0 {
		if this.match_checkValue(this.LoginName, POST_MARK_USER_NAME) != nil {
			return
		}

		this.LogAuth = vals.Get(POST_MARK_AUTHCODE)
		if this.checkValue(this.LogAuth, POST_MARK_AUTHCODE) != nil {
			return
		}

	}
	this.WalletInfo = &dbmodels.WalletDbT{}

	this.WalletInfo.NickName = vals.Get(POST_MARK_WALLET_NICKNAME)
	this.WalletInfo.PublicAddr = vals.Get(POST_MARK_WALLET_PUBADDR)

	if len(this.WalletInfo.NickName) == 0 && len(this.WalletInfo.PublicAddr) == 0 {
		this.ErrorMsg = errors.New("Nick name and public address are null!")
		return
	}
	if len(this.WalletInfo.NickName) > 0 {
		this.match_checkValue(this.WalletInfo.NickName, POST_MARK_WALLET_NICKNAME)
	}
}
