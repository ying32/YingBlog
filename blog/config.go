package blog

import (
	"encoding/json"
	"net"
	"os"
	"runtime"
)

const (
	APP_CONFIG_FILENAME = "app.conf"
	APP_DEFAULT_CONFIG  = "app.default.conf"
)

var (
	goVersion = runtime.Version()
)

// 定义可配置的服务器信息
type AppSvrConf struct {
	Port       int    `json:"port"`
	DB         string `json:"database"`
	SessionKey string `json:"sessionkey"`
	WebSite    string `json:"website"`
	UploadPath string `json:"uploadpath"`
}

// 全局配置表
var appSvrConf AppSvrConf

func readConfigFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&appSvrConf)
	if err != nil {
		f.Close()
		return nil
	}
	f.Close()
	return nil
}

func writeConfigFile() error {
	f, err := os.Create(APP_CONFIG_FILENAME)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	err = encoder.Encode(appSvrConf)
	if err != nil {
		f.Close()
		return err
	}
	f.Close()
	return nil
}

// 获取第一个
func getLocalIp() string {
	ip := "0.0.0.0"
	address, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, a := range address {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.To4().String() // 好像第一个就是活动中的主要网卡
			}
		}
	}
	return ip
}

// 判断文件或者目录是否存在
func isExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}
