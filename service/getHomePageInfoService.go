package service

import (
	"fmt"
	"myproject/util"
	"github.com/garyburd/redigo/redis"

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
	res, err := lock.DoWithLock(uinPublicDayDistributeLock, DefaultTimeout, conn, AddUserFunc{})
	fmt.Println("res: ", res)
	fmt.Println("err: ", err)

}

type AddUserFunc struct {

}

func (addUser AddUserFunc) Execute (conn redis.Conn)(m interface{}, err error) {
	return 233, nil
}
