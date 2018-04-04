# gowatch

Go 程序热重载工具，提升开发效率

通过监听当前目录下的相关文件变动，进行实时执行命令

### 安装使用

```go
go get github.com/yeqown/gowatch 
go install github.com/yeqown/gowatch
```

### 命令行

`gowatch init`

`gowatch run [command] [...cmdArgs]`

### 配置文件

`gowatch init` // 初始化配置文件

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
  - 

# 当前目录下需要排除的文件夹
exclude_paths:
  - vendor
  - testdata
  # - abspath also ok

```

>Forked and rewrite from [gowatch](https://github.com/silenceper/gowatch)
