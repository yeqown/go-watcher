# go-watcher
[![Go Report Card](https://goreportcard.com/badge/github.com/yeqown/go-watcher)](https://goreportcard.com/report/github.com/yeqown/go-watcher) [![GoReportCard](https://godoc.org/github.com/yeqown/go-watcher?status.svg)](https://godoc.org/github.com/yeqown/go-watcher)

Golang编写的热重载工具，自定义命令，支持监视文件及路径配置，环境变量配置。这是一个重复的轮子～

### 安装使用

```go
go install github.com/yeqown/go-watcher/cmd/go-watcher
```

or

```sh
brew tab yeqown/go-watcher
brew install go-watcher
```

### 命令行
```sh
➜  go-watcher git:(master) ✗ ./go-watcher -h   
NAME:
   go-watcher - A new cli application

USAGE:
   go-watcher [global options] command [command options] [arguments...]

VERSION:
   2.0.0

AUTHOR:
   yeqown@gmail.com

COMMANDS:
     init     generate a config file to specified postion
     run      execute a command, and watch the files, if any change to these files, the command will reload
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### 配置文件

```yml
watcher:                   # 监视器配置
  duration: 2000           # 文件修改时间间隔，只有高于这个间隔才回触发重载
  included_filetypes:      # 监视的文件扩展类型
  - .go                    # 
  excluded_regexps:        # 不被监视更改的文件正则表达式
  - ^.gitignore$
  - '*.yml$'
  - '*.txt$'
additional_paths: []       # 除了当前文件夹需要额外监视的文件夹
excluded_paths:            # 不需要监视的文件名，若为相对路径，只能对于当前路径生效
- vendor
- .git
envs:                      # 额外的环境变量
- GOROOT=/path/to/your/goroot
- GOPATH=/path/to/your/gopath
```

### 使用范例日志

```sh
➜  go-watcher git:(master) ✗ ./package/osx/go-watcher run -e "make" -c ./config.yml
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/cmd) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/cmd/go-watcher) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/internal) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/internal/command) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/internal/log) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/internal/testdata) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/internal/testdata/exclude) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/internal/testdata/testdata_inner) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/package) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/package/archived) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/package/linux) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/package/osx) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/resources) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/utils) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/utils/testdata) is under watching
[INFO] directory (/Users/yeqown/Projects/opensource/go-watcher/utils/testdata/testdata_inner) is under watching
rm -fr package
go build -o package/osx/go-watcher cmd/go-watcher/main.go
GOOS=linux GOARCH=amd64 go build -o package/linux/go-watcher cmd/go-watcher/main.go
mkdir -p package/archived
tar -zcvf package/archived/go-watcher.osx.tar.gz package/osx
a package/osx
a package/osx/go-watcher
tar -zcvf package/archived/go-watcher.linux.tar.gz package/linux
a package/linux
a package/linux/go-watcher
[INFO] command executed done!
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/osx/go-watcher) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/osx) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/linux/go-watcher) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/linux) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/archived/go-watcher.linux.tar.gz) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/archived) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/VERSION) is skipped, not target filetype
[INFO] [/Users/yeqown/Projects/opensource/go-watcher/cmd/go-watcher/main.go] changed
rm -fr package
mkdir -p package/osx
mkdir -p package/linux
echo "2.0.0" > VERSION
cp VERSION package/osx
cp VERSION package/linux
go build -o package/osx/go-watcher cmd/go-watcher/main.go
GOOS=linux GOARCH=amd64 go build -o package/linux/go-watcher cmd/go-watcher/main.go
mkdir -p package/archived
tar -zcvf package/archived/go-watcher.osx.tar.gz package/osx
a package/osx
a package/osx/go-watcher
a package/osx/VERSION
tar -zcvf package/archived/go-watcher.linux.tar.gz package/linux
a package/linux
a package/linux/go-watcher[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/osx) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/linux) is skipped, not target filetype

a package/linux/VERSION
[INFO] command executed done!
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/osx) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package/archived) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/VERSION) is skipped, not target filetype
[INFO] (/Users/yeqown/Projects/opensource/go-watcher/package) is skipped, not target filetype
^C[INFO] quit signal captured!
[INFO] go-watcher exited
➜  go-watcher git:(master) ✗ 
```