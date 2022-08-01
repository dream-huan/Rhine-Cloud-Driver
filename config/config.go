package config

import (
	"bufio"
	logger "golandproject/middleware/Log"
	"io"
	"os"
	"strconv"
	"strings"
)

var pwd, _ = os.Getwd()
var privatekey = ""
var originstorage int64
var jwtkey = ""

func GetOriginStorage() int64 {
	return originstorage
}

func GetPrivateKey() string {
	return privatekey
}

func GetJwtKey() string {
	return jwtkey
}

func ReadIni() {
	file, err := os.Open(pwd + "/setting.ini")
	if err != nil {
		logger.Errorw("读取ini配置文件出错！", "err", err)
	}
	r := bufio.NewReader(file)
	i := 1
	for {
		lineBytes, err := r.ReadBytes('\n')
		line := strings.TrimSpace(string(lineBytes))
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		if i == 4 || i == 9 || i == 14 {
			isok := 0
			storage := ""
			for _, v := range line {
				if v == ']' {
					break
				}
				if isok == 1 {
					if i == 4 {
						privatekey += string(v)
					}
					if i == 9 {
						jwtkey += string(v)
					}
					if i == 14 {
						storage += string(v)
					}
				}
				if v == '[' {
					isok = 1
				}
			}
			if i == 14 {
				originstorage, _ = strconv.ParseInt(storage, 10, 64)
			}
		}
		i = i + 1
	}
	//fmt.Printf("%s", privatekey)
	//fmt.Printf("%s", jwtkey)
	//fmt.Printf("%d", originstorage)
}
