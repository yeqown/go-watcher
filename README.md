# gw

Golang编写的热重载工具，自定义命令，支持监视文件及路径配置，环境变量配置。这是一个重复的轮子～

### 安装使用

```go
git clone git@github.com:yeqown/gw.git
cd path/to/gw && go build
mv gw $GOPATH/bin
```

### 命令行

生成配置文件`gw init`

利用gw执行命令 `gw run [command] [...cmdArgs]`，如：
	
	gw run go run main.go -conf ./configs/config.yml

### 配置文件

`gw init` // 初始化配置文件

```yml
# gowatch.yml

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

# 需要排除的文件夹，支持绝对路径和相对路径
exclude_paths:
  - ./vendor
  - ./testdata
  # - abspath is also ok

```

> Forked and rewrite from [gowatch](https://github.com/silenceper/gowatch)
