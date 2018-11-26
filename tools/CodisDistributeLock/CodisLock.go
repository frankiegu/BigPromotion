package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

var (
	delScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1] then 
		return redis.call("del", KEYS[1]) 
	else 
		return 0
	end`)
)

const(
	LOCK_SUCCESS = "OK"
	SET_IF_NOT_EXIST = "NX"
	SET_WITH_EX = "EX"
	RELEASE_SUCCESS = int64(1)

	GET_LOCK_FAIL = 1
	UNLOCK_FAIL = 2
)

type Lock struct {
	resource string
	token string
	conn redis.Conn
	timeout int
}

func (lock *Lock) tryLock()(ok bool, err error) {
	_, err = redis.String(lock.conn.Do("SET", lock.key(), lock.token, SET_WITH_EX, int(lock.timeout), SET_IF_NOT_EXIST))

	if err == redis.ErrNil {
		// The lock was not successful, it already exists
		ok = false
		err = nil
		return
	}

	if err != nil {
		ok  =  false
		err = err
		return
	}

	ok  = true
	err = nil
	return

}

func (lock *Lock) Unlock() (err error) {

	delScript.Do(lock.conn, lock.key(), lock.token)
	//_, err = lock.conn.Do("del", lock.key())
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

func (lock *Lock) AddTimeout(ex_time int64) (ok bool, err error) {
	ttl_time, err := redis.Int64(lock.conn.Do("TTL", lock.key()))

	if err != nil {
		log.Fatal("redis get failed: ", err)
	}

	if ttl_time > 0 {
		_, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, SET_WITH_EX, int(ttl_time + ex_time)))
		if err == redis.ErrNil {
			return false, nil

		}
		if err != nil {

			return false, err
		}
	}

	return true, nil
}

func TryLock(conn redis.Conn, resource string, token string, DefaultTimeout int) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(conn, resource, token, DefaultTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, token string, timeout int) (lock *Lock, ok bool, err error) {
	lock = &Lock{resource, token, conn, timeout}

	ok, err = lock.tryLock()
	if !ok || err != nil {
		lock = nil
	}
	return
}

func main() {
	fmt.Println("start")
	DefaultTimeout := 1

	conn, err := redis.Dial("tcp", "localhost:6379")

	//only one
	requestId := uuid.Must(uuid.NewV4())
	fmt.Println("uuid: ", requestId)
	lock, ok, err := TryLock(conn, "z", fmt.Sprintf("%d", requestId), int(DefaultTimeout))

	if err != nil {
		log.Fatal("Error while attempting lock")
	}

	if !ok {
		log.Fatal("challenge lock failed ... ")
	}
	lock.AddTimeout(1)

	time.Sleep(time.Duration(DefaultTimeout) * time.Second)
	fmt.Println("end")
	defer lock.Unlock()
}


