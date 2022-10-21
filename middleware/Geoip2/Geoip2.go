package Geoip2

import (
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"github.com/oschwald/geoip2-golang"
	"net"
	"os"
)

var db *geoip2.Reader
var pwd, _ = os.Getwd()
var baseURL = pwd + "/Geoip2/"

func IpQueryCity(ip string) string {
	if ip == "127.0.0.1" {
		return "本机地址"
	}
	//exePath, err := os.Executable()
	//res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	//fmt.Printf("%#v", res)
	db, err := geoip2.Open(baseURL + "GeoLite2-City.mmdb")
	queryIp := net.ParseIP(ip)
	record, err := db.City(queryIp)
	defer db.Close()
	if err != nil {
		logger.Errorf("查询ip错误:%#v", err)
		return "无数据"
	}
	return record.Country.Names["zh-CN"] + record.City.Names["zh-CN"]
}
