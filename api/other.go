package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"myctfplat/models"
	"myctfplat/tools"
)
//排名列表
func RanksGET(c *gin.Context){
	cookie,_ := c.Request.Cookie("SESSION")		//无语子的渲染
	bytesPass, _ := base64.StdEncoding.DecodeString(cookie.Value)
	tpass, _ := tools.AesDecrypt(bytesPass, tools.Aeskey)
	stpass := string(tpass)
	c.HTML(200,"navbar.html",gin.H{
		"IsLogin":true,
		"username":stpass,
	})
	state,users := models.UserHidden()	//隐藏需要隐藏的用户
	if state==models.WellOp{
		c.HTML(200,"ranks.html",gin.H{
			"users":users,
		})
	}

}


//首页渲染
func IndexGet(c *gin.Context){
	cookie,_ := c.Request.Cookie("SESSION")
	if cookie == nil{
		c.HTML(200,"navbar.html",gin.H{
			"IsLogin":false,
		})
		c.HTML(200,"index.html",nil)
	}else {
		//如果有cookie则解密用户名，渲染前端页面
		bytesPass, _ := base64.StdEncoding.DecodeString(cookie.Value)
		tpass, _ := tools.AesDecrypt(bytesPass, tools.Aeskey)
		stpass := string(tpass)
		c.HTML(200,"navbar.html",gin.H{
			"IsLogin":true,
			"username":stpass,
		})
		c.HTML(200,"index.html",nil)
	}

}

