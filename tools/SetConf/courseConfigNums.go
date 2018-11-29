package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
	"time"
)

type AddCourseNums struct {

	AddNums int
}

const(
	EtcdAddNumKey = "/z/addNums"
)

func SetAddNumsToEtcd(addNums int) {
	cli, err := client.New(client.Config{
		Endpoints:[]string{"http://127.0.0.1:2379"},
		HeaderTimeoutPerRequest:5*time.Second,
	})

	if err != nil {
		fmt.Println("connect failed, err:", err)
	}
	fmt.Println("connect succ")

	kapi := client.NewKeysAPI(cli)

	AddInfoConf := &AddCourseNums{addNums}



	data, err := json.Marshal(AddInfoConf)
	if err != nil {
		fmt.Println("json failed,", err)
		return
	}
	fmt.Println("data: ", string(data))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	respSet, err := kapi.Set(ctx, EtcdAddNumKey, string(data), nil)
	cancel()
	if err != nil {
		fmt.Println("put failed, ", err)
		return
	}
	fmt.Println("set operation : key: ", respSet.Node.Key, ", value: ", respSet.Node.Value)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	respGet, err := kapi.Get(ctx, EtcdAddNumKey, nil)
	cancel()

	if err != nil {
		fmt.Println("get failed err, ", err)
		return
	}

	fmt.Println("get operation : key: ", respGet.Node.Key, ", value: ", respGet.Node.Value)
}

func main() {
	SetAddNumsToEtcd(10)
}
