package service

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"myproject/bean"
	"myproject/util"
)

type HomePageService struct {

}

func (h *HomePageService) AddUser(uin int64, uinPublicDayDistributeLock string) {
	fmt.Println("uinPublicDayDistributeLock: ", uinPublicDayDistributeLock)
	fmt.Println("uin: ", uin)

	codisUtil := util.GetCodisUtilInstance()
	redisKey := codisUtil.GetCodisKey(util.USER_DETAILS)
	fmt.Println("redisKey: ", redisKey)

	lock := &util.Lock{}
	conn, err := redis.Dial("tcp", "localhost:6379")
	defer conn.Close()

	DefaultTimeout := 1
	res, err := lock.DoWithLock(uinPublicDayDistributeLock, DefaultTimeout, conn, AddUserFunc{}, redisKey, uin)
	fmt.Println("AddUser lock.DoWithLock res: ", res)
	fmt.Println("AddUser lock.DoWithLock err: ", err)

}

type AddUserFunc struct {

}

func (addUser AddUserFunc) Execute (conn redis.Conn, redisKey string, uin int64)(m interface{}, err error) {

	uinStr := string(uin)
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
