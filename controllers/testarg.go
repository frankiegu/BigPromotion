package controllers

import "github.com/beego"


type TestArgController struct {
	beego.Controller
}

type User1 struct {
	Username string
	Age string
	Email string
}

//http://localhost:8081/test/arg?id=10&name=zhangsan
func (c *TestArgController) Get() {
	id := c.GetString("id")
	c.Ctx.WriteString(id)

	name := c.Input().Get("name")
	c.Ctx.WriteString(name)
}




//func (c *TestArgController) Get() {
//	c.TplName = "index.tpl"
//}

func (c *TestArgController) Post() {
	u := User1{}
	if err := c.ParseForm(&u); err!=nil {

	}
	c.Ctx.WriteString("Username:" + u.Username + ", Age:" + u.Age + ", Email:" + u.Email)
}