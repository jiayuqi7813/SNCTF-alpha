package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"myctfplat/tools"
)


//个人信息
func MyselfGET(c *gin.Context){
	cookie,_ := c.Request.Cookie("SESSION")
	bytesPass, _ := base64.StdEncoding.DecodeString(cookie.Value)
	tpass, _ := tools.AesDecrypt(bytesPass, tools.Aeskey)
	stpass := string(tpass)
	//str := fmt.Sprintf("hello,%s",stpass)
	//不知道为啥IsLogin直接在下面渲染传不过去，直接这里再次渲染一次
	c.HTML(200,"navbar.html",gin.H{
		"IsLogin":true,
		"username": stpass,
	})
	//渲染页面，传数据
	c.HTML(200,"usersetting.html",gin.H{
		"IsLogin":true,
		"username": stpass,
	})
}
