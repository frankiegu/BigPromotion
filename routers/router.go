package routers

import (
	"myproject/controllers"
	"github.com/astaxie/beego"
)

func init() {
    //beego.Router("/", &controllers.MainController{}, "*:Get;post:Test")
    beego.Router("/testC", &controllers.TestController{}, "get:Get;post:Post")
	beego.Router("/test/arg", &controllers.TestArgController{}, "get:Get;post:Post")
	beego.Router("/test/orm", &controllers.TestModuleController{}, "get:Get;post:Post")
	beego.Router("/test/http", &controllers.TestHttpLibController{}, "get:Get;post:Post")
    beego.Router("/secKill", &controllers.SkillController{}, "*:SecKill")
    beego.Router("/secInfo", &controllers.SkillController{}, "*:SecInfo")
    beego.Router("/getHomePageInfo", &controllers.GetHomePageInfoController{}, "*:GetHomePageInfo")

}
