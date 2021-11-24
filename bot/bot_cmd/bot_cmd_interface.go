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
