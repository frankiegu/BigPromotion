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
		req := <-secKillConf.secReqChan
		conn := secKillConf.proxy2LayerRedisPool.Get()

		data, err := json.Marshal(req)
		//data == nil
		if err != nil {
			println("cao......json.Marshal failed.")
			logs.Error("json.Marshal failed, error:%v, req:%v", err, req)
			conn.Close()
			continue
		}

		_, err = conn.Do("LPUSH", "sec_que", string(data))
		if err != nil {
			println("cao.....lpush failed")
			logs.Error("lPUSH failed, error:%v, req:%v", err, req)
			conn.Close()
			continue
		}

		conn.Close()

	}
}

func ReadHandle() {

	for{
		conn := secKillConf.proxy2LayerRedisPool.Get()
		reply, err := conn.Do("RPOP", "recv_queue")
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

		var result SecResult
		err = json.Unmarshal([]byte(data), &result)

		if err != nil {
			logs.Error("json.Unmarshal failed, error:%v", err)
			conn.Close()
			continue
		}

		userKey := fmt.Sprint("%s_%s", result.UserId, result.ProductId)
		secKillConf.UserConnMapLock.Lock()
		resultChan, ok := secKillConf.UserConnMap[userKey]
		secKillConf.UserConnMapLock.Unlock()

		if !ok {
			conn.Close()
			logs.Warn("user not found:%v", err)
			continue
		}

		resultChan <- &result
		conn.Close()


	}
}