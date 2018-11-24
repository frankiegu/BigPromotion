package controllers

import "github.com/beego"

type TestController struct {
	beego.Controller
}

func (c *TestController) Get() {
	c.Ctx.WriteString("this is the first beego controller get function test")
}

func (c *TestController) Post() {
	c.Ctx.WriteString("this is the first beego controller post function test")
}