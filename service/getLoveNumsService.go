package service

import (
	"fmt"
	"myproject/util"
)

type LoveNumsService struct {

}

func (l *LoveNumsService) GetLoveNums(uin int64, uinPublicDayDistributeLock string) (curLoveNums int64) {

	redisKey := codisUtil.GetCodisKey(util.USER_DETAILS)
	redisCourseKey := codisUtil.GetCodisKey(util.COURSE_DETAILS)
	fmt.Println("redisKey: ", redisKey)
	fmt.Println("redisCourseKey: ", redisCourseKey)





	return 0

}
