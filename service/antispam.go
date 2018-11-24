package service

import (
	"fmt"
	"sync"
)

var (
	secLimitMgr = &SecLimitMgr {
		UserLimitMap:make(map[int]*SecLimit, 10000),
		IpLimitMap:make(map[string]*SecLimit, 10000),
	}
)



type SecLimitMgr struct {
	UserLimitMap map[int]*SecLimit
	IpLimitMap map[string]*SecLimit
	lock sync.Mutex
}

func antiSpam(req *SecRequst)(err error) {


	seckLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		seckLimit = &SecLimit{}
		secLimitMgr.UserLimitMap[req.UserId] = seckLimit
	}

	count := seckLimit.Count(req.AccessTime.Unix())

	// cur user can not request num can not more than 2
	if count > secKillConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request, because request too busy")
		return
	}

	limit, ok := secLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		limit = &SecLimit{}
		secLimitMgr.IpLimitMap[req.ClientAddr] = limit
	}

	ipcount := limit.Count(req.AccessTime.Unix())
	if ipcount > secKillConf.IpSecAccessLimit {
		err = fmt.Errorf("invalid request, because request ip too busy")
	}



	return
}

type SecLimit struct {
	count int
	curTime int64
}

func (p *SecLimit) Count(nowTime int64)(curCount int) {
	if p.curTime != nowTime {
		p.count = 1
		p.curTime = nowTime
		curCount = p.count
		return
	}

	p.count ++
	curCount = p.count

	return
}

func (p *SecLimit) Check(nowTime int64) int {
	if p.curTime != nowTime {
		return 0
	}

	return p.count
}




