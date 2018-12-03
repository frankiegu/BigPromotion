package service

import (
	"myproject/bean"
	"sync"
)

type ApplyInfo struct {

}

var (
	singleInstance *ApplyInfo
	onceSync sync.Once
)

func NewApplyInfoInstance() *ApplyInfo {
	onceSync.Do(func() {
		singleInstance = new(ApplyInfo)
	})
	return singleInstance
}

func (ap *ApplyInfo) CalHistoryExchange(userInfo *bean.UserDetails) int {
	sum := 0
	if userInfo.ExchangeLoveNumRecords == nil {
		return 0
	}
	for exchangeRecord := range userInfo.ExchangeLoveNumRecords {
		sum += exchangeRecord
	}
	return sum
}

