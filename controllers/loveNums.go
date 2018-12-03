package controllers

import (
	"fmt"
	"github.com/beego"
	"myproject/service"
	"myproject/util"
)

type GetLoveNumsController struct {
	beego.Controller
}

func (p *GetLoveNumsController) GetLoveNums() {

	//data := make(map[string]interface{})

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

	loveNUmsService := &service.LoveNumsService{}

	curLoveNums := loveNUmsService.GetLoveNums(uid, uinPublicDayDistributeLock)
	fmt.Println("curLoveNums: ", curLoveNums)

}
