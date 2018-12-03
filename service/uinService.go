package service

import "encoding/json"

func getUinStr(v []interface{}) string {
	uin := v[0]
	juin, _ := json.Marshal(uin)
	uinStr := string(juin)
	uinStrLen := len(uinStr)
	uinStr = uinStr[1:uinStrLen-1]
	return uinStr
}
