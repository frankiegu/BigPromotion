package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"myproject/service"
	"myproject/util"
	"strconv"
)

type GetHomePageInfoController struct {
	beego.Controller
}

func (p *GetHomePageInfoController) GetHomePageInfo() {

	data := make(map[string]interface{})

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
	}

	allRaiseCourseTimeNums := homepageService.GetAllRaiseNums(uinPublicDayDistributeLock)
	fmt.Println("allRaiseCourseTimeNums: ", allRaiseCourseTimeNums)
	realRaiseCourseTime, _ := json.Marshal(allRaiseCourseTimeNums)
	fmt.Println(string(realRaiseCourseTime))

	//etcd
	allConfigCourseNums := homepageService.GetConfigAllDonateCourseNums()
	fmt.Println("allConfigCourseNums: ", allConfigCourseNums)

	realRaiseCourseTimeNums, _ := strconv.ParseInt(string(realRaiseCourseTime), 10, 64)
	configCourseTimeNums,_ := strconv.ParseInt(allConfigCourseNums, 10, 64)

	allRaiseCourses := realRaiseCourseTimeNums + configCourseTimeNums

	data["CollectedTime"] = allRaiseCourses
	result["data"] = data
	return


}