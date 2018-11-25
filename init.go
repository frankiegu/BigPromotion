package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beego/logs"
	etcd_client "github.com/coreos/etcd/client"
	"github.com/garyburd/redigo/redis"
	"log"
	"myproject/service"
	"strings"
	"time"
)

var (
	redisPool * redis.Pool
	etcdClient * etcd_client.Client
)

func initSec(err error) {

	//log
	err = initLogger()
	if err != nil {
		logs.Error("init logger failed, err:%v", err)
		return
	} else {
		logs.Info("init logger success")
	}

	//etcd
	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed, err:%v", err)
		return
	} else {
		logs.Info("init etcd success")
	}


	//redis
	err = initRedis()
	if err != nil {
		logs.Error("init redis failed, err:%v", err)
		return
	} else {
		logs.Info("init redis success")
	}

	//get all product info
	err = loadSecConf()
	if err != nil {
		logs.Error("init sec conf failed, err:%v", err)
		return
	} else {
		logs.Info("init sec success")
	}

	logs.Info("initSecProductWatcher begin")

	/*
	include blacklist
	*/

	service.InitService(secKillConf)

	initSecProductWatcher()
	//Test_watch()
	return

}

func Test_watch() {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints : []string{secKillConf.EtcdConf.EtcdAddr},
		HeaderTimeoutPerRequest : time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed, err", err)
	}

	kapi := etcd_client.NewKeysAPI(cli)

	go func() {
		watcher := kapi.Watcher("/a", &etcd_client.WatcherOptions{
			Recursive: true,
		})
		index := 0
		for {
			resp, err := watcher.Next(context.Background())
			fmt.Println("happend action", index,":", resp.Action, "\nkey:", resp.Node.Key, "\nvalue:", resp.Node.Value, "\nerr:", err)

			index += 1;
		}
	}()


	time.Sleep(1 * time.Second)
	kapi.Set(context.Background(), "/a/b", "a-b", &etcd_client.SetOptions{})
	time.Sleep(1 * time.Second)
	kapi.Set(context.Background(), "/a/b/c", "a-b-c", &etcd_client.SetOptions{})
	time.Sleep(1 * time.Second)
	kapi.Set(context.Background(), "/a/d", "a-d", &etcd_client.SetOptions{})
	time.Sleep(1 * time.Second)
	kapi.Delete(context.Background(), "/a/d", &etcd_client.DeleteOptions{})
	time.Sleep(1 * time.Second)
	kapi.Delete(context.Background(), "/a/b/c", &etcd_client.DeleteOptions{})
	time.Sleep(1 * time.Second)
	kapi.Delete(context.Background(), "/a/b", &etcd_client.DeleteOptions{})
	time.Sleep(1 * time.Second)


}

//monitor etcd change
func initSecProductWatcher() {


	kapi := getKapi()
	//coroutine
	go watchSecProductKey(kapi)


}
func getKapi() etcd_client.KeysAPI {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints : []string{secKillConf.EtcdConf.EtcdAddr},
		HeaderTimeoutPerRequest : time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed, err", err)
		return nil
	}

	kapi := etcd_client.NewKeysAPI(cli)
	return kapi
}
func watchSecProductKey(kapi etcd_client.KeysAPI) {
	key := secKillConf.EtcdConf.EtcdSecProductKey
	watcher := kapi.Watcher(key, &etcd_client.WatcherOptions{
		Recursive:true,
	})
	logs.Debug("begin function watchSecProductKey")
	var secProductInfo []service.SecProductInfoConf
	var getConfSucc = true

	for {
		resp, err := watcher.Next(context.Background())
		if strings.EqualFold(resp.Action, "delete") {
			logs.Warn("key[%s] config deleted", key)
			continue
		} else if strings.EqualFold(resp.Action, "put") && strings.EqualFold(key, resp.Node.Key){
			err = json.Unmarshal([]byte(resp.Node.Value), &secProductInfo)
			if err != nil {
				logs.Error("key [%s], unmarshal[%s],err:%v", resp.Node.Key, resp.Node.Value, err)
				getConfSucc = false
				continue
			}
		}
		logs.Debug("get config from etcd, %s %q:%q\n", resp.Action, resp.Node.Key, resp.Node.Value)


	}

	if getConfSucc {
		logs.Debug("get config from etcd succ, %v", secProductInfo)
		updateSecProductInfo(secProductInfo)
	} else {
		logs.Error("get conf error in function watchSecProductKey")
	}

}
func updateSecProductInfo(secProductInfoConf []service.SecProductInfoConf) {

	var tmp map[int]*service.SecProductInfoConf = make(map[int]*service.SecProductInfoConf, 1024)
	for _, v := range secProductInfoConf {
		ProductInfo := v
		tmp[v.ProductId] = &ProductInfo
	}

	secKillConf.RwSecProductLock.Lock()
	secKillConf.SecProductInfoMap = tmp
	secKillConf.RwSecProductLock.Unlock()



}



func convertLogLevel(level string) int {
	switch (level) {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelDebug
}

func initLogger() (err error) {

	config := make(map[string]interface{})
	config["filename"] = secKillConf.LogPath
	config["level"] = convertLogLevel(secKillConf.LogLevel)

	configStr, err := json.Marshal(config)

	if err != nil {
		fmt.Println("Marshal failed, err:%v", err)
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

//etcdSecKeyPredix:/z/seckill
//productSecKey:/z/seckill/product
//blackSecKey:/z/seckill/blackListKey


func loadSecConf()(err error) {
	kapi := etcd_client.NewKeysAPI(*etcdClient)
	resp, err := kapi.Get(context.Background(), secKillConf.EtcdConf.EtcdSecProductKey, nil)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", secKillConf.EtcdConf.EtcdSecProductKey, err)
		return
	}

	var secProductInfo []service.SecProductInfoConf
	err = json.Unmarshal([]byte(resp.Node.Value), &secProductInfo)

	if err != nil {
		logs.Error("Unmarshal sec product info failed, err:%v", err)
	}

	log.Printf("\n  Unmarshal sec product info succeed, \n  key: %q \n, value: %q\n", resp.Node.Key, resp.Node.Value)
	updateSecProductInfo(secProductInfo)

	return err
}

func initEtcd() (err error) {

	logs.Info("secKillConf.etcdConf.etcdAddr : " + secKillConf.EtcdConf.EtcdAddr)
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints : []string{secKillConf.EtcdConf.EtcdAddr},
		HeaderTimeoutPerRequest : time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed, err: ", err)
		return
	}
	//kapi := etcd_client.NewKeysAPI(cli)


	etcdClient = &cli
	return nil
}



func initRedis() (err error) {

	redisPool = &redis.Pool{
		MaxIdle: secKillConf.RedisBlackConf.RedisMaxIdle,
		MaxActive:secKillConf.RedisBlackConf.RedisMaxActive,
		IdleTimeout:time.Duration(secKillConf.RedisBlackConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisBlackConf.RedisAddr)
		},
	}

	//redisPool = &redis.Pool{
	//	MaxIdle:64,
	//	MaxActive:0,
	//	IdleTimeout:300,
	//	Dial: func() (redis.Conn, error) {
	//		return redis.Dial("tcp", "127.0.0.1")
	//	},
	//}

	conn := redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed, err:%v", err)
	}

	return
}
