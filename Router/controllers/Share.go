package controllers

import (
	"github.com/gin-gonic/gin"
	"golandproject/Class"
	"golandproject/common"
	logger "golandproject/middleware/Log"
	"golandproject/middleware/Mysql"
	"golandproject/middleware/Redis"
	"net/http"
	"strconv"
	"time"
)

/*
1.重复分享禁止
2.dead_time和uid建立索引
3.password=_为无密码
4.password禁止出现除数字、字母外的其他字符
5.登录后才能使用转存到我的网盘功能


*/

func SplitSharePath(path, uid string, startid int64) (parentid int64) {
	var paths []string
	temp := ""
	temppath := ([]rune)(path)
	for _, v := range temppath {
		if v == '/' {
			if temp != "" {
				paths = append(paths, temp)
			}
			temp = ""
		} else {
			temp = temp + string(v)
		}
	}
	if temp != "" {
		paths = append(paths, temp)
	}
	parentid = startid
	if len(paths) > 1 {
		path = path[1:]
		parentid = Mysql.SplitPath(uid, paths[1], paths, parentid)
	}
	//for i := 1; i < len(paths); i = i + 1 {
	//	parentid = Mysql.SplitPath(uid, paths[i], parentid)
	//	if parentid == -1 {
	//		return -1
	//	}
	//}
	return parentid
}

func ShareFile(c *gin.Context) {
	var json Class.NewShare
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 分享文件:",
		"ip", c.ClientIP(),
		"uid", uid,
		"token", token,
		"ok", ok,
		"fileid", json.FileId,
		"deadtime", json.Deadtime,
		"password", json.Password,
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	key := Redis.NewShare(json.FileId, json.Deadtime, uid, json.Password)
	currenttime := time.Now()
	var deaddate string
	if json.Deadtime == 1 {
		deaddate = currenttime.AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
	} else if json.Deadtime == 2 {
		deaddate = currenttime.AddDate(0, 0, 30).Format("2006-01-02 15:04:05")
	} else {
		deaddate = "0000-00-00 00:00:00"
	}
	Mysql.ShareFile(uid, key, json.Password, deaddate, json.FileId)
	c.JSON(200, gin.H{
		"message": "OK",
		"key":     "https://pan.dreamxw.com/share/" + key,
		"status":  200,
	})
}

func GetShareFile(c *gin.Context) {
	isexist, fileid, password, uid := Redis.GetShare(c.Query("id"))
	path := c.Query("path")
	files, parentid := GetShareFileList(uid, path, fileid)
	logger.Infow("用户操作 获取分享文件:",
		"ip", c.ClientIP(),
		"id", c.Query("id"),
		"isexist", isexist,
		"path", c.Query("path"),
		"parentid", parentid,
	)
	if isexist == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   16,
		})
		return
	}
	if password == "" || c.Query("password") == password {
		c.JSON(200, gin.H{
			"message":  "OK",
			"status":   200,
			"uid":      uid,
			"isok":     true,
			"files":    files,
			"parentid": parentid,
		})
		Mysql.AddTime(c.Query("id"), 1)
	} else {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
			"uid":     uid,
		})
	}
}

func GetShareFileList(uid, path string, fileid int64) (files []Class.FileSystem, parentid int64) {
	//得到分享文件列表，构建文件系统
	//首先查看是文件夹还是文件，如果是文件夹就把内部文件拆开发送，如果不是就直接发送文件
	if path != "/" && Mysql.Isdir(fileid) {
		parentid = SplitSharePath(path, uid, fileid)
		files = Mysql.GetFileDirectory(uid, parentid)
	} else {
		parentid = fileid
		files = append(files, Mysql.GetFile(fileid))
	}
	return files, parentid
}

func CopyFiles(c *gin.Context) {
	var json Class.CopyFiles
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	isexist, fileid, password, _ := Redis.GetShare(json.Shareid)
	logger.Infow("用户操作 转存文件:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
		"id", json.Shareid,
		"isexist", isexist,
		"fileid", fileid,
		"password", password,
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	if isexist == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   16,
		})
		return
	}
	if !(password == "" || json.Password == password) {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   21,
		})
		return
	}
	isok := true
	for _, v := range json.FileIdList {
		if Mysql.VerifyFileAccess(v, fileid) == false {
			c.JSON(200, gin.H{
				"message": "NO",
				"status":  200,
				"error":   21,
			})
			return
		}
		if Mysql.FindPath(uid, v, json.Myparentid) == true {
			isok = false
			continue
		}
		if Mysql.CopyPath(v, json.Myparentid, uid) == false {
			isok = false
		}
	}
	if isok == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   17,
			"status":  200,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
		})
	}
}

func DownloadShareFile(c *gin.Context) {
	//验证文件是否所属分享id
	isexist, fileid, password, _ := Redis.GetShare(c.Query("shareid"))
	logger.Infow("用户操作 下载分享文件:",
		"ip", c.ClientIP(),
		"id", c.Query("shareid"),
		"isexist", isexist,
		"fileid", fileid,
		"password", password,
	)
	if isexist == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   16,
		})
		return
	}
	if !(password == "" || c.Query("password") == password) {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   21,
		})
		return
	}
	requestfileid, _ := strconv.ParseInt(c.Query("fileid"), 10, 64)
	if Mysql.VerifyFileAccess(requestfileid, fileid) == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   21,
		})
		return
	}
	path := Mysql.GetFilePath(requestfileid)
	key := Redis.AddDownloadKey(baseURL + path)
	c.JSON(200, gin.H{
		"message": "OK",
		"status":  200,
		"key":     key,
	})
	Mysql.AddTime(c.Query("shareid"), 2)
}

func GetMyShare(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 获取个人分享文件列表:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	c.JSON(200, gin.H{
		"message":    "OK",
		"status":     200,
		"sharelists": Mysql.GetMyShare(uid),
	})
}

func DeleteShare(c *gin.Context) {
	var json Class.DeleteShare
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 删除分享:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
		"shareid", json.Shareid,
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	isok := true
	isok = Mysql.DeleteShare(uid, json.Shareid) && isok
	isok = Redis.DeleteShare(json.Shareid) && isok
	if isok {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   20,
			"status":  200,
		})
	}
}
