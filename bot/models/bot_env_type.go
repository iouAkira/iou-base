package models

//BotEnv Bot环境变量集合
type BotEnv struct {
	RepoBaseDir           string
	DataBaseDir           string
	ContainerConfigPath   string
	BotToken              string
	BotAdminID            int64
	ReplyKeyboardFilePath string
	ReplyKeyBoard         map[string]string
	ContainerConfig       *ContainerConfig
}

var GlobalEnv = &BotEnv{}
