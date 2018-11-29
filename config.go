package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/logs"
	"myproject/service"
	"strings"
)



var(
	secKillConf = &service.SecKillConf{
		SecProductInfoMap:make(map[int]*service.SecProductInfoConf, 1024),
	}

	codisConfig = &service.CodisConf{

	}
)





func initConfig(err error) {

	redisBlackAddr := beego.AppConfig.String("redis_black_addr")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	logs.Debug("read config succ, redis addr:%v", redisBlackAddr)
	logs.Debug("read config succ, etcd addr:%v", etcdAddr)

	secKillConf.RedisBlackConf.RedisAddr = redisBlackAddr
	secKillConf.EtcdConf.EtcdAddr = etcdAddr

	if len(redisBlackAddr) == 0 || len(etcdAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s] or etcd[%s] config is null", redisBlackAddr, etcdAddr)
		return
	}

	redisMaxIdle, err := beego.AppConfig.Int("redis_black_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_idle error: %v", err)
	}

	secKillConf.RedisBlackConf.RedisMaxIdle = redisMaxIdle

	redisMaxActive, err := beego.AppConfig.Int("redis_black_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_active:%v", err)
		return
	}

	redisIdleTimeOut, err := beego.AppConfig.Int("redis_black_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_idle_timeout:%v", err)
		return
	}

	secKillConf.RedisBlackConf.RedisMaxIdle = redisMaxIdle
	secKillConf.RedisBlackConf.RedisMaxActive = redisMaxActive
	secKillConf.RedisBlackConf.RedisIdleTimeout = redisIdleTimeOut

	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}

	secKillConf.EtcdConf.Timeout = etcdTimeout

	secKillConf.EtcdConf.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if len(secKillConf.EtcdConf.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_key_prefix error:%v", err)
		return
	}

	addCourseNum := beego.AppConfig.String("etcd_add_course_key")
	secKillConf.EtcdConf.EtcdAddCourseKey = addCourseNum

	productKey := beego.AppConfig.String("etcd_product_key")
	if len(secKillConf.EtcdConf.EtcdSecProductKey) == 0 {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
	}

	if strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}


	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillConf.EtcdConf.EtcdSecKeyPrefix, productKey)

	secKillConf.LogPath = beego.AppConfig.String("log_path")
	secKillConf.LogLevel = beego.AppConfig.String("log_level")

	secKillConf.CookieSecretKey = beego.AppConfig.String("cookie_secretkey")
	secKillConf.UserSecAccessLimit, _ = beego.AppConfig.Int("user_sec_access_limit")
	secKillConf.IpSecAccessLimit, _ = beego.AppConfig.Int("ip_sec_access_limit")

	referList := beego.AppConfig.String("refer_whitelist")
	if len(referList) > 0 {
		secKillConf.ReferWhiteList = strings.Split(referList, ",")
	}

	ipLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit error:%v", err)
		return
	}

	secKillConf.IpSecAccessLimit = ipLimit
	//

	redisProxy2LayerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	logs.Debug("read config succ, redis addr:%v", redisProxy2LayerAddr)

	secKillConf.RedisProxy2LayerConf.RedisAddr = redisProxy2LayerAddr

	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("init config failed, redis[%s] config is null", redisBlackAddr)
		return
	}

	redisProxy2LayerMaxIdle, err := beego.AppConfig.Int("redis_proxy2layer_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_max_idle error: %v", err)
	}

	secKillConf.RedisProxy2LayerConf.RedisMaxIdle = redisProxy2LayerMaxIdle

	redisProxy2LayerMaxActive, err := beego.AppConfig.Int("redis_proxy2layer_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_max_active:%v", err)
		return
	}

	redisProxy2LayerIdleTimeOut, err := beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle_timeout:%v", err)
		return
	}

	secKillConf.RedisProxy2LayerConf.RedisMaxIdle = redisProxy2LayerMaxIdle
	secKillConf.RedisProxy2LayerConf.RedisMaxActive = redisProxy2LayerMaxActive
	secKillConf.RedisProxy2LayerConf.RedisIdleTimeout = redisProxy2LayerIdleTimeOut

	writeGoNums, err := beego.AppConfig.Int("write_proxy2layer_goroutine_nume")
	if err != nil {
		err = fmt.Errorf("init config failed, read write_proxy2layer_goroutine_nume:%v", err)
		return
	}
	secKillConf.WriteProxy2LayerGoroutineNume = writeGoNums

	ReadGoNums, err := beego.AppConfig.Int("read_layer2proxy_goroutine_nume")
	if err != nil {
		err = fmt.Errorf("init config failed, read read_layer2proxy_goroutine_nume:%v", err)
		return
	}
	secKillConf.ReadProxy2LayerGoroutineNums = ReadGoNums





	return
}