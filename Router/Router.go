package Router

import (
	"github.com/gin-gonic/gin"
	"golandproject/Router/controllers"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	var v1 = router.Group("/api/v1")
	v1.POST("/register", controllers.Register)
	v1.POST("/login", controllers.Login)
	v1.POST("/editpassword", controllers.EditPassword)
	v1.POST("/upload", controllers.Upload)
	v1.GET("/getinfo", controllers.GetInfo)
	v1.GET("/getiprecord", controllers.GetIpRecord)
	router.Run(":8888")
	return router
}
