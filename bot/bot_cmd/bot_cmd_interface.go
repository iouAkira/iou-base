package botcmd

// HandlerFunc 定义函数类型
type HandlerFunc func(*Context)

// HandlerFuncList 定义 HandlerFunc 函数类型切片
type HandlerFuncList []HandlerFunc

// ICommandHandler 定义所有路由 Handler 接口。
type ICommandHandler interface {
	RegCommandByChar(string, string, ...HandlerFunc) ICommandHandler
	RegCommand(Executable, string, ...HandlerFunc) ICommandHandler
	Handle(Executable, string, ...HandlerFunc) ICommandHandler
}

//Executable 对于一条指令来说需要用到以下两个方法，分别是Description和Run方法，
//Description 方法对当前指令的描述，返回值是一个字符串，
//Run 方法是执行当前指令的具体操作，
type Executable interface {
	Description(...string) string
	Run(...string) string
	Prefix() string
	SetCmd(string)
	GetCmd() string
}

// ParseExec 解析字符串参数
func ParseExec(args ...string) Executable {
	return nil
}
