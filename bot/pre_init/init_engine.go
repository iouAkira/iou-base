package preinit

import (
	botCmd "bot/bot_cmd"
	iouController "bot/iou-controller"
	"bot/models"
	"log"
)

// InitEngine 初始化bot引擎
func InitEngine() *botCmd.BotEngine {
	botEngine := botCmd.NewBotEngine()

	botEngine.Use(func(context *botCmd.Context) {
		if context.Update.Message != nil {
			log.Printf("[%s] %s", context.Update.Message.From.UserName, context.Update.Message.Text)
		}
		if context.Update.CallbackQuery != nil {
			log.Printf("[%s] %s", context.Update.CallbackQuery.From.UserName, context.Update.CallbackQuery.Data)
		}
	})
	//botEngine.RegCommandByChar("/", "cmd", iouController.SysCmdHandler(models.GlobalEnv))
	botEngine.RegCommandByChar("/", "help", iouController.HelpHandler(models.GlobalEnv))
	botEngine.RegCommandByChar("/", "start", iouController.HelpHandler(models.GlobalEnv))
	//botEngine.RegCommandByChar("/", "ddnode", iouController.DDNodeHandler(models.GlobalEnv))
	//botEngine.RegCommandByChar("/", "rdc", iouController.ReadCookieHandler(models.GlobalEnv))
	//botEngine.RegCommandByChar("/", "wskey", iouController.ReadWSKeyHandler(models.GlobalEnv))
	//engine.Cmd("ak", controller.AkController(model.Env))
	//engine.Cmd("dk", controller.DkController(model.Env))
	//engine.Cmd("clk", controller.ClearReplyKeyboardController(model.Env))
	//engine.Cmd("dl", controller.DownloadFileByUrlController(model.Env))
	//engine.Cmd("logs", controller.LogController(model.Env))
	botEngine.RegCommandByChar("/", "cancel", iouController.CancelController)
	// 注册一个未知命令响应函数
	//botEngine.RegCommandByChar("/", "unknow", iouController.UnknownController)

	return botEngine
}
