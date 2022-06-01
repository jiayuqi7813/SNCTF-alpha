package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"myctfplat/models"
	"myctfplat/tools"
	"net/http"
	"time"
)

func Error404(c *gin.Context){
	c.HTML(200,"404.html",gin.H{
	})
}

func LoginGET(c *gin.Context){
	c.HTML(200,"login.html",gin.H{

	})

}

func LoginPOST(c *gin.Context){
	name := c.PostForm("name")
	passwd := c.PostForm("passwd")
	var state models.State
	state = models.LoginUser(name,passwd)		//登陆函数
	if state == models.NoExistUser{
		c.HTML(200,"login.html",gin.H{
			"Error":true,
			"Msg":"用户名不存在",
		})
		c.Redirect(302,"/login")
	}
	if state ==models.PassWrong{
		c.HTML(200,"login.html",gin.H{
			"Error":true,
			"Msg":"密码错误",
		})
		c.Redirect(302,"/login")
	}
	if state == models.WellOp{
		expiration := time.Now()
		expiration = expiration.AddDate(0, 0, 1)
		//加密
		pass := []byte(name)
		xpass, err := tools.AesEncrypt(pass,tools.Aeskey)
		if err != nil {
			fmt.Println(err)
			return
		}
		pass64 := base64.StdEncoding.EncodeToString(xpass)
		//加密做cookie写入
		cookie := http.Cookie{Name:"SESSION",Value: pass64,Expires: expiration}
		http.SetCookie(c.Writer, &cookie)
		c.HTML(200,"login.html",gin.H{
			"Error":false,
			"IsLogin":true,
			"Msg":"登录成功",

		})
		c.Redirect(http.StatusTemporaryRedirect,"/myself")

	}

}

func LogOutPOST(c *gin.Context)  {
	// 设置cookie过期
	expiration := time.Now()
	expiration = expiration.AddDate(0, 0, -1)
	cookie := http.Cookie{Name: "SESSION", Value: "", Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
	c.HTML(200,"index.html",nil)
}

func RegisterGET(c *gin.Context){
	c.HTML(200,"register.html",nil)
}

func RegisterPOST(c *gin.Context){
	name := c.PostForm("name")
	email := c.PostForm("mail")
	passwd := c.PostForm("passwd")
	vrpasswd := c.PostForm("veripasswd")
	if passwd != vrpasswd {
		c.HTML(200,"register.html",gin.H{
			"Error":true,
			"Msg":"密码不一致",
		})
		c.Redirect(302,"/register")
		return
	}
	status := models.RegisterUser(name,passwd,email,0,true)
	if status==models.WellOp{
		c.HTML(200,"register.html",gin.H{
			"Error":false,
			"success":true,
			"Msg":"注册成功",
		})
		c.Redirect(302,"/login")
	}else{
		if status == models.EmailRepeat {
			c.HTML(200,"register.html",gin.H{
				"Error":true,
				"Msg":"邮箱已被注册",
			})
			c.Redirect(302,"/register")
		} else if status == models.UserRepeat {
			c.HTML(200,"register.html",gin.H{
				"Error":true,
				"Msg":"用户名已被注册",
			})
			c.Redirect(302,"/register")
		} else {
			c.JSON(200,gin.H{
				"msg":"数据库错误！",
			})
		}
	}
	c.Redirect(302,"/register")

}