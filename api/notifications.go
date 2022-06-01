package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"myctfplat/dbconf"
	"myctfplat/models"
	"myctfplat/tools"
)

func NoticeGET(c *gin.Context){
	//session验证解码获取用户名
	cookie,_ := c.Request.Cookie("SESSION")
	bytesPass, _ := base64.StdEncoding.DecodeString(cookie.Value)
	tpass, _ := tools.AesDecrypt(bytesPass, tools.Aeskey)
	stpass := string(tpass)
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	DB.AutoMigrate(models.Notifications{})
	var notices []models.Notifications
	c.HTML(200,"navbar.html",gin.H{		//导航栏
		"IsLogin":true,
		"username": stpass,
	})
	DB.Find(&notices)//数据检索
	c.HTML(200,"notice.html",gin.H{		//渲染
		"notice":notices,
	})
}