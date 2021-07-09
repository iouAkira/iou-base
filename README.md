# iou-tools

一个基础镜像，基于自己的使用需求总结，提供相对高度自定义的命令参数，实现自己要用的不同功能。
> 使用该镜像需要有一定的脚本阅读和动手能力，如果没有可以选择使用别人配置的仓库
- 该镜像容器启动是所有的功能依赖网络初始化，使用该镜像请确保自己网络的连通性。选择配置合适自己的源仓库地址。
- 如果要使用`crontab`定时任务的功能，请在各个仓库目录的根下面`iou-entry.sh`脚本中将任务写入到`/iouCron`的`.sh`文件内，文件命名参考`仓库名_功能目录名_cron.sh`，（镜像会默认会合并`/iouCron/`的所有`.sh`文件，合并写入系统定时任务摒弃定）
## docker-compose 环境变量

| Name |属性|说明|
| :---------: | ---- | ------------------------------------------------------------ |
| `APK_REPO` | 非必须 |部分地区用户连接`apk`官方源速度可能较慢，可通过此变量更换不同地区的源。例：`APK_REPO=mirrors.tuna.tsinghua.edu.cn`|
| `APK_ADD_PKG` | 非必须 |容器初始化启动的时候，想要安装的包，本镜像基于`alpine`，具体可安装的包可在[官方源查询](https://pkgs.alpinelinux.org/packages)，配置的安装包请确定正确，否则可能导致容器无法完成启动。例：`APK_ADD_PKG=curl&jq&wget`|
| `PIP_REPO` | 非必须 |部分地区用户连接`pip`官方源速度可能较慢，可通过此变量更换不同地区的源。例：`PIP_REPO=https://pypi.tuna.tsinghua.edu.cn/simple`|
| `NPM_REPO` | 非必须 |部分地区用户连接`npm`官方源速度可能较慢，可通过此变量更换不同地区的源。例：`NPM_REPO=https://registry.npm.taobao.org`|
| `GO_PROXY` | 非必须 |部分地区用户连接`go`资源速度可能较慢，可通过此变量更换不同地区 CND 镜像加速。例：`GO_PROXY=https://mirrors.aliyun.com/goproxy/`|
| `INIT_ENVS` | 非必须 |启动容器初始化想要安装的环境变量。例：`INIT_ENVS=node&python`|
|配置文件||(容器启动需要挂载一个`/data`文件夹，`/data`下面有`repos.json`需要使用的仓库信息配置和每个仓库对应的数据存放文件夹`repo1_data`、`repo2_data`类似等等)|
| `/data/repos.json` | 非必须 |默认为空启动一个什么功能都没有容器，可以手动进入容器配置想使用的功能。如果配置仓库地址，被使用仓库根需要包含一个`iou-entry.sh`里面包含的需要包含配置环境启动使用该仓库指令的脚本。当前`iou-entry.sh`里面可以继续嵌套调用`shell`配合完成自己想要的功能|
