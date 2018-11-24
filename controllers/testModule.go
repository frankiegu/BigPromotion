package controllers

import (
	"github.com/astaxie/beego"
	"myproject/models"
)

type TestModuleController struct {
	beego.Controller
}



func (c *TestModuleController) Get() {

	//user := models.User{Name:"wangwu"}
	//
	//c.Ctx.WriteString("operate mysql ... \n")
	//
	//models.AddUser(&user)
	//
	//c.Ctx.WriteString("call mmodel success")

	var users []models.User
	models.ReadUser(&users)

	c.Data["Title"] = "title"
	c.Data["Users"] = users
	c.Data["len"] = len(users)
	c.TplName = "test.tpl"

}




