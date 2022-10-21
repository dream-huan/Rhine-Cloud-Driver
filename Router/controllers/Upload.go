package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/dream-huan/Rhine-Cloud-Driver/Class"
	"github.com/dream-huan/Rhine-Cloud-Driver/common"
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Mysql"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Redis"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var pwd, _ = os.Getwd()
var baseURL = pwd + "/upload/"

func Mkdir(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	var json Class.Mkdir
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("JSON转化错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 新建文件夹:",
		"ip", c.ClientIP(),
		"uid", uid,
		"token", token,
		"ok", ok,
		"dirname", json.Filename,
		"parentid", json.Parentid,
	)
	matched, _ := regexp.Match(`/[/\\+?:*<>|"]/`, []byte(json.Filename))
	if matched == true || json.Filename == "" {
		//含有/[]+?:*<>|"的字符或空不能创建
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   13,
			"status":  200,
		})
		return
	}
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	if Mysql.Mkdir(uid, json.Filename+"/", json.Filename, json.Parentid) {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   13,
			"status":  200,
		})
	}
}

func RequestUpload(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 获取上传密钥:",
		"ip", c.ClientIP(),
		"uid", uid,
		"token", token,
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
	key := Redis.AddUploadKey()
	c.JSON(200, gin.H{
		"message": "OK",
		"status":  200,
		"key":     key,
	})
}

func Upload(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		logger.Errorf("获取form错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 上传文件:",
		"ip", c.ClientIP(),
		"uid", uid,
		"token", token,
		"ok", ok,
		"key", c.Param("key"),
		"filelen", len(form.File["upload[]"]),
		"parentid", form.Value["parentid"][0],
	)
	// 1.将上传时限增至12小时，文件上传完成即刻失效防止同时多号上传
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	if !Redis.GetUploadKey(c.Param("key")) {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   15,
			"status":  200,
		})
		return
	}
	files := form.File["upload[]"]
	parentid, err := strconv.ParseInt(form.Value["parentid"][0], 10, 64)
	if err != nil {
		logger.Errorf("对parentid执行ParseInt错误:%#v", err)
		return
	}
	for _, file := range files {
		//md5值计算，同md5值的文件合并节省储存空间
		fileContent, err := file.Open()
		if err != nil {
			logger.Errorf("文件打开错误:%#v", err)
			return
		}
		byteContainer, err := ioutil.ReadAll(fileContent)
		if err != nil {
			logger.Errorf("文件转化为字节容器错误:%#v", err)
			return
		}
		md5result := md5.Sum(byteContainer)
		md5Str := fmt.Sprintf("%x", md5result)
		if fileid, filesize := Mysql.Md5Query(md5Str); fileid != 0 {
			if !Mysql.JudgeStorage(uid, filesize) {
				c.JSON(200, gin.H{
					"message": "NO",
					"status":  200,
				})
				return
			}
			if !Mysql.JudgeStorage(uid, filesize) {
				c.JSON(200, gin.H{
					"message": "NO",
					"error":   23,
					"status":  200,
				})
				return
			}
			if Mysql.CheckIsExistSame(file.Filename, parentid) {
				c.JSON(200, gin.H{
					"message": "NO",
					"error":   24,
					"status":  200,
				})
				return
			}
			Mysql.AddOldFile(uid, file.Filename, md5Str, "/", fileid, filesize, parentid)
			Mysql.ChangeStorage(uid, filesize, "+")
			continue
		}
		newFileName := uid + "_" + common.RandStringRunes(8) + "_"
		for _, v := range file.Filename {
			if v == ',' {
				newFileName = newFileName + "_"
			} else {
				newFileName = newFileName + string(v)
			}
		}
		if Mysql.CheckIsExistSame(file.Filename, parentid) {
			c.JSON(200, gin.H{
				"message": "NO",
				"error":   24,
				"status":  200,
			})
			return
		}
		c.SaveUploadedFile(file, baseURL+uid+"/"+newFileName)
		fi, err := os.Stat(baseURL + uid + "/" + newFileName)
		if err != nil {
			logger.Errorf("文件状态错误:%#v", err)
		}
		if !Mysql.JudgeStorage(uid, fi.Size()) {
			c.JSON(200, gin.H{
				"message": "NO",
				"error":   23,
				"status":  200,
			})
			return
		}
		Mysql.AddFile(uid, newFileName, md5Str, "/", file.Filename, fi.Size(), parentid)
		Mysql.ChangeStorage(uid, fi.Size(), "+")
	}
	Redis.DelUploadKey(c.Param("key"))
	c.JSON(200, gin.H{
		"message": "OK",
		"file":    len(files),
		"status":  200,
	})
}

func UploadAvatar(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 上传头像:",
		"ip", c.ClientIP(),
		"uid", uid,
		"token", token,
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
	file, err := c.FormFile("avatar")
	if err != nil {
		logger.Errorf("form获取错误:%#v", err)
	}
	c.SaveUploadedFile(file, baseURL+uid+".png")
	c.JSON(200, gin.H{
		"name":   uid + ".png",
		"status": "done",
	})
}
