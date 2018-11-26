package service

import "myproject/util"

type HomePageService struct {

}

func (h *HomePageService) AddUser(uin int64, uinPublicDayDistributeLock string) {
	println("uinPublicDayDistributeLock: ", uinPublicDayDistributeLock)
	println("uin: ", uin)

	codisUtil := util.GetCodisUtilInstance()
	redisKey := codisUtil.GetCodisKey(util.USER_DETAILS)
	println("redisKey: ", redisKey)


}
