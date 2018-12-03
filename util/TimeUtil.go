package util

import (
	"time"
)

func GetNowDay() (nowDay string) {

	loc, _ := time.LoadLocation("Asia/Shanghai")
	/*
	2006 means yyyy
	1 means month
	2 means day
	*/
	nowDay =  time.Now().In(loc).Format("200612")
	return
}