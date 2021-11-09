package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	rs "repo_sync/utils"
	"runtime/debug"
)

func main() {
	// ENV VER=0.1 \
    // MNT_DIR=/data \
    // CRON_FILE_PATH=/iouCron \
	mntDir="/data"
	repoBaseDir="/iouRepos"
	
	if os.Getenv("MNT_DIR") != "" {
		mntDir = os.Getenv("MNT_DIR")
	}
	if os.Getenv("REPOS_DIR") != "" {
		repoBaseDir = os.Getenv("REPOS_DIR")
	}
	repoConfigJson=fmt.Sprintf("%v/repos.json",mntDir)
	if os.Getenv("REPOS_CONFIG") != "" {
		repoConfigJson = os.Getenv("REPOS_CONFIG")
	}
	
	readRepoConfig(repoConfigJson,repoBaseDir)
}

func readRepoConfig(repoConfigJson string,repoBaseDir string) {
	if rs.Exists(repoConfigJson) {
		log.Printf("检测到仓库配置文件 %v，开始同步仓库操作。", repoConfigJson)
		var repoConfig rs.ReposConfig
		f, _ := ioutil.ReadFile(repoConfigJson)
		if err := json.Unmarshal(f, &repoConfig); err == nil {
			succCnt, failCnt := 0, 0
			for i, repo := range repoConfig.Repos {
				if repo.RepoPrivate {
					if repo.GitAccount != "" && repo.GitToken != "" {
						log.Printf("↓↓↓↓↓↓↓↓↓↓↓↓ 第%v个仓库，名字为%v，为私有库，账户、Token已配置，开始同步", i+1, repo.RepoName)
						errSr := SyncRepo(repo.RepoURL, fmt.Sprintf("%v/%v", repoBaseDir, repoDirName), repo.GitAccount, repo.GitToken)
						if errSr == nil {
							succCnt += 1
						} else {
							failCnt += 1
						}
					} else {
						log.Printf("第%v个仓库，名字为%v，为私有库，但是账户、Token未配置，同步失败", i+1, repo.RepoName)
						failCnt += 1
					}
				} else {
					log.Printf("↓↓↓↓↓↓↓↓↓↓↓↓ 第%v个仓库，名字为%v，为公开仓库，开始同步", i+1, repo.RepoName)
					errSr := SyncRepo(repo.RepoURL, repo.RepoName, repo.GitAccount, repo.GitToken)
					if errSr == nil {
						succCnt += 1
					} else {
						failCnt += 1
					}
				}

			}
			log.Printf("\n仓库同步已完成！成功%v个,失败%v个", succCnt, failCnt)
		} else {
			log.Printf("读取仓库配置文件出错，跳过同步仓库操作。请检查 %v 文件配置是否正确", repoConfigJson)
		}
	} else {
		log.Printf("仓库配置文件 %v 不存在，跳过同步仓库操作。", repoConfigJson)
	}
}
func SyncRepo(repoUrl string, repoPath string, gitAccount string, gitToken string, ) error {
	if rs.Exists(repoPath) {
		log.Printf("脚本仓库目录已存在，执行pull")
		return pullRepo(repoPath, gitAccount, gitToken)
	} else {
		log.Printf("脚本仓库目录不存在，执行clone")
		return cloneRepo(repoUrl, repoPath, gitAccount, gitToken)
	}

}

func cloneRepo(url string, directory string, gitAccount string, gitToken string) error {
	log.Printf("git clone %s to %s", url, directory)

	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: gitAccount,
			Password: gitToken,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	rs.CheckIfError(err)
	if err != nil {
		return err
	}

	ref, err := r.Head()

	rs.CheckIfError(err)

	commit, err := r.CommitObject(ref.Hash())
	rs.CheckIfError(err)
	fmt.Println(commit)
	return nil
}

func pullRepo(path string, gitAccount string, gitToken string) error {
	//对异常状态进行捕获并输出到缓冲区
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover %v\n", err)
			debug.PrintStack()
		}
	}()
	resetHard(path) //还原本地修改操作放到shell_default_scripts.sh里面
	// We instantiate a new repository targeting the given path (the .git folder)
	r, errP := git.PlainOpen(path)
	rs.CheckIfError(errP)

	// Get the working directory for the repository
	w, errW := r.Worktree()
	rs.CheckIfError(errW)

	//Pull the latest changes from the origin remote and merge into the current branch
	errPull := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      true,
		Auth: &http.BasicAuth{
			Username: gitAccount,
			Password: gitToken,
		}})
	if errPull != nil {
		if errPull.Error() == "already up-to-date" {
			log.Printf("已经是最新代码，暂无更新。")
		} else if errPull.Error() == "authentication required" {
			log.Printf("用户密码登陆失败，更新失败。")
		} else {
			log.Printf(errPull.Error())
			return errPull
		}
	} else {
		rs.CheckIfError(errPull)
		// 获取最后一次提交的信息。
		ref, errH := r.Head()
		rs.CheckIfError(errH)

		commit, errC := r.CommitObject(ref.Hash())
		rs.CheckIfError(errC)
		log.Printf("%v", commit)
	}
	return nil
}

func resetHard(path string) {
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
		log.Printf(err.Error())
	}
	if err = command.Wait(); err != nil {
		log.Printf(err.Error())
	} else {
		//fmt.Println(command.ProcessState.Pid())
		//fmt.Println(command.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
		log.Printf("还原本地修改（新增文件不受影响）防止更新冲突...")
	}
}
