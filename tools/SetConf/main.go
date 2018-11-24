package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
	"time"
)

type SecInfoConf struct {
	ProductId int
	StartTime int
	EndTime int
	Status int
	Total int
	Left int
}

const(
	EtcdKey = "/z/secskill/product"
)

func SetLogConfToEtcd() {
	cli, err := client.New(client.Config{
		Endpoints:[]string{"http://127.0.0.1:2379"},
		HeaderTimeoutPerRequest:5*time.Second,
	})



	if err != nil {
		fmt.Println("connect failed, err:", err)
	}
	fmt.Println("connect succ")

	kapi := client.NewKeysAPI(cli)

	var SecInfoConfArr []SecInfoConf
	// array append method :
	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ProductId:1082,
			StartTime:1541473511,
			EndTime:1617582365,
			Status:0,
			Total:10000,
			Left:10000,
		},
	)

	SecInfoConfArr = append(
		SecInfoConfArr,
		SecInfoConf{
			ProductId:1092,
			StartTime:1541473511,
			EndTime:1617582365,
			Status:0,
			Total:900,
			Left:900,
		},
	)

	data, err := json.Marshal(SecInfoConfArr)
	if err != nil {
		fmt.Println("json failed,", err)
		return
	}
	fmt.Println("data: ", string(data))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	respSet, err := kapi.Set(ctx, EtcdKey, string(data), nil)
	cancel()
	if err != nil {
		fmt.Println("put failed, ", err)
		return
	}
	fmt.Println("set operation : key: ", respSet.Node.Key, ", value: ", respSet.Node.Value)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	respGet, err := kapi.Get(ctx, EtcdKey, nil)
	cancel()

	if err != nil {
		fmt.Println("get failed err, ", err)
		return
	}

	fmt.Println("get operation : key: ", respGet.Node.Key, ", value: ", respGet.Node.Value)


}

func main() {
	SetLogConfToEtcd()
}
