package models

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type PostErrorInfo struct {
	ErrorMsg error
}

func (this *PostErrorInfo) checkValue(val, flag string) (err error) {
	if len(val) == 0 {
		this.ErrorMsg = errors.New("Error format [ Require parameter \"" + flag + "\" ]!")
		logmodels.Logger.Error(logmodels.SPrintError("PostErrorInfo : checkValue", "%v", this.ErrorMsg))
		err = this.ErrorMsg
	}
	return
}

func (this *PostErrorInfo) match_checkValue(val, flag string) (err error) {
	if len(val) == 0 {
		this.ErrorMsg = errors.New("Error format [ Require parameter \"" + flag + "\" ]!")
	} else {
		check := strings.ToLower(val)
		str := `(?:=)|(?:\$)|(?:\*)|(?::)|(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute|admin|administrator|root|id)\b)`
		re, err := regexp.Compile(str)
		if err != nil {
			this.ErrorMsg = err
		}
		if re.MatchString(val) {
			this.ErrorMsg = errors.New("Error format [  \"" + flag + "\" has invalid characters \"" + re.FindString(check) + "\"]!")
		}
	}

	if this.ErrorMsg != nil {
		logmodels.Logger.Error(logmodels.SPrintError("PostErrorInfo : checkValue", "%v", this.ErrorMsg))
		err = this.ErrorMsg
	}
	return
}

type AccountBaseInfo struct {
	PostErrorInfo
	LoginName string
	IsExist   bool
	UserInfo  *dbmodels.UserDbT
}

func (this *AccountBaseInfo) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_USEREMAIL: this.LoginName,
			RT_M_EXIST:     this.IsExist,
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

func (this *AccountBaseInfo) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 检查数据库
	this.IsExist = false
	uid, e := DBManagerInst.Quary_Check_Email_Exist(this.LoginName)
	if e == nil {
		this.IsExist = uid > 0

		if this.IsExist {
			this.UserInfo, e = DBManagerInst.Quary_UserInfo(uid)
		}
	}
	this.ErrorMsg = e
	return e
}

func (this *AccountBaseInfo) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_CHECK_USERNAME)
	this.match_checkValue(this.LoginName, POST_MARK_CHECK_USERNAME)
}
