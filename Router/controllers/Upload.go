package controllers

import (
	"github.com/gin-gonic/gin"
	"golandproject/middleware/Jwt"
	"log"
)

//不对最大容量进行限制
//router.MaxMultipartMemory = 16 << 20
func Upload(c *gin.Context) {
	token, _ := c.Cookie("token")
	ip := Jwt.TokenGetIp(token)
	if token == "" || ip == "0.0.0.0" {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	if !Jwt.TokenValid(token, ip) {
		c.JSON(401, gin.H{
			"message": "NO",
			"error":   3,
			"status":  401,
		})
		return
	}
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	for _, file := range files {
		log.Println(file.Filename)
		// Upload the file to specific dst.
		c.SaveUploadedFile(file, "D:/golandproject/upload/"+file.Filename)
	}
	//c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	c.JSON(200, gin.H{
		"message": "OK",
		"file":    len(files),
		"status":  200,
	})
}
