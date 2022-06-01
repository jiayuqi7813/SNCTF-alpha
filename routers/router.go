package routers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"myctfplat/api"
	"myctfplat/dbconf"
	"myctfplat/models"
	"myctfplat/tools"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//连接数据库
func link() {
	err := dbconf.Inimysql()
	if err != nil{
		panic(err)
	}
}

func Initrouter(){
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	var files []string
	filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	router.LoadHTMLFiles(files...)
	router.GET("/",api.IndexGet)
	router.Static("/static","static")
	//router.POST("/login")
	router.GET("/register",api.RegisterGET)	//注册页面
	router.POST("/register",api.RegisterPOST) //注册
	router.POST("/login",api.LoginPOST)	//登录
	router.GET("/login",api.LoginGET)
	router.NoRoute(api.Error404)
	protect := router.Group("")
	protect.Use(AuthRequired())		//cookie保护
	{
		protect.GET("/myself",api.MyselfGET)	//个人中心
		protect.GET("/logout",api.LogOutPOST)	//登出
		protect.GET("/challenges",api.ChallengeGET) //题目列表
		protect.POST("/challenges",api.ChallengeGET) //题目列表
		protect.POST("/submitflag",api.SubmitFlag)	//提交flag
		protect.GET("/notifications",api.NoticeGET) 	//公告捏
		protect.GET("/ranks",api.RanksGET)//积分排名
	}

	//v1 := router.Group("admin")
	//{
		//v1.GET("/register",api.Register)
	//}
	router.Run(":9000")
}

//中间件
func AuthRequired() gin.HandlerFunc{
	return func(c *gin.Context) {
		link()
		DB  := dbconf.Dbctf
		DB.SingularTable(true)
		DB.AutoMigrate(&models.User{})
		user := new(models.User)
		cookie,_ := c.Request.Cookie("SESSION")
		if cookie == nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"msg":"请先登录",
			})
			c.Abort()
		}
		bytesPass, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登陆"})
			c.Abort()
		}
		tpass, err := tools.AesDecrypt(bytesPass, tools.Aeskey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登陆"})
			c.Abort()
		}
		stpass := string(tpass)
		err = DB.First(user,"username = ?",stpass).Error //判断用户么是否存在
		if err != nil{
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登陆"})
			c.Abort()
		}

		c.Next()
	}
}