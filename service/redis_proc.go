package service

import (
	"encoding/json"
	"fmt"
	"github.com/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
)

func WriteHandle() {
	for{
		//req := <-secKillConf.secReqChan
		conn := secKillConf.proxy2LayerRedisPool.Get()
		for res := range secKillConf.secReqChan {

			println("res: ", )
			data, err := json.Marshal(res)
			if err != nil {
				println("cao......json.Marshal failed.")
				logs.Error("json.Marshal failed, error:%v, req:%v", err, res)
				conn.Close()
				continue
			}
			_, err = conn.Do("LPUSH", "sec_queue", string(data))
			if err != nil {
				println("cao.....lpush failed")
				logs.Warn("lPUSH failed, error:%v", err)
				conn.Close()
				continue
			}
			conn.Close()

		}
		conn.Close()


	}
}

func ReadHandle() {

	for{
		conn := secKillConf.proxy2LayerRedisPool.Get()
		reply, err := conn.Do("RPOP", "sec_queue")
		data, err := redis.String(reply, err)

		if err == redis.ErrNil {
			time.Sleep(time.Second)
			conn.Close()
			continue
		}

		logs.Debug("rpop from redis succ, data:%s", string(data))

		if err != nil {
			//logs.Error("RPOP failed, error:%v", err)
			conn.Close()
			continue
		}
		println("rpop from redis succ, data: ", string(data))

		var result SecResult
		err = json.Unmarshal([]byte(data), &result)

		if err != nil {
			logs.Error("json.Unmarshal failed, error:%v", err)
			conn.Close()
			continue
		}


		println("secKillConf.UserConnMap get data")
		for k, v := range secKillConf.UserConnMap {
			println("k: ", k)
			println("v:")
			println("usrid: ", v.UserId)
			println("productid: ", v.ProductId)
			println("code: ", v.Code)
			println("token: ", v.Token)
		}

		userKey := fmt.Sprintf("%s_%s", result.UserId, result.ProductId)
		println("userKey: ", userKey)
		secKillConf.UserConnMapLock.Lock()
		BasicInfo := secKillConf.UserConnMap[userKey]
		fmt.Println(BasicInfo)

		//simulate fill business data
		buzData := &SecResponse{}
		buzData.Token = BasicInfo.Token
		buzData.Code = BasicInfo.Code
		buzData.ProductId = BasicInfo.ProductId
		buzData.UserId = BasicInfo.UserId
		buzData.TokenTime = 1234

		businessResponse <- buzData

		//simulate get code
		BasicInfo.Code = 0
		secKillConf.UserConnMapLock.Unlock()


		//BasicInfo <- &result
		BasicInfo = result
		conn.Close()


	}
}