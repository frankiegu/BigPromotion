package service

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"myproject/bean"
	"myproject/util"
	"strconv"
	"time"
	etcd_client "github.com/coreos/etcd/client"
	"context"

)

type AddCourseNums struct {
	AddNums int64
}

type HomePageService struct {

}
var(
	codisUtil = util.GetCodisUtilInstance()
	DefaultTimeout = 1

)

func (h *HomePageService) AddUser(uin int64, uinPublicDayDistributeLock string) {
	fmt.Println("uinPublicDayDistributeLock: ", uinPublicDayDistributeLock)
	fmt.Println("uin: ", uin)

	//codisUtil := util.GetCodisUtilInstance()
	redisKey := codisUtil.GetCodisKey(util.USER_DETAILS)
	fmt.Println("redisKey: ", redisKey)

	lock := &util.Lock{}
	conn, err := getRedisConn()
	defer conn.Close()

	res, err := lock.DoWithLock(uinPublicDayDistributeLock, DefaultTimeout, conn, AddUserFunc{}, redisKey, uin)
	fmt.Println("AddUser lock.DoWithLock res: ", res)
	fmt.Println("AddUser lock.DoWithLock err: ", err)

}

func getRedisConn() (redis.Conn, interface{}) {
	return redis.Dial("tcp", "localhost:6379")
}

func (h *HomePageService) GetAllRaiseNums(uinPublicDayDistributeLock string) interface{} {

	//codisUtil := util.GetCodisUtilInstance()
	redisKey := codisUtil.GetCodisKey(util.GLOBAL_DONATE_DETAILS)
	lock := &util.Lock{}
	conn, err := getRedisConn()
	defer conn.Close()

	res, err := lock.DoWithLock(uinPublicDayDistributeLock, DefaultTimeout, conn, GetAllRaiseNumsFunc{}, redisKey)
	if err != nil {
		res = nil
		fmt.Println("err: ", err)
	}
	return res
}
func (h *HomePageService) GetConfigAllDonateCourseNums() string {

	kapi := getKapi()
	addCourseKey := secKillConf.EtcdConf.EtcdAddCourseKey
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	respGet, err := kapi.Get(ctx, addCourseKey, nil)

	cancel()

	if err != nil {
		fmt.Println("get failed err, ", err)
		return "0"
	}

	fmt.Println("get operation : key: ", respGet.Node.Key, ", value: ", respGet.Node.Value)
	respStr := respGet.Node.Value
	addCourseNumsStruct := &AddCourseNums{}
	json.Unmarshal([]byte(respStr), &addCourseNumsStruct)
	configNums := addCourseNumsStruct.AddNums
	fmt.Println("configNums: ", configNums)
	configNUmsStr := strconv.FormatInt(configNums, 10)
	return configNUmsStr

}

func getKapi() etcd_client.KeysAPI {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints : []string{secKillConf.EtcdConf.EtcdAddr},
		HeaderTimeoutPerRequest : time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})

	if err != nil {
		fmt.Println("connect etcd failed, err", err)
		return nil
	}

	kapi := etcd_client.NewKeysAPI(cli)
	return kapi
}



type GetAllRaiseNumsFunc struct {

}

func (getAllRaiseNums GetAllRaiseNumsFunc) Execute (conn redis.Conn, redisKey string, v... interface{}) (m interface{}, err error) {

	allRaiseNumStr, err := redis.String(conn.Do("get", redisKey))
	if allRaiseNumStr == "" {
		fmt.Println("allRaiseNumStr is null, init with 0 now ...")
		conn.Do("set", redisKey, "0")
		allRaiseNumStr, err = redis.String(conn.Do("get", redisKey))
	}
	fmt.Println("allRaiseNumStr: ", allRaiseNumStr)
	if err != nil {
		return "0", err
	}
	return allRaiseNumStr, nil
}


type AddUserFunc struct {

}

func (addUser AddUserFunc) Execute (conn redis.Conn, redisKey string, v... interface{})(m interface{}, err error) {


	uinStr := getUinStr(v)
	fmt.Println("uinStr:", uinStr)

	userDetail, err := redis.String(conn.Do("hget", redisKey, uinStr))

	if userDetail == "" {
		fmt.Println("userDetail is null")
		usr := &bean.UserDetails{}
		busr, e := json.Marshal(usr)

		if e != nil {
			fmt.Println("err: ", e)
			err = e
			m = "when hset, json Marshal failed"
			return
		}

		fmt.Println("business data: ", string(busr))

		conn.Do("hset", redisKey, uinStr, string(busr))
		hgetUserDetail, hgetErr := redis.String(conn.Do("hget", redisKey, uinStr))

		if hgetErr != nil {
			err = hgetErr
			m = "when hget, error happened"
			return
		}
		fmt.Println("hgetUserDetail: ")
		fmt.Println(hgetUserDetail)

		err = nil
		m = hgetUserDetail
		return

	}

	fmt.Println("userDetail has already been inserted")
	printUserDetail := userDetail
	fmt.Println("userDetail that already inserted is : ", printUserDetail)
	return userDetail, nil
}
func getUinStr(v []interface{}) string {
	uin := v[0]
	juin, _ := json.Marshal(uin)
	uinStr := string(juin)
	uinStrLen := len(uinStr)
	uinStr = uinStr[1:uinStrLen-1]
	return uinStr
}
