package main

import (
	"bot/models"
	preinit "bot/pre_init"
)

func main() {
	// 读取加载程序需要使用的环境变量
	preinit.LoadBotEnv()
	// ddutils.ExecUpCommand(upParams)
	engine := preinit.InitEngine()
	engine.Run(models.GlobalEnv.BotToken, models.GlobalEnv.BotAdminID)
}
