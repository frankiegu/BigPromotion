package util

import (
	"fmt"
	"github.com/astaxie/beego"
)


func getCodisTestMode()(b bool) {
	codisTestMode, err := beego.AppConfig.Bool("codis.testMode")
	if err != nil {
		err = fmt.Errorf("init config failed, read codis test mode:%v", err)
		return
	}
	return codisTestMode
}