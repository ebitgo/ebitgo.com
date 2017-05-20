package models

import (
	"errors"
	// "fmt"
	"net/url"

	"github.com/ebitgo/ebitgo.com/models/logger"
)

type ReslutOperation struct {
	Error           interface{} `json:"error"`
	Data            interface{} `json:"data"`
	OperationResult IOperationInterface
}

func (this *ReslutOperation) GetResultData() *ReslutOperation {
	if this.OperationResult != nil {
		return this.OperationResult.GetResultData()
	}

	this.Error = map[string]interface{}{
		RT_M_ERROR: "Post message is undefined",
	}
	return this
}

func (this *ReslutOperation) QuaryExecute() (err error) {
	if this.OperationResult == nil {
		err = errors.New("Can not find out result object!")
	} else {
		err = this.OperationResult.QuaryExecute()
	}
	return
}

func (this *ReslutOperation) DecodeContext(vals url.Values) {
	this.OperationResult = nil
	// fmt.Println(vals)

	postType := vals.Get(POST_TYPE_FLAG)

	if len(postType) == 0 {
		return
	}

	switch postType {
	case PT_CHECK_USERNAME:
		this.OperationResult = &AccountBaseInfo{}

	case PT_USER_REGISTRATION:
		this.OperationResult = &AccountVerification{}

	case PT_USER_LOGIN:
		this.OperationResult = &UserLoginOperation{}

	case PT_CHECK_USERAUTH:
		this.OperationResult = &AccountAuthVerify{}

	case PT_USER_ACTIVE:
		this.OperationResult = &AccountActivation{}

	case PT_USER_MODIFY_PW:
		this.OperationResult = &AccountModifyPassword{}

	case PT_CHECK_GA_OPERA:
		this.OperationResult = &AccountUserGAConfig{}

	case RT_USER_ADD_WALLET_ACCOUNT:
		this.OperationResult = &WalletAddOper{}

	case PT_GET_WALLETS:
		this.OperationResult = &WalletQuaryAllOper{}

	case RT_UPDATE_WALLET_INFO:
		this.OperationResult = &WalletUpdateOper{}

	case RT_USER_DEL_WALLET_ACCOUNT:
		this.OperationResult = &WalletDeleteOper{}

	case RT_UPDATE_WALLET_SSKEY:
		this.OperationResult = &WalletUpdateSKEYOper{}

	case PT_SEARCH_WALLETS:
		this.OperationResult = &FriendSearchOperation{}

	case RT_USER_ADD_FRIEND:
		this.OperationResult = &FriendAddOperation{}

	case PT_GET_FRIENDS:
		this.OperationResult = &FriendQuaryAllOper{}

	case RT_USER_DEL_FRIEND:
		this.OperationResult = &FriendDeleteOper{}

	default:
		err := errors.New("Operation type is not defined!")
		logmodels.Logger.Error(logmodels.SPrintError("ReslutOperation : DecodeContext", "%v", err))
		return
	}
	this.OperationResult.DecodeContext(vals)
}
