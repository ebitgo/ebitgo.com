package routers

import (
	"github.com/astaxie/beego"
	"github.com/ebitgo/ebitgo.com/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/account", &controllers.AccountController{})
	beego.Router("/recaptcha", &controllers.ReCaptchaController{})
	beego.Router("/awallet", &controllers.WalletController{})
}
