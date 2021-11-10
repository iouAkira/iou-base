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
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@                                @@"
echo "@@  开始时间 $(date +'%Y-%m-%d %H:%M:%S')  @@"
echo "@@                                @@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
if [ "$1" ]; then
    up_cmd=$1
    if [ "$APK_REPO" ]; then
        sed -i "s/dl-cdn.alpinelinux.org/$APK_REPO/g" /etc/apk/repositories | sed -e "s/^/[容器首次启动初始化] /"
    fi
    if [ "$APK_ADD_PKG" ]; then
        apk add $(echo $APK_ADD_PKG | tr "&" " ") | sed -e "s/^/[容器首次启动初始化] /"
    fi
    if [ "$INIT_ENVS" ]; then
        for env in $(echo "$INIT_ENVS" | tr "&" " "); do
            "${env}Init" | sed -e "s/^/[容器首次启动初始化] /"
            if [ $? -ne 0 ]; then
                echo "[容器首次启动初始化] [$env]环境初始化出错❌，重启后继续尝试初始化"
                exit 1
            else
                echo "[容器首次启动初始化] [$env]环境初始化完成✅"
            fi
        done
    fi
fi
# 是否指定了仓库信息配置文件的路径，未制定，自动寻找MNT_DIR根木repos.json
if [ -z "$REPOS_CONFIG" ]; then
    export REPOS_CONFIG=$MNT_DIR/repos.json
fi
echo "========================================读取[$REPOS_CONFIG]仓库相关配置，并进行执行配置========================================"
# 同步仓库
echo "-e"
repo_sync | sed -e "s/^/[仓库更新\/同步] /"

echo -e "\n>>>>>>>>>>>>>>>>>>>>执行仓库入口脚本\n"
for repoInx in $(cat $REPOS_CONFIG | jq .repos | jq 'keys|join(" ")' | sed "s/\"//g"); do
    cd "$REPOS_DIR"
    repoName=$(cat $REPOS_CONFIG | jq -r ".repos | .[$repoInx] | .repo_name")
    repoBranch=$(cat $REPOS_CONFIG | jq -r ".repos | .[$repoInx] | .repo_branch")
    repoEntrypoint=$(cat $REPOS_CONFIG | jq -r ".repos | .[$repoInx] | .repo_entrypoint")
    cd "$REPOS_DIR/$repoName"
    if [ -z "$repoBranch" ]; then
        echo "[$repoName] 仓库未指定分支，使用当前默认分支"
    else
        echo "[$repoName] 仓库切换到指定的[$repoBranch]分支..."
        git checkout $repoBranch | sed -e "s/^/[$repoName] /"
    fi
    if [ $(echo $repoEntrypoint | sed "s/null//g") ]; then
        echo "[$repoName] 仓库已配置指定的入口shell脚本文件"
        if expr "$repoEntrypoint" : 'http.*' &>/dev/null; then
            echo "[$repoName] 仓库配置指定的入口shell脚本为远程脚本，开始下载远程脚本 $repoEntrypoint"
            wget -O iou-entry.sh "$repoEntrypoint" | sed -e "s/^/[$repoName]/"
            echo "[$repoName] 仓库配置指定的入口shell脚本下载完成，开始执行..."
            sh iou-entry.sh | sed -e "s/^/[$repoName\/iou-entry.sh] /"
            echo "[$repoName] 仓库配置指定的入口shell脚本，执行结束..."
        else
            if [ ! -f "$repoEntrypoint" ]; then
                echo "[$repoName] 仓库配置指定的入口shell脚本为挂载脚本文件，但是挂载文件$repoEntrypoint不存在，跳过..."
            else
                echo "[$repoName] 仓库配置指定的入口shell脚本为挂载脚本文件，开始执行..."
                cp -rf "$repoEntrypoint" ./iou-entry.sh
                sh iou-entry.sh | sed -e "s/^/[$repoName\/iou-entry.sh] /"
                echo "[$repoName] 仓库配置指定的入口shell脚本为挂载脚本文件，开始结束..."
            fi
        fi
    else
        if [ -f "$REPOS_DIR/$repoName/iou-entry.sh" ]; then
            echo "[$repoName] 仓库为默认入口shell脚本，开始执行..."
            sh iou-entry.sh | sed -e "s/^/[$repoName\/iou-entry.sh] /"
            echo "[$repoName] 仓库为默认入口shell脚本，执行结束..."
        else
            echo "[$repoName] 仓库为默认入口shell脚本iou-entry.sh不存在，跳过..."
        fi
    fi
    echo "-e"
done

echo "========================================================$REPOS_CONFIG]配置结束================================================="

firstFile="y"
echo "05 * * * * entrypoint.sh >> $MNT_DIR/entrypoint.log 2>&1 " >"$CRON_FILE_PATH/entrypoint_cron.sh"

for cronFile in $(ls "$CRON_FILE_PATH" | grep ".sh" | grep -v "merge_all_cron.sh" | tr "\n" " "); do
    cd $CRON_FILE_PATH
    if [ $firstFile == "y" ]; then
        echo "#[$cronFile]文件任务列表" >"$CRON_FILE_PATH/merge_all_cron.sh"
        firstFile="n"
    else
        echo "#[$cronFile]文件任务列表" >>"$CRON_FILE_PATH/merge_all_cron.sh"
    fi
    cat $cronFile >>"$CRON_FILE_PATH/merge_all_cron.sh"
    echo "-e" >>"$CRON_FILE_PATH/merge_all_cron.sh"
done

echo "-e"
echo "[定时任务列表] 更新定时任务列表"
crontab "$CRON_FILE_PATH/merge_all_cron.sh"
crontab -l | sed -e "s/^/[定时任务列表] /"
echo "-e"

echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@                                @@"
echo "@@  完成时间 $(date +'%Y-%m-%d %H:%M:%S')  @@"
echo "@@                                @@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"

if [ "$up_cmd" ]; then
    echo "keep running..."
    crond -f
fi
