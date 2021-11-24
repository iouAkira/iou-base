package preinit

import (
	"log"
	"os"

	ioutools "bot/iou-tools"
	"bot/models"
)

// LoadBotEnv 加载程序需要使用的环境变量
func LoadBotEnv() {
	if ioutools.CheckDirOrFileIsExist("") {
		var containerConfig models.ContainerConfig

		models.GlobalEnv.ContainerConfig = &containerConfig
	} else {
		log.Printf("[ERROR] 程序配置目录不存在，无法启动Bot，退出。")
		os.Exit(0)
	}
}
