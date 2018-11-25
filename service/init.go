package service

import (
	"strconv"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/beego/logs"
)



func initRedisProcessFunc()(err error) {
	for i:=0; i< secKillConf.WriteProxy2LayerGoroutineNume; i++ {
		go WriteHandle()
	}

	for i:=0; i< secKillConf.ReadProxy2LayerGoroutineNums; i++ {
		go ReadHandle()
	}

	return
}




func loadBlackList()(err error) {

	err = initBlackList();
	if err != nil{
		logs.Error("init black redis failed, err:%v", err)
		return
	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", "idblacklist")
	idlist, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}

	for _, v := range idlist {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id[%v]", id)
			continue
		}

		secKillConf.idBlackMap[id] = true
	}


	reply, err = conn.Do("hgetall", "ipblacklist")
	iplist, err := redis.Strings(reply, err)
	if err != nil {
		logs.Warn("hget all failed, err:%v", err)
		return
	}

	for _, v := range iplist {
		secKillConf.ipBlackMap[v] = true
	}

	go SyncIdBlackList()
	go SyncIpBlackList()

	return
}



func initProxy2LayerRedis()(err error) {
	secKillConf.proxy2LayerRedisPool = &redis.Pool{
		MaxIdle:secKillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:secKillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout:time.Duration(secKillConf.RedisProxy2LayerConf.RedisIdleTimeout)*time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisProxy2LayerConf.RedisAddr)
		},
	}

	conn := secKillConf.proxy2LayerRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}


func initBlackList()(err error) {
	secKillConf.blackRedisPool = &redis.Pool{
		MaxIdle:secKillConf.RedisBlackConf.RedisMaxIdle,
		MaxActive:secKillConf.RedisBlackConf.RedisMaxActive,
		IdleTimeout:time.Duration(secKillConf.RedisBlackConf.RedisIdleTimeout)*time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisBlackConf.RedisAddr)
		},
	}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
		return
	}

	return
}

func SyncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()

	for {
		conn :=secKillConf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackIplist", time.Second)
		ip, err := redis.String(reply, err)

		if err != nil {
			continue
		}
		curTime := time.Now().Unix()
		ipList = append(ipList, ip)

		if len(ipList) > 200 || curTime - lastTime > 100 {
			secKillConf.RwSecProductLock.Lock()
			for _,v := range ipList {
				secKillConf.ipBlackMap[v] = true
			}
			secKillConf.RwSecProductLock.Unlock()

			lastTime = curTime
			logs.Info("sync ip list from redis succ, ip(%v)", ipList)

		}

	}
}

func SyncIdBlackList() {
	var idList []string
	lastTime := time.Now().Unix()

	for {
		conn :=secKillConf.blackRedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "blackIdlist", time.Second)
		id, err := redis.String(reply, err)

		if err != nil {
			continue
		}
		curTime := time.Now().Unix()
		idList = append(idList, id)

		if len(idList) > 200 || curTime - lastTime > 100 {
			secKillConf.RwSecProductLock.Lock()
			for _,v := range idList {
				id, err := strconv.Atoi(v)
				if err != nil {
					logs.Warn("invalid user id[%v]", id)
					continue
				}
				secKillConf.idBlackMap[id] = true
			}
			secKillConf.RwSecProductLock.Unlock()

			lastTime = curTime
			logs.Info("sync id list from redis succ, id(%v)", idList)

		}

	}
}