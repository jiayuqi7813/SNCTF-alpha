package models

import (
	"myctfplat/dbconf"
	"myctfplat/tools"
	"time"
)

//用户的数据模型.
type User struct {
	Id	int
	Mark	int //用户分数
	Name	string		//显示姓名
	Email	string		//邮箱
	Stuid	string //学号
	Username	string //用户名
	Hashpass	string //哈希后的密码
	Identity	int //标识是管理员(1)还是普通用户(0)
	IfActive	bool//是否激活
	IfHidden	bool//是否隐藏
	CreatedTime	time.Time `gorm:"default:null"`
	UpdatedTime	time.Time `gorm:"default:null"`
}
//连接数据库
func link() {
	err := dbconf.Inimysql()
	if err != nil{
		panic(err)
	}
}

func LoginUser(username string,password string) (state State){
	link()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	DB.AutoMigrate(&User{})
	user := new(User)
	err := DB.First(user,"username = ?",username).Error //判断用户么是否存在
	if err != nil{
		state = NoExistUser  //用户不存在
		return
	}
	if user.Hashpass != tools.MD5(password){
		state = PassWrong //密码错误
		return
	}else {
		state = WellOp
		return
	}
}

func RegisterUser(username string,password string,email string,ifadmin int,ifactive bool)(state State){
	link()
	DB := dbconf.Dbctf
	DB.SingularTable(
		true)
	DB.AutoMigrate(&User{})
	user := new(User)
	user.Username = username
	user.Hashpass= tools.MD5(password)		//hashpasswd
	user.Email = email				//邮箱
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()
	if ifadmin == 1{
		user.IfActive = true
		user.IfHidden = true
	}else{
		user.IfActive = ifactive
		user.IfHidden = false
	}
	user.Mark=0					//分数初始化
	var err error
	//判断用户名是否重复
	err = DB.First(user,"Username = ?",username).Error

	if err == nil{
		state = UserRepeat
		return
	}
	//判断邮箱是否重复
	err = DB.First(user,"Email = ?",email).Error
	if err == nil{
		state = EmailRepeat
		return
	}
	//创建数据
	result := DB.Create(user)
	if result.Error == nil{
		state = WellOp
		return
	}else {
		state =DatabaseErr
		return
	}

}

func UserHidden()(state State,user []User){
	dbconf.DBLink()
	DB  := dbconf.Dbctf
	DB.SingularTable(true)
	DB.Table("user").Where("if_hidden = ?", 0).Order("mark desc").Find(&user)
	state=WellOp
	return
}
