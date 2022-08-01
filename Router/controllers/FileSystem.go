package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golandproject/Class"
	"golandproject/common"
	logger "golandproject/middleware/Log"
	"golandproject/middleware/Mysql"
	"golandproject/middleware/Redis"
	"net/http"
	"os"
	"strconv"
)

func SplitPath(path, uid string) (parentid int64) {
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
	parentid = Mysql.SplitPath(uid, uid, paths, 0)
	return parentid
}

//todo:不能移动到自己的子目录下，后端需要限制
func MoveFiles(c *gin.Context) {
	var json Class.MoveFiles
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
	logger.Infow("用户操作 移动文件:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
		"fileidlist", json.FileIdList,
		"newparentid", json.Newparentid,
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
	for _, v := range json.FileIdList {
		if Mysql.MovePath(v, json.Newparentid, uid) == false {
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

func GetFileDirectory(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 构建个人文件系统:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
		"path", c.Query("path"),
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	path := c.Query("path")
	parentid := SplitPath(path, uid)
	//parentid := Mysql.FindPath(uid, path)
	if parentid == -1 {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   10,
			"status":  200,
		})
		return
	}
	files := Mysql.GetFileDirectory(uid, parentid)
	c.JSON(200, gin.H{
		"files":    files,
		"parentid": parentid,
	})
}

func KeyDownloadFile(c *gin.Context) {
	if len(c.Param("key")) != 64 {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   15,
		})
		return
	}
	isexist, value := Redis.GetDownloadKey(c.Param("key"))
	if isexist == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"status":  200,
			"error":   15,
		})
		return
	}
	fmt.Printf("%#v", value)
	_, err := os.Open(value)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   11,
			"status":  200,
		})
	} else {
		filename := ""
		for _, i := range value {
			if i == '/' {
				filename = ""
			} else {
				filename = filename + string(i)
			}
		}
		//c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		//c.Header("Content-Transfer-Encoding", "binary")
		c.File(value)
	}
}

func DownloadFile(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 获取文件下载密钥:",
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
	fileid, err := strconv.ParseInt(c.Query("fileid"), 10, 64)
	if err != nil {
		logger.Errorf("对fileid执行ParseInt错误:%#v", err)
		return
	}
	if Mysql.VerifyFileOwnership(uid, fileid) {
		path := Mysql.GetFilePath(fileid)
		key := Redis.AddDownloadKey(baseURL + path)
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
			"key":     key,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   14,
			"status":  200,
		})
	}
	//_, err := os.Open(baseURL + path + filename)
	////fmt.Printf("%#v", baseURL+path+filename)
	//if err != nil {
	//	fmt.Printf("%#v", err)
	//	c.JSON(200, gin.H{
	//		"message": "NO",
	//		"error":   11,
	//		"status":  200,
	//	})
	//} else {
	//	//c.Header("Content-Type", "application/octet-stream")
	//	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	//	//c.Header("Content-Transfer-Encoding", "binary")
	//	c.File(baseURL + path + filename)
	//}
}

func DeleteFile(c *gin.Context) {
	var json Class.DeletedFile
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("JSON转化错误:%#v", err)
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 删除文件:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
		"fileid", json.FileId,
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	fileid := json.FileId
	if Mysql.DeleteFile(uid, fileid) {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   12,
			"status":  200,
		})
	}
}

func DeleteFiles(c *gin.Context) {
	var json Class.DeletedFiles
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("JSON转化错误:%#v", err)
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 删除多个文件:",
		"ip", c.ClientIP(),
		"token", token,
		"uid", uid,
		"ok", ok,
		"fileidlist", json.DeletedFiles,
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
	for _, v := range json.DeletedFiles {
		if Mysql.Isdir(v) {
			if Mysql.DeleteDir(uid, v) == false {
				isok = false
			}
		} else {
			if Mysql.DeleteFile(uid, v) == false {
				isok = false
			}
		}
	}
	if isok == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   17,
			"status":  200,
		})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
		})
	}
}

func GetAvatar(c *gin.Context) {
	if c.Query("id") != "" {
		isexist, _, _, uid := Redis.GetShare(c.Query("id"))
		if isexist == false {
			c.JSON(200, gin.H{
				"message": "NO",
				"status":  200,
				"error":   16,
			})
			return
		}
		_, err := os.Stat(baseURL + uid + ".png")
		if err == os.ErrNotExist {
			logger.Errorf("获取分享者头像文件错误:%#v", err)
			c.JSON(200, gin.H{
				"message": "NO",
				"error":   22,
				"status":  200,
			})
			return
		}
		c.Header("Content-Disposition", "attachment; filename=avatar.png")
		c.File(baseURL + uid + ".png")
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	_, err = os.Stat(baseURL + uid + ".png")
	if err == os.ErrNotExist {
		logger.Errorf("获取本人头像文件错误:%#v", err)
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   22,
			"status":  200,
		})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=avatar.png")
	c.File(baseURL + uid + ".png")
}
