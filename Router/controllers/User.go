package controllers

import (
	"github.com/gin-gonic/gin"
	"golandproject/Class"
	"golandproject/middleware/Jwt"
	"golandproject/middleware/Mysql"
	"golandproject/middleware/Recaptcha"
	"net/http"
)

func Register(c *gin.Context) {
	var json Class.Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	isok = Mysql.AddUser(json.Uid, json.Password, json.Email)
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
	//if err := c.ShouldBindJSON(&json); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	uid := json.Uid
	password := json.Password
	token, _ := c.Cookie("token")
	//获取IP地址
	ip := c.ClientIP()
	if token == "" && (json.RecaptchaToken == "" || !Recaptcha.VerifyToken(json.RecaptchaToken)) {
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   7,
			"status":  200,
		})
		return
	}
	//密码验证自动延长token时间
	if token == "" && uid != "" && password != "" && Mysql.VerifyPassword(uid, password) == true {
		token, _ := Jwt.GenerateToken(json.Uid, ip)
		c.JSON(200, gin.H{
			"message": "OK",
			"error":   -1,
			"status":  200,
			"token":   token,
		})
		Mysql.AddLoginRecord(uid, ip)
	} else if token != "" { //使用token登录则不延长
		ip = Jwt.TokenGetIp(token)
		if Jwt.TokenValid(token, ip) {
			c.JSON(200, gin.H{
				"message": "OK",
				"error":   -1,
				"status":  200,
			})
			Mysql.AddLoginRecord(Jwt.TokenGetUid(token), ip)
		} else {
			c.JSON(401, gin.H{
				"message": "NO",
				"error":   3,
				"status":  401,
			})
		}
	} else { //密码也不行 token也不行 登录失败
		c.JSON(200, gin.H{
			"message": "NO",
			"error":   1,
			"status":  200,
		})
	}
}

func EditPassword(c *gin.Context) {
	var json Class.Login
	_ = c.ShouldBindJSON(&json)
	//if err := c.ShouldBindJSON(&json); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	password := json.Password
	token, _ := c.Cookie("token")
	passwordLen := len(json.Password)
	//密码验证自动延长token时间
	uid := Jwt.TokenGetUid(token)
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
	token, _ := c.Cookie("token")
	//获取IP地址
	ip := c.ClientIP()
	ip = Jwt.TokenGetIp(token)
	if Jwt.TokenValid(token, ip) {
		info := Mysql.GetInfo(Jwt.TokenGetUid(token))
		c.JSON(200, gin.H{
			"uid":         info.Uid,
			"email":       info.Email,
			"create_time": info.Create_time,
			"message":     "OK",
			"error":       -1,
			"status":      200,
		})
	} else {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
	}
}

func GetIpRecord(c *gin.Context) {
	token, _ := c.Cookie("token")
	//获取IP地址
	ip := c.ClientIP()
	ip = Jwt.TokenGetIp(token)
	//密码验证自动延长token时间
	if Jwt.TokenValid(token, ip) {
		iprecord := Mysql.GetIpRecord(Jwt.TokenGetUid(token))
		c.JSON(200, iprecord)
	} else {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
	}
}
