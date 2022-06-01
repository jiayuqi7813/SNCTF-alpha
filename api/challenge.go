package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"myctfplat/dbconf"
	"myctfplat/models"
	"myctfplat/tools"
	"net/http"
	"strconv"
	"time"
)



var Status int
var AlertMsg string

//flag提交接口
func SubmitFlag(c *gin.Context){
	//session验证解码获取用户名
	cookie,_ := c.Request.Cookie("SESSION")
	bytesPass, _ := base64.StdEncoding.DecodeString(cookie.Value)
	tpass, _ := tools.AesDecrypt(bytesPass, tools.Aeskey)
	stpass := string(tpass)
	//数据库操作
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	var uid int
	//根据用户查询uid
	result,_:= DB.Debug().Table("user").Select("id").Where("username = ?",stpass).Rows()
	for result.Next(){
		err := result.Scan(&uid)
		if err != nil {
			log.Fatal(err)
		}
	}
	//uid := c.PostForm("uid") //获取用户id
	flag := c.PostForm("flag")	//获取用户提交flag值
	cid := c.PostForm("cid")	//获取题目id
	intcid,_ :=strconv.Atoi(cid)		//强制转换，他妈的
	//检测flag是否合法
	state,submiss := models.ChallengeVerify(uid,intcid,flag)
	if state == models.DatabaseErr{
		c.JSON(200,gin.H{
			"error":"非法数据！",
		})
		return
	}
	//题目不存在
	if state == models.NoSuchSubject{
		c.JSON(200,gin.H{
			"error":"题目不存在！",
		})
		return
	}
	if state == models.HasRightSubmit{
		//题目已提交
		Status = 1
		AlertMsg = "此题目已提交，请勿重复提交"
		c.Redirect(http.StatusTemporaryRedirect,"/challenges")//307
	}

	//提交记录统一赋值
	submiss.UserID=uid
	submiss.ChallengeID=intcid
	submiss.IP=c.ClientIP()
	submiss.Flag=flag
	submiss.Time=int(time.Now().Unix())
	//flag正确记录赋值
	solve := models.Solve{
		UserID: uid,
		ChallengeID: intcid,
		Time: int(time.Now().Unix()),
	}
	if state == models.FlagWrong{
		//flag错误
		submiss.IsBool=false
		DB.Table("submission").Create(&submiss)//添加flag提交数据
		Status = 3
		AlertMsg = "flag错误！"
		c.Redirect(http.StatusTemporaryRedirect,"/challenges")//307
	}

	if state == models.WellOp{
		//flag正确
		submiss.IsBool=true
		DB.Table("submission").Create(&submiss)//添加flag提交数据
		DB.Table("solve").Create(&solve)		//正确记录
		//给用户增加分数
		models.AddUserScore(uid,intcid)
		//获取分数差
		reducedScore,_ := models.EditChallengeScore(intcid)
		//削减其他用户分数
		models.UpdateUserScores(reducedScore,intcid)
		Status = 2
		AlertMsg = "flag提交成功！"
		c.Redirect(http.StatusTemporaryRedirect,"/challenges")//307
	}




}


//练习场页面接口
func ChallengeGET(c *gin.Context){
	//session验证解码获取用户名
	cookie,_ := c.Request.Cookie("SESSION")
	bytesPass, _ := base64.StdEncoding.DecodeString(cookie.Value)
	tpass, _ := tools.AesDecrypt(bytesPass, tools.Aeskey)
	stpass := string(tpass)
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	var uid int
	//根据用户查询uid
	result,_:= DB.Debug().Table("user").Select("id").Where("username = ?",stpass).Rows()
	for result.Next(){
		err := result.Scan(&uid)
		if err != nil {
			log.Fatal(err)
		}
	}
	state,challenges := models.ChallengeHidden()	//把题目隐藏
	_,category := models.ChallengeCategory()		//获取题目类型
	for index,challenge := range challenges{
		//fmt.Println(index,challenge.ID)
		challenges[index].SolveNum ,_ =models.GetSolverCount(challenge.ID)	//获取题目解题数
		if hasAlreadySolved(uid,challenge.ID) == true{		//获取题目是否已解出
			challenges[index].IsBool = true
		}else {
			challenges[index].IsBool=false
		}
	}

	if state ==models.WellOp{
		c.HTML(200,"navbar.html",gin.H{		//导航栏
			"IsLogin":true,
			"username": stpass,
		})

		c.HTML(200,"challenge.html",gin.H{		//主页面
			"category":category,
			"challenges":challenges,
			"status" : Status,
			"msg":AlertMsg,
		})

	}


}
//判断题目是否被解出
func hasAlreadySolved(uid int, cid int) (exists bool) {
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	result,_ :=DB.Raw("SELECT EXISTS(SELECT 1 FROM solve WHERE user_id = ? AND challenge_id = ?);",uid,cid).Rows()
	for result.Next(){
		result.Scan(&exists)
	}
	return exists
}

