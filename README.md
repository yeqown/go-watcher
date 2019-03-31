# go-watcher

Golang编写的热重载工具，自定义命令，支持监视文件及路径配置，环境变量配置。这是一个重复的轮子～

### 安装使用

```go
go install github.com/yeqown/go-watcher/cmd/go-watcher
```

### 命令行
```
NAME: P
   go-watcher - A new cli application

USAGE:
   go-watcher [global options] command [command options] [arguments...]

VERSION:
   1.1.0

AUTHOR:
   yeqown@gmail.com

COMMANDS:
     init     complete a task on the list
     run      execute a command, and watch the files, if any change to these files, the command will reload
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  load configuration from FILE, --default=./config.yml (default: "./config.yml")
   --help, -h              show help
   --version, -v           print the version
```

### 配置文件

`go-watcher init` // 初始化配置文件

```yml
# go-watcher.yml

# 需要监听的除当前目录以外的目录
extern_paths:
  - $PATH/project/demo

# 需要排除的文件表达式
exclude_regexps:
  - ".yml$"
  - ".txt$"

# 热加载命令的环境变量
envs:
  - GOOS=linux
  - GOPATH=/your/gopath
  - GOROOT=/usr/local/go

# 需要排除的文件夹，支持绝对路径和相对路径, 默认添加了当前系统用户的所有环境变量
exclude_paths:
  - ./vendor
  - ./testdata
  # - abspath is also ok

```
