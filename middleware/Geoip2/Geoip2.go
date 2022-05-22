package Geoip2

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

var db *geoip2.Reader

func IpQueryCity(ip string) string {
	if ip == "127.0.0.1" {
		return "本机地址"
	}
	//exePath, err := os.Executable()
	//res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	//fmt.Printf("%#v", res)
	db, err := geoip2.Open("D:/golandproject/middleware/Geoip2/GeoLite2-City.mmdb")
	queryIp := net.ParseIP(ip)
	record, err := db.City(queryIp)
	defer db.Close()
	if err != nil {
		return "无数据"
	}
	return record.Country.Names["zh-CN"] + record.Subdivisions[0].Names["zh-CN"] + record.City.Names["zh-CN"]
}
