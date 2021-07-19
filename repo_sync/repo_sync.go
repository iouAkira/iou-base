package main

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
	"os/exec"
	"runtime/debug"
)

var (
	repoConfigJson = "/data/repos.json"
	//repoUrl 仓库地址
	//repoUrl = "https://gitee.com/lxk0301/jd_scripts"
	repoUrl = ""
	//baseScriptsPath clone到本地的绝对路径
	baseScriptsPath = "/scripts"
	githubUserName  = ""
	githubToken     = ""
)

func SyncRepo() {
	if os.Getenv("DDBOT_VER") != "" {
		if Exists(baseScriptsPath) {
			Info("脚本仓库目录已存在，执行pull")
			repoPull(baseScriptsPath)
		} else {
			Info("脚本仓库目录不存在，执行clone")
			repoClone(repoUrl, baseScriptsPath)
		}
	} else {
		Info("为了避免程序内置的用户名密码被滥用，所以会有使用场景检查，当前环境不符合使用要求，请更新镜像。")
	}
}

//Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func repoClone(url string, directory string) {

	// Clone the given repository to the given directory
	Info("git clone %s to %s", url, directory)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: githubUserName, // yes, this can be anything except an empty string
			Password: githubToken,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}

func repoPull(path string) {
	//对异常状态进行补货并输出到缓冲区
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover %v\n", err)
			debug.PrintStack()
		}
	}()
	//resetHard(path) //还原本地修改操作放到shell_default_scripts.sh里面
	// We instantiate a new repository targeting the given path (the .git folder)
	r, errP := git.PlainOpen(path)
	CheckIfError(errP)

	// Get the working directory for the repository
	w, errW := r.Worktree()
	CheckIfError(errW)

	//Pull the latest changes from the origin remote and merge into the current branch
	errPull := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      true,
		Auth: &http.BasicAuth{
			Username: githubUserName,
			Password: githubToken,
		}})
	if errPull != nil {
		if errPull.Error() == "already up-to-date" {
			Info("已经是最新代码，暂无更新。")
		} else if errPull.Error() == "authentication required" {
			Info("用户密码登陆失败，更新失败。")
		} else {
			Info(errPull.Error())
		}
	} else {
		CheckIfError(errPull)
		// 获取最后一次提交的信息。
		ref, errH := r.Head()
		CheckIfError(errH)

		commit, errC := r.CommitObject(ref.Hash())
		CheckIfError(errC)
		Info("%v", commit)
	}
}

func resetHard(path string) {
	//var execResult string
	var cmdArguments []string
	resetCmd := []string{"git", "-C", path, "reset", "--hard"}

	for i, v := range resetCmd {
		if i >= 1 {
			cmdArguments = append(cmdArguments, v)
		}
	}
	command := exec.Command(resetCmd[0], cmdArguments...)
	outInfo := bytes.Buffer{}
	command.Stdout = &outInfo
	err := command.Start()
	if err != nil {
		Info(err.Error())
	}
	if err = command.Wait(); err != nil {
		Info(err.Error())
	} else {
		//fmt.Println(command.ProcessState.Pid())
		//fmt.Println(command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
		Info("还原本地修改（新增文件不受影响）防止更新冲突.....\n%v", outInfo.String())
	}
}

// CheckIfError
func CheckIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

type ReposConfig struct {
	Repos []struct {
		RepoName     string `json:"repo_name"`
		RepoURL      string `json:"repo_url"`
		RepoBranch   string `json:"repo_branch"`
		RepoDataPath string `json:"repo_data_path"`
	} `json:"repos"`
	Ext string `json:"ext"`
}
