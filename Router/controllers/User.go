package controllers

import (
	"github.com/dream-huan/Rhine-Cloud-Driver/Class"
	"github.com/dream-huan/Rhine-Cloud-Driver/common"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Jwt"
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Mysql"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Recaptcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Register(c *gin.Context) {
	var json Class.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Errorf("从JSON提取值错误:%#v", err)
		return
	}
	logger.Infow("用户操作 注册:",
		"ip", c.ClientIP(),
		"uid", json.Uid,
		"password", json.Password,
		"email", json.Email,
		"recaptchatoken", json.RecaptchaToken,
	)
	if json.RecaptchaToken == "" || !Recaptcha.VerifyToken(json.RecaptchaToken) {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   7,
			"status":  200,
		})
		return
	}
	//判断是否符合规范
	isok := true
	uidLen := len(json.Uid)
	passwordLen := len(json.Password)
	emailLen := len(json.Email)
	if uidLen >= 1 && uidLen <= 20 && passwordLen >= 1 && passwordLen <= 20 && emailLen >= 1 && emailLen <= 50 {
		for _, i := range json.Uid {
			if !(i >= '0' && i <= '9') {
				isok = false
				break
			}
		}
	} else {
		isok = false
	}
	if isok == false {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   6,
			"status":  200,
		})
		return
	}
	err := os.Mkdir(baseURL+json.Uid, os.ModePerm)
	if err != nil {
		logger.Errorf("创建新文件夹错误:%#v", err)
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   1,
			"status":  200,
		})
		return
	}
	isok = Mysql.AddUser(json.Uid, json.Password, json.Email)
	Mysql.Mkdir(json.Uid, "/", json.Uid, 0)
	if isok == true {
		c.JSON(200, gin.H{
			"message": "OK",
			"error":   -1,
			"status":  200,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   1,
			"status":  200,
		})
	}
}

func Login(c *gin.Context) {
	var json Class.Login
	_ = c.ShouldBindJSON(&json)
	uid := json.Uid
	password := json.Password
	logger.Infow("用户操作 登录:",
		"ip", c.ClientIP(),
		"uid", json.Uid,
		"password", json.Password,
	)
	if json.RecaptchaToken == "" || !Recaptcha.VerifyToken(json.RecaptchaToken) {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   7,
			"status":  200,
		})
		return
	}
	if uid != "" && password != "" && Mysql.VerifyPassword(uid, password) == true {
		token, _ := Jwt.GenerateToken(json.Uid, c.ClientIP())
		c.JSON(200, gin.H{
			"message": "OK",
			"error":   -1,
			"status":  200,
			"token":   token,
		})
		Mysql.AddLoginRecord(uid, c.ClientIP())
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   1,
			"status":  200,
		})
	}
}

func EditPassword(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	var json Class.Login
	err = c.ShouldBindJSON(&json)
	if err != nil {
		logger.Errorf("从JSON提取值错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 修改密码:",
		"ip", c.ClientIP(),
		"uid", uid,
		"token", token,
		"ok", ok,
		"newpassword", json.Password,
	)
	if ok == false {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	password := json.Password
	passwordLen := len(json.Password)
	if passwordLen < 1 || passwordLen > 20 {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   8,
			"status":  200,
		})
		return
	}
	if Mysql.EditPassword(uid, password) {
		c.JSON(200, gin.H{
			"message": "OK",
			"status":  200,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   9,
			"status":  200,
		})
	}
}

func GetInfo(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 获取个人信息:",
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
	info := Mysql.GetInfo(uid)
	c.JSON(200, gin.H{
		"uid":          info.Uid,
		"email":        info.Email,
		"create_time":  info.Create_time,
		"usedstorage":  info.Usedstorage,
		"totalstorage": info.Totalstorage,
		"message":      "OK",
		"error":        -1,
		"status":       200,
	})
}

func GetIpRecord(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		logger.Errorf("token获取错误:%#v", err)
		return
	}
	ok, uid := common.VerifyUser(c.ClientIP(), token)
	logger.Infow("用户操作 获取IP记录:",
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
	iprecord := Mysql.GetIpRecord(Jwt.TokenGetUid(token))
	c.JSON(200, iprecord)
}
