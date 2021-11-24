package botcmd

import (
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// botEngine 指令前缀相关方法定义
type IPrefixFunc interface {
	GetCommandPrefixs() []string
	GetPrefix(string) string
}

// BotEngine 中注入了Bot所需要的一些，相当与一个大框架，使用时候需要New()来对Engine初始化操作
type BotEngine struct {
	HandlerFuncs     []HandlerFunc
	handlerPrfixList HandlerPrefixList
	basePath         string
	engine           *BotEngine
	pool             sync.Pool
	botToken         string
	botAdminID       int64
	bot              *tgbotapi.BotAPI
}

// New 返回一个 Engine 实体,初始化操作并不包含任何路由和中间件
func NewBotEngine() *BotEngine {
	botEngine := &BotEngine{
		HandlerFuncs: nil,
		basePath:     "/",
	}
	botEngine.engine = botEngine
	return botEngine
}

// GetCommandPrefixs 获取所有的命令前缀
func (engine *BotEngine) GetCommandPrefixs() []string {
	var prefixs []string
	for _, v := range engine.handlerPrfixList {
		prefixs = append(prefixs, v.handlerPrfix.Prefix())
	}
	return prefixs
}

// GetPrefix 获取指定命令前缀
func (engine *BotEngine) GetPrefix(word string) string {
	prefixs := engine.GetCommandPrefixs()
	for _, v := range prefixs {
		if strings.HasSuffix(word, v) {
			return v
		}
	}
	return ""
}

// Run botEngine 启动
func (botEngine *BotEngine) Run(token string, adminId int64) {
	botEngine.botToken = token
	botEngine.botAdminID = adminId
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
	// botEngine.handleHTTPRequest(c)
}
