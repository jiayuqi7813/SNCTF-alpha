package models

import (
	"fmt"
	"math"
	"myctfplat/dbconf"
)

type Challenge struct {
	ID 			int     //题目id
	Name        string	//题目名称
	Score       int		//题目分数
	Flag        string	//flag
	Description string	//描述
	Category    string	//分类
	File		string //附件
	FileLink	string //附件链接
	Tags        string	//标签
	Hints       string	//提示
	IsDocker 	bool 	//是否是动态靶机
	Visible    	bool		//是否可见
	IsBool		bool 	//临时flag与否
	SolveNum	int		//临时解题数量
}
type DockerChallenge struct {
	ID           int 			//一般路过id
	Cid			 int			//题目id
	Name         string			//题目名称
	Image_name   string			//镜像名称
	Private_port string			//镜像所需端口，web：80，pwn：9999
}
//用户调用的动态靶机数据库
type DClist struct {
	ID           int 			//一般路过id
	Cid			 int			//题目id
	Did			 int 			//动态靶机id
	Name         string			//题目名称
	Port         string			//启用端口
	Flag		 string			//flag
	Time		 int			//创建时间
}


//题目附件的数据库模型
type SubjectFile struct{
	Id	int
	ChallengeID	int        //对应题目的ID
	FileName	string//下载下来的文件名
	Md5FileName	string//Md5后存储的文件名
}
//单次flag提交记录
type Submission struct {
	ID          int  		//id
	UserID      int			//用户id
	ChallengeID int			//题目id
	Flag        string		//flag
	IP          string		//ip
	Time        int			//提交时间（时间戳
	IsBool		bool		//是否正确
}

//正确flag提交次数
type Solve struct {
	ID			int			//id
	UserID 	 	int			//用户id
	ChallengeID int			//题目id
	Time        int			//提交时间（时间戳
}
//题目分类筛选专用
type cate struct {
	Category string
}

//给予用户增加分数
func AddUserScore(uid int, cid int) error {
	var newScore int	//新分数
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	rows,_:=DB.Table("challenge").Select("score").Where("id = ?",cid).Rows()
	for rows.Next(){
		rows.Scan(&newScore)
	}

	DB.Debug().Raw("UPDATE user SET mark=mark+? WHERE id=?", newScore,uid).Scan(&User{})
	return nil

}

//用户分数削减
func UpdateUserScores(reducedScore, cid int) error {
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	//妈的动态削减，cnm
	DB.Raw("UPDATE user SET mark=mark-? WHERE EXISTS(SELECT 1 FROM solve WHERE user.id=solve.user_id AND solve.challenge_id=?);", reducedScore,cid).Scan(&User{})
	return nil

}

//分数动态计算规则
func EditChallengeScore(cid int)(reducedScore int, err error){
	var currentScore int	//原始分数
	//经典链接数据库
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	//获取题目原始分数
	rows,_:=DB.Table("challenge").Select("score").Where("id = ?",cid).Rows()
	for rows.Next(){
		rows.Scan(&currentScore)
	}
	solverCount, err := GetSolverCount(cid)	//获取解题人数
	//新分数，公式采用：https://github.com/o-o-overflow/scoring-playground
	newScore := int(100 + (1000-100)/(1.0+float64(solverCount)*0.04*math.Log(float64(solverCount))))
	reducedScore = currentScore - newScore	//分数差
	//对题目分数进行更新
	DB.Model(&Challenge{}).Where("id = ?",2).Update("score",newScore)

	return reducedScore,err

}
//获取题目解题人数
func GetSolverCount(id int) (count int, err error) {
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	rows,_:=DB.Table("solve").Select("COUNT(*)").Where("challenge_id = ?",id).Rows()
	for rows.Next(){
		err = rows.Scan(&count)
	}
	return count,err
}
//判断题目是否存在
func isChallengeExisted(id int) (exists bool) {
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	result,_ := DB.Raw("SELECT EXISTS(SELECT 1 FROM challenge WHERE id = ?);",id).Rows()
	for result.Next(){
		result.Scan(&exists)
	}
	return exists
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



//flag提交是否正确合法
func ChallengeVerify(userid int,challengeid int,flag string)(state State,submission Submission){
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	DB.AutoMigrate(&Challenge{})
	var challenge Challenge
	if userid == 0 ||challengeid == 0{		//判断参数是否合法
		state =DatabaseErr
		return
	}
	//判断题目是否存在
	if !isChallengeExisted(challengeid){
		state = NoSuchSubject
		return
	}
	//判断题目是否被解出
	if hasAlreadySolved(userid,challengeid){
		state = HasRightSubmit
		return
	}
	DB.Select("flag").Where("id = ?",challengeid).Find(&challenge)
	fmt.Println("数据库中的flag:",challenge.Flag)
	fmt.Println("提交的flag为：",flag)
	if challenge.Flag != flag{
		state = FlagWrong
		return
	}else {
		state = WellOp
		return
	}

}

//获取题目类别
func ChallengeCategory()(state State,category []cate){

	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	DB.AutoMigrate(&Challenge{})
	//DB.Debug().Select("distinct category").Find(&cate{})
	//DB.Debug().Raw("SELECT distinct category FROM challenge").Scan(&category)
	DB.Debug().Select("distinct category").Find(&Challenge{}).Scan(&category)//内容存入专有分类
	//fmt.Printf("%#v",&category)
	state = WellOp
	return
}

//获取所有未隐藏的题目
func ChallengeHidden()(state State,challenge []Challenge){
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	DB.AutoMigrate(&Challenge{})
	DB.Select([]string{"id","name","score","description","category","tags","hints"}).Where("visible = ?", 1).Find(&challenge)

	state=WellOp
	return

}