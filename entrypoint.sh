#!/bin/sh

function pythonInit() {
  apk add --update python3-dev py3-pip
  if [ $PIP_REPO ]; then
    pip3 config set global.index-url $PIP_REPO
  fi
  pip3 install --upgrade pip
}
function nodeInit() {
  apk add --update nodejs npm
  if [ $NPM_REPO ]; then
    npm config set registry $NPM_REPO
  fi
}
function golangInit() {
  apk add --no-cache go
  if [ $GO_PROXY ]; then
    go env -w GO111MODULE=on
    go env -w GOPROXY=$GO_PROXY,direct
  fi
}

if [ "$1" ]; then
  up_cmd=$1
  if [ "${APK_REPO}" ]; then
    sed -i "s/dl-cdn.alpinelinux.org/${APK_REPO}/g" /etc/apk/repositories
  fi
  if [ "$APK_ADD_PKG"]; then
    apk add "$(echo $APK_ADD_PKG | tr "&" " ")"
  fi
  if [ "$INIT_ENVS" ]; then
    for env in $(echo "$INIT_ENVS" | tr "&" " "); do
      "${env}Init"
      if [ $? -ne 0 ]; then
        echo "${env}环境初始化出错❌，重启后继续尝试初始化"
        exit 1
      else
        echo "${env}环境初始化完成✅"
      fi
    done
  fi
fi

echo "------------------------------------------------读取/data/repos.json仓库相关配置，并进行执行配置------------------------------------------------"
repo_sync

echo "执行仓库入口脚本"
for repoInx in $(cat /data/repos.json | jq .repos | jq 'keys|join(" ")' | sed "s/\"//g"); do
  cd "$REPOS_DIR"
  repoName=$(cat /data/repos.json | jq ".repos | .[$repoInx] | .repo_name")
  repoBranch=$(cat /data/repos.json | jq ".repos | .[$repoInx] | .repo_branch")
  cd "$REPOS_DIR/$repoName"
  echo "切换到 $repoBranch 分支..."
  git checkout $repoBranch
  echo "执行$repoName仓库下面的程序入口脚本"
  sh iou-entry.sh
done
echo "--------------------------------------------------/data/repos.json配置结束---------------------------------------------------"

firstFile="y"
for cronFile in $(ls "$CRON_FILE_DIR" | grep ".sh" | tr "\n" " "); do
  if [ $firstFile == "y" ]; then
    echo "#$cronFile cron list" >"$CRON_FILE_DIR/merge_all_cron.sh"
    firstFile="n"
  else
    echo "#$cronFile cron list" >"$CRON_FILE_DIR/merge_all_cron.sh"
  fi
  cat -e $cronFile >>"$CRON_FILE_DIR/merge_all_cron.sh"
done

if [ "$up_cmd" ]; then
  echo "set crontab lsit"
  crontab "$CRON_FILE_DIR/merge_all_cron.sh"
  echo "keep run..."
  crond -f
else
  echo "默认定时任务执行结束。"
fi
