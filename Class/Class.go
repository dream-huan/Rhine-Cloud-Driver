package Class

import "github.com/dgrijalva/jwt-go"

//数据库用户结构体
type User struct {
	Uid          string
	Password     string
	Email        string
	Create_time  string
	Usedstorage  int64
	Totalstorage int64
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

//文件系统结构体
type FileSystem struct {
	FileId      int64
	Create_time string
	Uid         string
	Path        string
	Isdir       bool
	Md5         string
	Parentid    int64
	Originid    int64
	FileName    string
	FileSize    int64
	Valid       bool
	OldFileName string
}

////上传文件路由结构体
//type UploadFile struct {
//	KEY string `uri:"key" binding:"required"`
//}
//
////下载文件路由结构体
//type DownloadFile struct {
//	KEY string `uri:"key" binding:"required"`
//}

//删除文件结构体
type DeletedFile struct {
	FileId int64
}

//批量删除文件结构体
type DeletedFiles struct {
	DeletedFiles []int64 `form:"deletefiles" json:"deletefiles"`
}

//创建新建文件夹结构体
type Mkdir struct {
	Parentid int64  `form:"parentid" json:"parentid"`
	Filename string `form:"filename" json:"filename"`
}

//移动文件结构体
type MoveFile struct {
	Newparentid int64 `form:"newparentid" json:"newparentid"`
	FileId      int64 `form:"fileid" json:"fileid"`
}

//批量移动文件结构体
type MoveFiles struct {
	Newparentid int64   `form:"newparentid" json:"newparentid"`
	FileIdList  []int64 `form:"fileidlist" json:"fileidlist"`
}

//批量拷贝文件结构体
type CopyFiles struct {
	Myparentid int64   `form:"myparentid" json:"myparentid"`
	FileIdList []int64 `form:"fileidlist" json:"fileidlist"`
	Shareid    string  `form:"shareid" json:"shareid"`
	Password   string  `form:"password" json:"password"`
}

//新建分享结构体
type NewShare struct {
	Password string `form:"password" json:"password"`
	FileId   int64  `form:"fileid" json:"fileid"`
	Deadtime int64  `form:"time" json:"time"`
}

type DeleteShare struct {
	Shareid string `form:"shareid" json:"shareid"`
}

//文件分享结构体
type MyShare struct {
	Shareid       string
	Filename      string
	Isdir         bool
	DeadDate      string
	ViewTimes     int64
	DownloadTimes int64
	Password      string
}
