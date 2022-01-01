package botcmd

import (
	"fmt"
	"strings"
)

//Command 指令结构体
type Command struct {
	Cmd    string   // 执行命令的名称
	Help   string   // 执行命令的介绍
	Params []string //执行命令参数
	prefix string   //命令前缀
}

//CommandHelp 给每一个命令添加一个执行前缀并动态添加到命令列表,支持多参数(pram1,pram2...)
func CommandHelp(prefix string, args ...string) string {
	if len(args) == 0 {
		return prefix
	}
	var result string
	if !strings.HasPrefix(strings.Trim(args[0], ""), "CommandPrefix") {
		result += prefix
	}
	result += "%s"
	return fmt.Sprintf(result, strings.Join(args, " "))
}

//Description 对执行执行的使用说明,支持多参数(pram1,pram2...)
func (c *Command) Description(args ...string) string {
	return fmt.Sprintf("%s %s", c.Run(args...), c.Help)
}

//Run  执行之前可以将更新参数传入函数中达到拓展功能效果,支持多参数(pram1,pram2...)
func (c *Command) Run(args ...string) string {
	args = append([]string{c.Cmd}, args...)
	return CommandHelp(c.Prefix(), args...)
}

//Prefix 定义了指令前缀
func (c *Command) Prefix() string {
	return c.prefix
}

func (c *Command) SetCmd(str string) {
	c.Cmd = str
}

func (c *Command) GetCmd() string {
	return c.Cmd
}

func (c *Command) RunToNext(f func(args ...string) string) string {
	return f()

}

var HelpCmd = Command{Cmd: "help", Help: "使用说明"}
var RdcCmd = Command{Cmd: "rdc", Help: "列出用户"}

// MyCommands 默认实现的一些命令
var MyCommands = []Command{
	HelpCmd,
	RdcCmd,
	{Cmd: "logs", Help: "下载日志"},
}

// GetCmd 用于快速添加生成命令
func GetCmd(c string) Command {
	for _, cmd := range MyCommands {
		if c == cmd.Cmd {
			return cmd
		}
	}
	return Command{Cmd: c, Help: CommandHelp(c)}

}
