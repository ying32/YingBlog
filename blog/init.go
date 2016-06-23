package blog

import (
	"fmt"

	"github.com/ying32/qqwry"
)

// 全局 纯真ip数据库
var QQwryData *qqwry.QQWry

func init() {
	// 检查并写默认配置
	if !isExists(APP_CONFIG_FILENAME) {
		readConfigFile(APP_DEFAULT_CONFIG)
		// 写新的配置
		writeConfigFile()
	} else {
		// 载加配置
		if err := readConfigFile(APP_CONFIG_FILENAME); err != nil {
			fmt.Println("读app.conf配置文件失败。")
			panic(err)
		}
	}
	// 初始session
	sessionStore = newCookieStore([]byte(appSvrConf.SessionKey))

	// 初始纯真ip数据库文件
	QQwryData = qqwry.NewQQWry("qqwry.dat")
}
