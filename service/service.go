package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/beego/logs"
	"time"
)



const (
	ProductStatusNamal = 0
	ProductStatusSaleOut = 1
	ProductStatusForceSaleOut = 2
)

func NewSecRequst()(secRequst *SecRequst) {
	//create an object
	secRequst = &SecRequst{
		ResultChan: make(chan *SecResult, 1),
	}
	return
}


func SecInfo(productId int)(data []map[string]interface{}, code int, err error) {

	secKillConf.RwSecProductLock.Lock()
	defer secKillConf.RwSecProductLock.Unlock()
	logs.Debug("sec info config is [%v]", secKillConf.SecProductInfoMap[productId])
	item, code, err := SecInfoById(productId)
	data = append(data, item)
	return
}


/*
	data: map array
 */
func SecInfoList()(data []map[string]interface{}, code int, err error){

	for _,v := range secKillConf.SecProductInfoMap {
		secKillConf.RwSecProductLock.Lock()
		item,_,err := SecInfoById(v.ProductId)
		secKillConf.RwSecProductLock.Unlock()
		if (err != nil) {
			logs.Error("get product_id[%d] failed, err:%v", v.ProductId, err)
			continue
		}
		data = append(data, item)
	}

	return
}


func SecInfoById(productId int)(data map[string]interface{}, code int, err error) {

	v, ok := secKillConf.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}

	start := false
	end := false
	status := "sec kill is already start and not end"

	now := time.Now().Unix()
	if now - v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill is not start"
	}

	if now - v.StartTime > 0 {
		start = true
	}

	if now - v.EndTime > 0 {
		start = false
		end = true
		status = "sec kill is already end"
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "product is sale out"
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start_time"] = start
	data["end_time"] = end
	data["status"] = status

	return
}

func userCheck(req *SecRequst)(err error){

	found := false
	for _,refer := range  secKillConf.ReferWhiteList {

		if req.ClientReference == "" {
			// temporary pass the userCheck
			found = true
			break
		}

		if refer == req.ClientReference {
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("invalid request")
		logs.Warn("user[%d] is reject by refer, req[%v]", req.UserId, req)
	}

	authData := fmt.Sprintf("%d:%s", req.UserId, secKillConf.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	//temporary pass the cookie check
	if req.UserAuthSign == "" {
		return
	}

	if authSign != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
	}
	return
}

func SecKill(req *SecRequst)(data map[string]interface{}, code int, err error) {

	secKillConf.RwSecProductLock.Lock()
	defer secKillConf.RwSecProductLock.Unlock()

	//check user
	err = userCheck(req)
	if err != nil {
		code = ErrUserCheckAuthFailed
		logs.Warn("userId[%d] invalid,check failed, req[%v]", req.UserId, req)
		return
	}

	//anti-brush
	err = antiSpam(req)
	data = make(map[string]interface{})
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userId[%d] invalid,check failed, req[%v]", req.UserId, req)
		data["product_id"] = "failed"
		fmt.Println("invalid request, because busy")
		return
	}

	data, code, err = SecInfoById(req.ProductId)
	//map to json
	str, _ := json.Marshal(data)
	println("sec info data: ", string(str))

	if err != nil {
		logs.Warn("user id[%d] secInfoBy id failed, req[%v]", req.UserId, req)
		return
	}

	if code != 0 {
		logs.Warn("userId[%d] secInfoById failed, code[%d] req[%v]", req.UserId, code, req)
		return
	}


	userKey := fmt.Sprintf("%s_%s", req.UserId, req.ProductId)

	secKillConf.UserConnMap[userKey] = req.ResultChan

	//write to redis
	secKillConf.secReqChan <- req

	//wait 10s
	ticker := time.NewTicker(time.Second * 2)

	defer func() {
		ticker.Stop()
		secKillConf.UserConnMapLock.Lock()
		delete(secKillConf.UserConnMap,  userKey)
		secKillConf.UserConnMapLock.Unlock()
	}()

	select {

		case <- ticker.C:
			code = ErrProcessTimeout
			err = fmt.Errorf("request timeout")
			return
		case <- req.CloseNotify:
			code = ErrClientClosed
			err = fmt.Errorf("client already closed")
			return
		case result := <-req.ResultChan:
			println("wocao........")
			code = result.Code
			data["product_id"] = result.ProductId
			data["token"] = result.Token
			data["user_id"] = result.UserId
			return

	}
	return

}