package Router

import (
	"github.com/dream-huan/Rhine-Cloud-Driver/Router/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	var v1 = router.Group("/api/v1")
	v1.POST("/register", controllers.Register)
	v1.POST("/login", controllers.Login)
	v1.POST("/editpassword", controllers.EditPassword)
	v1.POST("/upload/:key", controllers.Upload)
	v1.POST("/mkdir", controllers.Mkdir)
	v1.POST("/newshare", controllers.ShareFile)
	v1.POST("/uploadavatar", controllers.UploadAvatar)
	v1.POST("/movefile", controllers.MoveFiles)
	v1.POST("/copyfile", controllers.CopyFiles)
	v1.GET("/getmyshare", controllers.GetMyShare)
	v1.GET("/getavatar", controllers.GetAvatar)
	v1.GET("/getinfo", controllers.GetInfo)
	v1.GET("/downloadsharefile", controllers.DownloadShareFile)
	v1.GET("/getiprecord", controllers.GetIpRecord)
	v1.GET("/shareinfo", controllers.GetShareFile)
	//directory := v1.Group("directory")
	//{
	//	directory.GET("*path", controllers.GetShareFileDirectory)
	//}
	v1.GET("/uploadrequest", controllers.RequestUpload)
	v1.GET("/directory", controllers.GetFileDirectory)
	v1.GET("/download", controllers.DownloadFile)
	v1.GET("/downloadfile/:key", controllers.KeyDownloadFile)
	v1.DELETE("/deletefile", controllers.DeleteFile)
	v1.DELETE("/deletefiles", controllers.DeleteFiles)
	v1.DELETE("/deleteshare", controllers.DeleteShare)
	router.Run(":8888")
	return router
}
