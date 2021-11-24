package preinit

import (
	botcmd "bot/bot_cmd"
)

// 初始化bot引擎
func InitEngine() *botcmd.BotEngine {
	botEngine := botcmd.NewBotEngine()
	return botEngine
}
