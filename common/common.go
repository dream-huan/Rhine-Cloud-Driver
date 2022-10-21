package common

import (
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Jwt"
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"math/rand"
	"os"
	"time"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成n位随机数
func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//文件读取

// 文件删除
func DeleteFile(path string) (err error) {
	err = os.Remove(path)
	if err != nil {
		logger.Errorf("文件删除错误:%#v", err)
		return err
	}
	return nil
}

// 文件移动
func MoveFile(oldpath, newpath string) {
	_ = os.Rename(oldpath, newpath)
}

func VerifyUser(ip string, token string) (bool, string) {
	uid := Jwt.TokenGetUid(token)
	if uid == "" {
		return false, ""
	}
	if !Jwt.TokenValid(token, ip) {
		return false, ""
	}
	return true, uid
}
