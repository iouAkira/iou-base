package botcmd

// HandlerPrefix 是单一指令的集合，包含了当前指令的命令名和 当前 root 上存在的指令数组
type HandlerPrefix struct {
	executable Executable       //handlerPrefix 是指的命令前缀,如 /cmd /help中的 "/"指令接口
	commands   *CommandExecList //commands   是handlerPrefix指令下的各种命令的集合
}

// HandlerPrefixList 是指令组群，通常是需要挂载到 Engine 上来实现多个指令的启用 [[">help",">ls"],["/cmd","/echo"]]
type HandlerPrefixList []*HandlerPrefix

//get 获取当前前缀下的指令
func (hpList HandlerPrefixList) get(handlerPrefix Executable) *HandlerPrefix {
	if len(hpList) == 0 {
		return nil
	}
	for _, hp := range hpList {
		if hp.executable.Prefix() == handlerPrefix.Prefix() {
			return hp
		}
	}
	return nil
}

// CommandExec 为指令节点，每个指令节点都对应
type CommandExec struct {
	commandStr string          // 指令
	handlers   HandlerFuncList // 回调处理函数
}

type CommandExecList []*CommandExec

// addCommandNode 添加指令并封装加入函数
func (execList *CommandExecList) addCommandExec(path string, handlers HandlerFuncList) {
	cmdPath := path
	var hasPath bool
	for _, s := range *execList {
		hasPath = false
		if s.commandStr == cmdPath {
			hasPath = true
			break
		}
	}
	if !hasPath {
		*execList = append(*execList, &CommandExec{commandStr: cmdPath, handlers: handlers})
	}
}

func (receiver *HandlerPrefix) to(cmdStr string) *CommandExec {
	if len(*receiver.commands) == 0 {
		return nil
	}
	for _, tree := range *receiver.commands {
		if tree.commandStr == cmdStr {
			return tree
		}
	}
	return nil
}
