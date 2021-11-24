package models

// 容器配置
type ContainerConfig struct {
	BotHandlerToken string `json:"botHandlerToken"`
	BotAdminID      string `json:"botAdminID"`
	Repos           []Repo `json:"repos"`
}

// 仓库配置信息
type Repo struct {
	RepoName       string       `json:"repo_name"`
	RepoURL        string       `json:"repo_url"`
	RepoBranch     string       `json:"repo_branch"`
	RepoEntrypoint string       `json:"repo_entrypoint"`
	RegCommands    []RegCommand `json:"reg_commands"`
	RepoPrivate    bool         `json:"repo_private"`
	GitAccount     string       `json:"git_account"`
	GitToken       string       `json:"git_token"`
}

// 仓库注册bot指令需要配置的信息
type RegCommand struct {
	Prefix             string `json:"prefix"`
	Name               string `json:"name"`
	Help               string `json:"help"`
	HanderFunc         string `json:"handerFunc"`
	ControllerFilePath string `json:"controllerFilePath"`
}
