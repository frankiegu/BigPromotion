package controllers

import (
	"fmt"
	"github.com/beego"
	"github.com/beego/logs"
	"myproject/service"
	"strings"
	"time"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	//p.Data["json"] = "sec Kill"
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	fmt.Println("productId: ", productId)
	fmt.Println("err: ", err)
	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid product_id"
		return
	} else {
		source := p.GetString("src")
		authcode := p.GetString("authcode")
		secTime := p.GetString("time")
		nance := p.GetString("nance")

		secRequest := service.NewSecRequst()
		secRequest.AuthCode = authcode
		secRequest.Nance = nance
		secRequest.ProductId = productId
		secRequest.SecTime = secTime
		secRequest.Source = source
		//secRequest.UserId,_ = strconv.Atoi(p.Ctx.GetCookie("UserId"))
		secRequest.UserId, _ = p.GetInt("user_id")
		secRequest.UserAuthSign = p.Ctx.GetCookie("userAuthSign")
		secRequest.AccessTime = time.Now()

		//userKey := fmt.Sprintf("%s_%s", secRequest.UserId, secRequest.ProductId)




		//10.108.42.200:39098 p.Ctx.Request.RemoteAddr
		if len(p.Ctx.Request.RemoteAddr) > 0 {
			secRequest.ClientAddr = strings.Split(p.Ctx.Request.RemoteAddr, ":")[0]
		}

		secRequest.ClientReference = p.Ctx.Request.Referer()
		//secRequest.CloseNotify = p.Ctx.ResponseWriter.CloseNotify()

		logs.Debug("client request:[%v]", secRequest)



		if err != nil {
			result["code"] = service.ErrInvalidRequest
			result["message"] = fmt.Sprintf("invalid cookie:userId")
			return
		}

		data, code, err := service.SecKill(secRequest)




		if err != nil {
			result["code"] = code
			result["message"] = err.Error()
			return
		}

		result["data"] = data
		result["code"] = code
		return

	}



}

func (p *SkillController) SecInfo() {
	//p.Data["json"] = "sec info"
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()
	if err != nil {

		data,code,err := service.SecInfoList()
		if err != nil {
			result["code"] = code
			result["message"] = err.Error
			logs.Error("Invalid request, get product_id failed, err:%v", err)
			return
		}
		result["code"] = code
		result["data"] = data
		return

	} else {
		data, code, err := service.SecInfo(productId)
		if err != nil {
			result["code"] = code
			result["message"] = err.Error
			logs.Error("invalid request, get product_id failed,err:%v", err)
			return
		}
		result["code"] = code
		result["data"] = data
	}




}

