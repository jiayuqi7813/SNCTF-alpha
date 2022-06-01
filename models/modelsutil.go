package models

type State int

const (
	WellOp  State = iota	//Everything is ok
	DatabaseErr              // 数据库内部错误
	NoSuchKey


	//用户状态
	PassWrong // 密码错误
	UserRepeat            // 已经存在用户（注册时）
	EmailRepeat            // 已经存在Email（注册时）
	NoExistUser   //用户不存在
	MarkEditWrong //修改分数失败
	NoActive //未激活
	FailActive //激活失败
	ActiveRepeat//重复激活
	NewAndOldDiff//新旧密码不一致

	//题目状态
	NoSuchSubject
	NoSuchId

	//提交flag状态
	FlagWrong //flag错误
	NoRightSubmit //没有成功提交记录
	HasRightSubmit //有成功提交记录

	//题目文件状态
	FileDeleteError//题目文件删除失败
)
