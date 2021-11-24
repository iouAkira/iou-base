package botcmd

import (
	"fmt"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// auth @clamp
type Context struct {
	Request          *tgbotapi.BotAPI
	Update           *tgbotapi.Update
	mu               sync.RWMutex
	Keys             map[string]interface{}
	engine           *BotEngine
	HandlerPrefixStr string
}

// auth @clamp
func (c *Context) reset() {
	*c = Context{}
}

// Get
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// Send 发送消息给Bot
func (c *Context) Send(chat tgbotapi.Chattable) (tgbotapi.Message, error) {
	return c.Request.Send(chat)
}

func (c *Context) Req(chat tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return c.Request.Request(chat)
}

//Vars 获取当前路由下的参数信息，当前用空格分割并返回给一个切片数组
func (c *Context) Vars() []string {
	var cc string
	if c.Update.Message != nil {
		cc = c.Update.Message.Text
	} else {
		cc = c.Update.CallbackQuery.Data
	}
	cmd, err := ParseCmd(cc, c.engine)
	if err != nil {
		return nil
	}
	return cmd.Params
}

//VarsCallback 获取当前路由下的参数信息，当前用空格分割并返回给一个切片数组
func (c *Context) Varswiyh() []string {
	var cc string
	if c.Update.Message != nil {
		cc = c.Update.Message.Text
	} else {
		cc = c.Update.CallbackQuery.Data
	}
	cmd, err := ParseCmd(cc, c.engine)
	if err != nil {
		return nil
	}
	return cmd.Params
}

//Message 获取当前路由下的消息信息
func (c *Context) Message(ctx *Context) *tgbotapi.Message {
	var cc *tgbotapi.Message
	if c.Update.Message != nil {
		cc = ctx.Update.Message
	} else {
		cc = ctx.Update.CallbackQuery.Message
	}
	return cc
}

// ParseCmd 解析命令参数,对已经首字节为"/"进行裁剪，保留非空参数,并且把剩余按空格切分口存入 Command.Params 中
func ParseCmd(cmd string, engine IPrefixFunc) (Command, error) {
	cmd = strings.Trim(cmd, " ")
	if cmd == "" {
		return Command{}, fmt.Errorf("命令不能为空")
	}
	commandPrefix := engine.GetPrefix(cmd)
	if commandPrefix == "" {
		return Command{}, fmt.Errorf("前缀无法识别")
	}
	//log.Printf("ParseCmd commandPrefix: %v", commandPrefix)
	if !strings.HasPrefix(cmd, commandPrefix) {
		return Command{}, fmt.Errorf("非法命令")
	}
	cmd = strings.Trim(cmd, commandPrefix)
	cmdMsgSplit := strings.Split(cmd, " ")
	var arr []string
	for _, v := range cmdMsgSplit {
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	cmdST := Command{prefix: commandPrefix, Cmd: cmdMsgSplit[0]}
	//log.Println(cmdST)
	//log.Printf("cmdST:%+v", cmdST)
	cmdST.Params = cmdMsgSplit[1:]
	return cmdST, nil
}
