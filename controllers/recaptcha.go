package controllers

import (
	"github.com/astaxie/beego"
	"github.com/ebitgo/ebitgo.com/models"
	"github.com/ebitgo/ebitgo.com/models/logger"
)

type ReCaptchaController struct {
	beego.Controller
}

func (this *ReCaptchaController) Get() {
	this.Data["Website"] = "eBitGo.com"
	this.Data["Email"] = "QQ群: 452779719"
	this.TplName = "index.tpl"
}

func (this *ReCaptchaController) Post() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Accept")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
	equation := this.Input().Get("Equation")

	// fmt.Println("Equation = ", equation)
	logmodels.Logger.Trace(logmodels.SPrintTrace("ReCaptchaController : Post", "Equation = %s", equation))

	if equation == "quary" {
		equa := &models.EquationStruct{}
		this.Data["json"] = equa.CreateEquation().GetResultData()
		this.ServeJSON()
	}
}
