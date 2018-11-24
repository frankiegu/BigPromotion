package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/logs"
	_ "myproject/routers"
)





func main() {

	var errConfig error
	initConfig(errConfig)
	if errConfig != nil {
		panic(errConfig)
		return
	}
	logs.Info("initConfig succ")

	var errSec error
	initSec(errSec)
	if errSec != nil {
		panic(errSec)
		return
	}

	logs.Info("initSec succ")


	beego.Run()
}

