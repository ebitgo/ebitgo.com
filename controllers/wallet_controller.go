package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ebitgo/ebitgo.com/models"
)

type WalletController struct {
	beego.Controller
}

func (this *WalletController) Get() {
	this.Data["Website"] = "eBitGo.com"
	this.Data["Email"] = "QQ群: 452779719"

	// results := &models.ReslutOperation{}

	// results.DecodeContext(this.Input())
	// err := results.QuaryExecute()

	// if err == nil {
	// 	this.Data["dataresult"] = results.GetResultData()
	// } else {
	// 	this.Data["dataresult"] = map[string]interface{}{
	// 		"Error": err.Error(),
	// 	}
	// }
	this.TplName = "index.tpl"
}

func (this *WalletController) Post() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Accept")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")

	results := &models.ReslutOperation{}

	results.DecodeContext(this.Input())
	err := results.QuaryExecute()
	if err == nil {
		this.Data["json"] = results.GetResultData()
	} else {
		this.Data["json"] = map[string]interface{}{
			"Error": err.Error(),
		}
	}
	this.ServeJSON()
}
