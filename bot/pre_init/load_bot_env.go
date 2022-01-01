package preinit

import (
	iouTools "bot/iou-tools"
	"bot/models"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

// LoadBotEnv 加载程序需要使用的环境变量
func LoadBotEnv() {
	defaultRepoBaseDir := "/iouRepos"
	defaultDataBaseDir := "/data"
	//configFile := "/data/config.json"
	iouConfigPath := "/Users/akira/iou-config.json"

	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量
	flag.StringVar(&iouConfigPath, "config", iouConfigPath, fmt.Sprintf("默认为[%v],如果iou-config.json文件不存在于该默认路径，请使用-config指定，否则程序将不启动。", iouConfigPath))
	flag.Parse()
	log.Printf("-config 启动参数值:[%v];\n", iouConfigPath)

	if iouTools.CheckDirOrFileIsExist(iouConfigPath) {
		models.GlobalEnv.IouConfigPath = iouConfigPath

		var readConfig models.IouConfig
		f, _ := ioutil.ReadFile(iouConfigPath)
		if err := json.Unmarshal(f, &readConfig); err != nil {
			log.Fatalf("读取[%v]配置文件内容出错，退出启动", iouConfigPath)
		}
		models.GlobalEnv.IouConfig = &readConfig

		if models.GlobalEnv.IouConfig.RepoBaseDir == "" {
			log.Printf("未查找到容器内仓库文件夹存放根目录配置，使用默认仓库根目录[%v]", defaultRepoBaseDir)
			models.GlobalEnv.IouConfig.RepoBaseDir = defaultRepoBaseDir
		} else {
			log.Printf("容器内仓库文件夹存放根目录[%v]", models.GlobalEnv.IouConfig.RepoBaseDir)
		}

		if models.GlobalEnv.IouConfig.DataBaseDir == "" || iouTools.CheckDirOrFileIsExist(models.GlobalEnv.IouConfig.DataBaseDir) {
			log.Printf("未查找到容器内数据文件夹存放根目录配置，使用默认仓库根目录[%v]", defaultDataBaseDir)
			models.GlobalEnv.IouConfig.DataBaseDir = defaultDataBaseDir
		} else {
			log.Printf("容器内数据文件夹存放根目录[%v]", defaultDataBaseDir)
		}

		if models.GlobalEnv.IouConfig.BotHandlerToken == "" && models.GlobalEnv.IouConfig.BotAdminID <= 0 {
			log.Fatalf("请检查交互管理BOT配置信息是否完整。")
		}
		replyKeyBoard := map[string]string{
			"查看系统进程⛓": "/cmd ps -ef|grep -v 'grep\\| ts\\|/ts\\| sh'",
			"查看帮助说明📝": ">help",
		}
		models.GlobalEnv.ReplyKeyBoard = replyKeyBoard

	} else {
		log.Fatal("程序配置目录不存在，无法读取相关配置，退出启动。")
	}
}
