package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"myproject/service"
	"myproject/util"
)

type GetHomePageInfoController struct {
	beego.Controller
}

func (p *GetHomePageInfoController) GetHomePageInfo() {



	uid, err := p.GetInt64("uid")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		result["code"] = service.ErrUserNull
		result["message"] = "invalid user_id"
		return
	}


	uinPublicDayDistributeLock := util.PUBLIC_DAY_DISTRIBUTE_LOCK
	uinPublicDayDistributeLock += fmt.Sprintf("%d", uid)

	homepageService := &service.HomePageService{}
	if uid != 0 {

		homepageService.AddUser(uid, uinPublicDayDistributeLock)
		//addUser(uid, uinPublicDayDistributeLock)
	}



}