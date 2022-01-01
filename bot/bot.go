package main

import (
	botCmd "bot/bot_cmd"
	"bot/models"
	preInit "bot/pre_init"
)

func main() {
	// 读取加载程序需要使用的环境变量
	preInit.LoadBotEnv()
	// ddUtils.ExecUpCommand(upParams)
	engine := preInit.InitEngine()
	engine.Run(
		models.GlobalEnv.IouConfig.BotHandlerToken,
		models.GlobalEnv.IouConfig.BotAdminID,
		botCmd.DebugMode(false),
		botCmd.TimeOut(60),
	)
}
