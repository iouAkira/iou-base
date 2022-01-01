package botcmd

import (
	iouTools "bot/iou-tools"
	"fmt"
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// IPrefixFunc  指令前缀相关方法定义
type IPrefixFunc interface {
	GetCommandPrefixList() []string
	GetPrefix(string) string
}

// BotEngine 中注入了Bot所需要的一些，相当与一个大框架，使用时候需要New()来对Engine初始化操作
type BotEngine struct {
	engine            *BotEngine
	isRoot            bool
	HandlerFuncList   HandlerFuncList
	handlerPrefixList HandlerPrefixList
	basePath          string
	pool              sync.Pool
	bot               *tgbotapi.BotAPI
	botToken          string
	botAdminID        int64
	options           options
}

//配置项目
type options struct {
	debug   bool
	timeout int
}

// NewBotEngine 返回一个 Engine 实体,初始化操作并不包含任何路由和中间件
func NewBotEngine() *BotEngine {
	botEngine := &BotEngine{
		HandlerFuncList: nil,
		basePath:        "/",
	}
	botEngine.engine = botEngine
	return botEngine
}

// 编译检查 CommandHandler 是否实现 ICommandHandler 接口
var _ ICommandHandler = &BotEngine{}

// Handle [ICommandHandler]接口定义的方法
func (botEngine *BotEngine) Handle(httpMethod Executable, command string, handlers ...HandlerFunc) ICommandHandler {
	return botEngine.addHandle(httpMethod, command, handlers)
}

// RegCommandByChar [ICommandHandler]接口定义的方法；启动初始化注册监听指令
func (botEngine *BotEngine) RegCommandByChar(commandPrefix string, command string, handlers ...HandlerFunc) ICommandHandler {
	return botEngine.addHandle(&Command{prefix: commandPrefix, Cmd: command}, command, handlers)
}

// RegCommand 对消息命令的扩展，只要实现 Executable 接口即可对当前指令进行处理
// 例如: 如果你想实现一个 “>hit 2” , “>”即为 当前指令的prefix, 后期可以支持多字段的prefix 并且emoji亦可
func (botEngine *BotEngine) RegCommand(cmd Executable, command string, handlers ...HandlerFunc) ICommandHandler {
	return botEngine.addHandle(cmd, command, handlers)
}

//
func (botEngine *BotEngine) addHandle(cmdMethod Executable, command string, handlers HandlerFuncList) ICommandHandler {
	handlers = botEngine.combineHandlers(handlers)
	botEngine.engine.addCommand(cmdMethod, command, handlers)
	log.Printf("[CommandHandler] 注册指令: %v %v，帮助：%v", cmdMethod.Prefix(), command, cmdMethod.Description())
	return botEngine.returnObj()
}

// 注册中间件
func (botEngine *BotEngine) combineHandlers(handlers HandlerFuncList) HandlerFuncList {
	finalSize := len(botEngine.HandlerFuncList) + len(handlers)
	mergedHandlers := make(HandlerFuncList, finalSize)
	copy(mergedHandlers, botEngine.HandlerFuncList)
	copy(mergedHandlers[len(botEngine.HandlerFuncList):], handlers)
	return mergedHandlers
}

// addCommand 将所有handlerPrefix相同(指令前缀)的指令合并到一颗树组合中。
func (botEngine *BotEngine) addCommand(pHandlerPrefix Executable, command string, handlers HandlerFuncList) {
	//查询 handlerPrefix 节点是否存在  例如 "/" "@" ">"，
	//对为空对象是初始化添加节点组
	handlerPrefix := botEngine.handlerPrefixList.get(pHandlerPrefix)
	if handlerPrefix == nil {
		hp := new(HandlerPrefix)
		hp.executable = pHandlerPrefix
		hp.commands = &CommandExecList{}
		botEngine.handlerPrefixList = append(botEngine.handlerPrefixList, hp)
		handlerPrefix = botEngine.handlerPrefixList.get(pHandlerPrefix)
	}
	handlerPrefix.commands.addCommandExec(command, handlers)
}

// GetCommandPrefixList 获取所有的命令前缀
func (botEngine *BotEngine) GetCommandPrefixList() []string {
	var prefixList []string
	for _, v := range botEngine.handlerPrefixList {
		prefixList = append(prefixList, v.executable.Prefix())
	}
	return prefixList
}

// GetPrefix 获取指定命令前缀
func (botEngine *BotEngine) GetPrefix(word string) string {
	prefixList := botEngine.GetCommandPrefixList()
	for _, v := range prefixList {
		if strings.HasSuffix(word, v) {
			return v
		}
	}
	return ""
}

// Use adds middleware to the commandHandler, see example code in GitHub.
func (botEngine *BotEngine) Use(middleware ...HandlerFunc) ICommandHandler {
	botEngine.HandlerFuncList = append(botEngine.HandlerFuncList, middleware...)
	return botEngine.returnObj()
}

func (botEngine *BotEngine) returnObj() ICommandHandler {
	if botEngine.isRoot {
		return botEngine.engine
	}
	return botEngine
}

// Run botEngine 启动
func (botEngine *BotEngine) Run(token string, adminId int64, opts ...EngineConfig) {
	botEngine.botToken = token
	botEngine.botAdminID = adminId
	botEngine.options = defaultConfig
	for _, opt := range opts {
		opt(&botEngine.options)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	bot.Debug = false
	if err != nil {
		log.Panicf("start bot failed with some error %v", err)
	}
	log.Printf("Bot stared，Bot info ==> %s %s[%s]", bot.Self.FirstName, bot.Self.LastName, bot.Self.UserName)
	c := &Context{Request: bot}
	c.reset()
	c.Request = bot
	c.engine = botEngine
	botEngine.pool.Put(c)
	botEngine.bot = bot
	botEngine.handleHTTPRequest(c)
}

// handleHTTPRequest
func (botEngine *BotEngine) handleHTTPRequest(c *Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = botEngine.options.timeout
	updates := c.Request.GetUpdatesChan(u)
	for update := range updates {
		//log.Println(update.Message)
		//log.Println(update.CallbackQuery)
		if update.Message != nil && update.Message.Document != nil {
			fmt.Printf("%+v\n", "sss")
			continue
		}
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		c.mu.Lock()
		c.Update = &update
		c.mu.Unlock()
		botEngine.handleRequest(c)
	}
}

// handleRequest 对请求数据处理进不同的路由中,由解析器将用户返回的消息解析加工到 Executable 中,
func (botEngine *BotEngine) handleRequest(c *Context) {
	var msg string
	if c.Update.Message != nil {
		log.Printf("来源Message: %s", c.Update.Message.Text)
		msg = c.Update.Message.Text
	} else {
		log.Printf("来源CallbackQuery: %s", c.Update.CallbackQuery.Data)
		msg = c.Update.CallbackQuery.Data
	}
	msg = strings.Trim(msg, " ")
	//对路由集合遍历的查询开头与请求一致的指令
	var hp *HandlerPrefix
	var msgPrefix string
	for _, tree := range botEngine.handlerPrefixList {
		log.Printf("匹配指令前缀: %v", tree.executable.Prefix())
		log.Printf("匹配指令前缀: %v", tree.commands)
		if strings.HasPrefix(msg, tree.executable.Prefix()) {
			hp = tree
			msgPrefix = tree.executable.Prefix()
		}
	}
	log.Printf("%v", hp.commands)
	if hp != nil && len(*hp.commands) > 0 {
		hasCommand := false
		log.Printf("入口命令：%s", msgPrefix)
		for _, route := range *hp.commands {
			if fmt.Sprintf("%s%s", hp.executable.Prefix(), route.commandStr) == iouTools.CleanCmd(msg, 0)[0] {
				hasCommand = true
				log.Printf("存在命令 %s", route.commandStr)
				log.Printf("等待命令 %s", msg)
				log.Printf("接受命令 %s", msg)

				for _, handler := range route.handlers {
					handler(c)
				}
			}
		}
		if !hasCommand {
			log.Printf("不存在命令 %s", msg)
			unknownMsg := "/unknown"
			hp = botEngine.handlerPrefixList.get(&Command{prefix: "/"})
			for _, route := range *hp.commands {
				if fmt.Sprintf("%s%s", hp.executable.Prefix(), route.commandStr) == iouTools.CleanCmd(unknownMsg, 0)[0] {
					log.Printf("存在命令 %s", route.commandStr)
					log.Printf("等待命令 %s", unknownMsg)
					log.Printf("接受命令 %s", unknownMsg)
					for _, handler := range route.handlers {
						handler(c)
					}
				}
			}
		}
	}
}

func (botEngine *BotEngine) getRoot(msg string) (string, *HandlerPrefix) {
	var msgPrefix string
	var hp *HandlerPrefix
	//对前缀命令[root]的匹配，并将成功匹配的数据保存到hp中
	for _, tree := range botEngine.handlerPrefixList {
		if strings.HasPrefix(msg, tree.executable.Prefix()) {
			log.Printf("匹配指令前缀: %v", tree.executable.Prefix())
			hp = tree
			msgPrefix = tree.executable.Prefix()
		}
	}
	return msgPrefix, hp
}

type EngineConfig func(opt *options)

var defaultConfig = options{
	debug:   false,
	timeout: 60,
}

func DebugMode(debug bool) EngineConfig {
	return func(config *options) {
		config.debug = debug
	}
}

func TimeOut(t int) EngineConfig {
	return func(monitor *options) {
		monitor.timeout = t
	}
}
