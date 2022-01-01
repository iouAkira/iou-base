package models

//BotEnv Bot环境变量集合
type BotEnv struct {
	IouConfigPath         string
	ReplyKeyboardFilePath string
	ReplyKeyBoard         map[string]string
	IouConfig             *IouConfig
}

var GlobalEnv = &BotEnv{}
