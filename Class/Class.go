package Class

import "github.com/dgrijalva/jwt-go"

//数据库用户结构体
type User struct {
	Uid         string
	Password    string
	Email       string
	Create_time string
}

//登录参数结构体
type Login struct {
	Uid            string `form:"uid" json:"uid"`
	Password       string `form:"password" json:"password"`
	Email          string `form:"email" json:"email"`
	Token          string `form:"token" json:"token"`
	RecaptchaToken string `form:"recaptchaToken" json:"recaptchaToken"`
}

//JWT个性化claims结构体
type CustomClaims struct {
	Uid string
	jwt.StandardClaims
}

//跨域结构体
type Cors struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
}

//google验证码回调结构体
type RecaptchaToken struct {
	Success      bool   `json:"success"`
	Challenge_ts string `json:"challenge_ts"`
	Hostname     string `json:"hostname"`
}

//ip登录记录结构体
type IpRecord struct {
	Id   int64
	Uid  string
	Time string
	Ip   string
	City string
}
