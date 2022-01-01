package ioutools

import (
	"bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// MakeReplyKeyboard 构建快捷回复按钮
func MakeReplyKeyboard(config *models.BotEnv) tgbotapi.ReplyKeyboardMarkup {
	if CheckDirOrFileIsExist(config.ReplyKeyboardFilePath) {
		cookiesFile, err := ioutil.ReadFile(config.ReplyKeyboardFilePath)
		if err != nil {
			log.Printf("读取%v快捷回复配置文件出错。。%s", config.ReplyKeyboardFilePath, err)
		}

		lines := strings.Split(string(cookiesFile), "\n")
		//log.Printf("%v", lines)
		for _, line := range lines {
			lineSpt := strings.Split(line, "===")
			if len(lineSpt) > 1 {
				config.ReplyKeyBoard[lineSpt[0]] = lineSpt[1]
			}
		}
	}

	var replyKeyBoardKeys []string
	for k := range config.ReplyKeyBoard {
		replyKeyBoardKeys = append(replyKeyBoardKeys, k)
	}
	sort.Strings(replyKeyBoardKeys)

	var allRow [][]tgbotapi.KeyboardButton
	var keys []string

	for i, k := range replyKeyBoardKeys {
		//log.Printf("%v %v", i, replyKeyBoardKeys)
		keys = append(keys, k)
		if len(keys) == 2 || i == len(replyKeyBoardKeys)-1 {
			var row []tgbotapi.KeyboardButton
			for _, vi := range keys {
				row = append(row, tgbotapi.KeyboardButton{Text: vi})
			}
			allRow = append(allRow, row)
			keys = keys[0:0]
		}
	}
	replyKeyboards := tgbotapi.NewReplyKeyboard(allRow...)
	return replyKeyboards
}

// LoadReplyKeyboardMap 更新快捷回复按钮全局配置
func LoadReplyKeyboardMap(config *models.BotEnv) {
	if CheckDirOrFileIsExist(config.ReplyKeyboardFilePath) {
		cookiesFile, err := ioutil.ReadFile(config.ReplyKeyboardFilePath)
		if err != nil {
			log.Printf("读取%v快捷回复配置文件出错。。%s", config.ReplyKeyboardFilePath, err)
		}
		lines := strings.Split(string(cookiesFile), "\n")
		for _, line := range lines {
			lineSpt := strings.Split(line, "===")
			if len(lineSpt) > 1 {
				config.ReplyKeyBoard[lineSpt[0]] = lineSpt[1]
			}
		}
	}
}

func CleanCmd(cmd string, offset int) []string {
	cmdMsgSplit := strings.Split(cmd[offset:], " ")
	var arr []string
	for _, v := range cmdMsgSplit {
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	return arr
}
