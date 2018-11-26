package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

type RedisConf struct {
	RedisAddr string
	RedisMaxIdle int
	RedisMaxActive int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr string
	Timeout int
	//etcdSecKey string
	EtcdSecKeyPrefix string
	EtcdSecProductKey string

}

type CodisConf struct {
	TestMode bool
}


type SecKillConf struct {
	RedisBlackConf RedisConf
	RedisProxy2LayerConf RedisConf
	EtcdConf EtcdConf
	LogPath string
	LogLevel string
	SecProductInfoMap map[int]*SecProductInfoConf
	RwSecProductLock sync.RWMutex
	CookieSecretKey string
	UserSecAccessLimit int
	IpSecAccessLimit int
	ReferWhiteList []string
	ipBlackMap map[string]bool
	idBlackMap map[int]bool
	blackRedisPool *redis.Pool
	proxy2LayerRedisPool *redis.Pool
	RwBlackLock sync.RWMutex
	WriteProxy2LayerGoroutineNume int
	ReadProxy2LayerGoroutineNums int

	secReqChan chan *SecRequst
	secReqChanSize int

	//UserConnMap map[string]chan *SecResult
	UserConnMap map[string]SecResult
	UserConnMapLock sync.Mutex
}


type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime int64
	Status int
	Total int
	left int
}

type SecResult struct {
	ProductId int
	UserId int
	Code int
	Token string
}

type SecRequst struct {
	ProductId int
	Source string
	AuthCode string
	SecTime string
	Nance string
	UserId int
	AccessTime time.Time
	UserAuthSign string
	ClientAddr string
	ClientReference string

	//CloseNotify chan bool
	//BasicInfo chan *SecResult
	BasicInfo SecResult

}

type SecResponse struct {
	ProductId int
	UserId int
	Code int
	Token string
	TokenTime int64
}


