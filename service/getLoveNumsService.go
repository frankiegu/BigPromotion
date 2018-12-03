package service

import (
	"encoding/json"
	"fmt"
	"myproject/bean"
	"myproject/util"
	"github.com/garyburd/redigo/redis"

)

const(
	NINE_MINS int64 = 9 * 60
)

type LoveNumsService struct {

}

func (l *LoveNumsService) GetLoveNums(uin int64, uinPublicDayDistributeLock string) (curLoveNums int64) {

	redisKey := codisUtil.GetCodisKey(util.USER_DETAILS)
	redisCourseKey := codisUtil.GetCodisKey(util.COURSE_DETAILS)
	fmt.Println("redisKey: ", redisKey)
	fmt.Println("redisCourseKey: ", redisCourseKey)

	lock := &util.Lock{}
	conn, err := getRedisConn()
	defer conn.Close()

	res, err := lock.DoWithLock(uinPublicDayDistributeLock, DefaultTimeout, conn, getLoveNumsFunc{}, redisKey, uin, redisCourseKey)

	fmt.Println(res)

	if err != nil {
		return 0
	}
	return 0

}

type getLoveNumsFunc struct {

}

func (getLoveNums getLoveNumsFunc) Execute (conn redis.Conn, redisKey string, v... interface{})(m interface{}, err error) {

	res := 0
	uinStr := getUinStr(v)
	redisCourseKey := getRedisCourseKey(v)
	fmt.Println("uinStr: ", uinStr)
	fmt.Println("redisCoursekey: ", redisCourseKey)

	userDetailCheck, err := redis.String(conn.Do("hget", redisKey, uinStr))
	if userDetailCheck == "" {
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

		m = res
		err = fmt.Errorf("it is a new user, so it's love num is 0")
		return
	}

	// there is this person

	nowDay := util.GetNowDay()
	fmt.Println("nowDay: ", nowDay)
	userDetail, _ := redis.String(conn.Do("hget", redisKey, uinStr))
	fmt.Println("hgetUserDetail: ", userDetail)

	//already test, user corrent
	user := &bean.UserDetails{}
	json.Unmarshal([]byte(userDetail), &user)

	//testPrintUserObject, _ := json.Marshal(user)
	//fmt.Println("user: ", string(testPrintUserObject))

	learnTIme2LoveNums := getLearnTime2LoveNums(user, uinStr)
	fmt.Println("learnTIme2LoveNums: ", learnTIme2LoveNums)
	//userSignUp2LoveNums := getUserSignUp2LoveNums(user, nowDay, uinStr, )


	return "233", nil

}
func getRedisCourseKey(v []interface{}) string {
	courseKey := v[1]
	jCourseKey, _ := json.Marshal(courseKey)
	courseKeyStr := string(jCourseKey)
	courseKeyStrLen := len(courseKeyStr)
	courseKeyStr = courseKeyStr[1:courseKeyStrLen-1]
	return courseKeyStr

}
func getLearnTime2LoveNums(userInfo *bean.UserDetails, uinStr string) int {
	curDayLearTime := getCurDayLearnTIme(uinStr)
	previousLearnTime := getPreviousLearnTime(userInfo)

	fmt.Println("curDayLearTime: ", curDayLearTime)
	fmt.Println("previousLearnTime: ", previousLearnTime)
	allLearnTime := previousLearnTime + curDayLearTime
	applyInfo := NewApplyInfoInstance()
	previousExchangeLoveNums := applyInfo.CalHistoryExchange(userInfo)

	learnTime2LoveNums := int((allLearnTime - int64(previousExchangeLoveNums) * NINE_MINS) / NINE_MINS)

	return learnTime2LoveNums
}

func getPreviousLearnTime(userInfo *bean.UserDetails) int64 {
	var previousLearnTime int64 = 0
	if userInfo.PreviousLearnTime != 0 {
		previousLearnTime = userInfo.PreviousLearnTime
	}
	return previousLearnTime

}
func getCurDayLearnTIme(uinStr string) int64 {
	//rpc : user center, but in order to make a demo, here we set as 60
	return 60
}