package models

import (
	// "errors"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ebitgo/ebitgo.com/models/databasemodels"
	"github.com/ebitgo/ebitgo.com/models/logger"
	"github.com/go-gomail/gomail"
)

type VerifiFormular struct {
	Value1   string
	Value2   string
	Value3   string
	Operator string
	ViCode   string
}

func (this *VerifiFormular) Verify() bool {
	if len(this.ViCode) < 1 {
		return false
	}

	var v1, v2, v3 int

	if this.Value1 == "?" {
		v1, _ = strconv.Atoi(this.ViCode)
	} else {
		v1, _ = strconv.Atoi(this.Value1)
	}

	if this.Value2 == "?" {
		v2, _ = strconv.Atoi(this.ViCode)
	} else {
		v2, _ = strconv.Atoi(this.Value2)
	}

	if this.Value3 == "?" {
		v3, _ = strconv.Atoi(this.ViCode)
	} else {
		v3, _ = strconv.Atoi(this.Value3)
	}

	if this.Operator == "-" {
		return (v1 - v2) == v3
	}
	return (v1 + v2) == v3
}

type AccountVerification struct {
	AccountBaseInfo
	VerifInfo *VerifiFormular
}

func (this *AccountVerification) GetResultData() *ReslutOperation {
	ret := &ReslutOperation{}

	if this.ErrorMsg == nil {
		ret.Data = map[string]interface{}{
			RT_M_SUCCESS:   true,
			RT_M_USEREMAIL: this.LoginName,
		}
	}
	if this.ErrorMsg != nil {
		ret.Error = map[string]interface{}{
			RT_M_ERROR: this.ErrorMsg.Error(),
		}
	}
	return ret
}

func (this *AccountVerification) QuaryExecute() error {
	if this.ErrorMsg != nil {
		return this.ErrorMsg
	}

	// 计算校验码 如果正确，信息入库
	if this.VerifInfo.Verify() == false {
		err := errors.New("ViCode is not valid!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountVerification : QuaryExecute", "%v", err))
		this.ErrorMsg = err
		return err
	}

	// 先检查有没有当前登录用户
	err := this.AccountBaseInfo.QuaryExecute()
	if err != nil {
		this.ErrorMsg = err
		return err
	}
	// 用户存在
	if this.IsExist {
		err = errors.New("User exist! Registration failure!")
		logmodels.Logger.Error(logmodels.SPrintError("AccountUserGAConfig : QuaryExecute", "%v", err, this.LoginName))
		this.ErrorMsg = err
		return err
	}

	this.UserInfo.ActiveFlag = randomString(20)

	err = DBManagerInst.Add_RegUser(this.LoginName, this.UserInfo)
	if err != nil {
		this.ErrorMsg = err
		return err
	}

	this.sendActiveEmail()

	return nil
}

func (this *AccountVerification) DecodeContext(vals url.Values) {
	this.LoginName = vals.Get(POST_MARK_REG_USERNAME)
	if this.match_checkValue(this.LoginName, POST_MARK_REG_USERNAME) != nil {
		return
	}

	this.UserInfo = new(dbmodels.UserDbT)
	this.UserInfo.Password = vals.Get(POST_MARK_REG_PASSWORD)
	if this.checkValue(this.UserInfo.Password, POST_MARK_REG_PASSWORD) != nil {
		return
	}

	if this.VerifInfo == nil {
		this.VerifInfo = &VerifiFormular{}
	}

	this.VerifInfo.ViCode = vals.Get(POST_MARK_REG_VICODE)
	if this.checkValue(this.VerifInfo.ViCode, POST_MARK_REG_VICODE) != nil {
		return
	}

	this.VerifInfo.Value1 = vals.Get(POST_MARK_REG_VALUE1)
	if this.checkValue(this.VerifInfo.Value1, POST_MARK_REG_VALUE1) != nil {
		return
	}

	this.VerifInfo.Value2 = vals.Get(POST_MARK_REG_VALUE2)
	if this.checkValue(this.VerifInfo.Value1, POST_MARK_REG_VALUE1) != nil {
		return
	}

	this.VerifInfo.Value3 = vals.Get(POST_MARK_REG_VALUE3)
	if this.checkValue(this.VerifInfo.Value3, POST_MARK_REG_VALUE3) != nil {
		return
	}

	this.VerifInfo.Operator = vals.Get(POST_MARK_REG_OPERA)
	this.checkValue(this.VerifInfo.Operator, POST_MARK_REG_OPERA)
}

func (this *AccountVerification) sendActiveEmail() error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "support@eBitGo.com")
	mail.SetHeader("To", this.LoginName)
	mail.SetHeader("Subject", "eBitGo.com : verify your account")
	formatStr := "Hello <strong>%s</strong>:<br>" +
		"<p>Your account login name is %s, " +
		"and you have to active this account, " +
		"after activation, your account will be useful!</p>" +
		"<p>你的登录名为 %s, 请尽快激活账户, 激活后你可以使用文明提供的所有功能!</p><br>" +
		"<p>Your acitvation : <a href=\"http://eBitGo.com/activeaccount.html\" target=\"_black\">Click to acitve</a></p>" +
		"<p>点击这里进行激活 : <a href=\"http://eBitGo.com/activeaccount.html\" target=\"_black\">去这里激活</a></p><br>" +
		"<p>Please keep <strong>fid</strong> and <strong>tid</strong> properly preserved, you can use when you forget the password or change the password.</p>" +
		"<p>请将“fid”和“tid”妥善保存，可以在忘记密码或修改密码时使用。</p><br>" +
		"<p>Your fid = <strong>%s</strong> </p>" +
		"<p>Your tid = <strong>%d</strong> </p><br>" +
		"<p>If you can not open , please open url : http://eBitGo.com/activeaccount.html </p>" +
		"<p>如果不能打开网页，请用浏览器直接访问 : http://eBitGo.com/activeaccount.html </p><br>" +
		"<p><i>Activate the link is valid for <strong>one day</strong>.</i></p>" +
		"<p><i>当前激活ID有效期为 <strong>1天</strong>.</i></p><br>" +
		"<p><i>Please do not reply to this message.</i></p>" +
		"<p><i>请不要回复此邮件。</i></p><br>" +
		"<p><h4>support@eBitGo.com</h4></p>" +
		"<p><h3>http://www.eBitGo.com</h3></p>"
	emailBody := fmt.Sprintf(formatStr,
		this.LoginName, this.LoginName, this.LoginName, this.UserInfo.ActiveFlag, this.UserInfo.Id)
	mail.SetBody("text/html", emailBody)
	// smtp.exmail.qq.com
	// 海外用户别名主机请设置为：hwsmtp.exmail.qq.com
	dailer := gomail.NewPlainDialer("smtp.mxhichina.com", 465, "support@eBitGo.com", "Clock1234")
	dailer.SSL = true
	if err := dailer.DialAndSend(mail); err != nil {
		logmodels.Logger.Info(logmodels.SPrintInfo("AccountVerification : sendActiveEmail", "%v", err))
		return err
	}
	return nil
}
