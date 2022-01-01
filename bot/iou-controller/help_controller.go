package iou_controller

import (
	"log"

	botCmd "bot/bot_cmd"
	iouUtils "bot/iou-tools"
	models "bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpHandler 使用说明独立控制器
func HelpHandler(env *models.BotEnv) botCmd.HandlerFunc {
	return func(ctx *botCmd.Context) {
		readme := "🙌 <b>使用说明</b> v1.0.0\n"
		//创建信息
		helpMsg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, readme)
		//tgbotapi.ChatRecordAudio
		//修改信息格式
		helpMsg.ParseMode = tgbotapi.ModeHTML
		//创建回复键盘结构体
		tkbs := iouUtils.MakeReplyKeyboard(env)
		//赋值给ReplyMarkup[快速回复]
		helpMsg.ReplyMarkup = tkbs
		//发送消息
		if _, err := ctx.Send(helpMsg); err != nil {
			log.Println(err)
		}
	}
}

// CancelController 取消按钮回复信息
func CancelController(ctx *botCmd.Context) {
	if ctx.Update.CallbackQuery != nil {
		c := ctx.Update.CallbackQuery
		edit := tgbotapi.NewEditMessageText(c.Message.Chat.ID, c.Message.MessageID, "操作已经取消")
		_, _ = ctx.Send(edit)
	}
}
