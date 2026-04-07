# tieba-sign

百度贴吧自动签到工具。

## 功能特性

- 自动获取用户关注的所有贴吧
- 自动签到所有未签到的贴吧
- 并发签到（4 个 goroutine）
- 可配置日志级别
- 调试模式下保存请求/响应日志

## 技术栈

- **语言**: Go 1.24.1+
- **依赖**: github.com/pelletier/go-toml v1.9.

## 目录结构

```
tieba-sign/
├── src/                          # Go 源代码
│   ├── main.go                   # 程序入口
│   ├── config/config.go          # 配置文件加载
│   ├── client/client.go         # HTTP 客户端
│   ├── tieba/                    # 贴吧相关功能
│   │   ├── forum.go              # 论坛数据结构
│   │   ├── fetcher.go           # 获取关注贴吧、TBS、签到
│   │   ├── filter.go             # 过滤未签到贴吧
│   │   └── signer.go             # 签到器
│   ├── log/logger.go            # 日志系统
│   ├── go.mod
│   └── go.sum
└── packaging/                    # 打包配置
    ├── archlinux/                # Arch Linux AUR
    └── common/systemd/           # systemd 服务文件
```

## 前置要求

- Go 1.24.1 或更高版本
- 百度账号 BDUSS 和 STOKEN（获取方式见下文）

## 获取 BDUSS 和 STOKEN

1. 使用浏览器登录百度贴吧 (https://tieba.baidu.com)
2. 打开开发者工具（F12），切换到 Network 面板
3. 刷新页面，找到任意 tieba.baidu.com 的请求
4. 查看请求头中的 Cookie，找到 `BDUSS` 和 `STOKEN` 的值

## 配置

配置文件位于 `~/.config/tieba-sign/config.toml`（如未设置 `XDG_CONFIG_HOME`）。

### 创建配置文件

```bash
mkdir -p ~/.config/tieba-sign
```

创建 `config.toml` 文件，内容如下：

```toml
BDUSS = "your-bduss-value"
STOKEN = "your-stoken-value"

[log]
level = "info"  # 可选: debug, info, notice, warning, error
```

### 配置说明

| 字段      | 描述                  | 必填 |
| --------- | --------------------- | ---- |
| BDUSS     | 百度贴吧登录凭证      | 是   |
| STOKEN    | 百度贴吧安全令牌      | 是   |
| log.level | 日志级别，默认 notice | 否   |

### 日志级别

从低到高：debug, info, notice, warning, error, crit, alert, emerg

## 快速开始

### 编译运行

```bash
cd src
go build -o tieba-sign .
./tieba-sign
```

或直接运行：

```bash
cd src
go run .
```

## 运行效果

首次运行输出示例：

```
NOTICE: 正在获取关注贴吧列表
NOTICE: 待签到: 5/10
INFO: 正在签到: 魔兽世界
INFO: 正在签到: 英雄联盟
INFO: 正在签到: 原神
INFO: 正在签到: 崩坏星穹铁道
INFO: 正在签到: 米哈游
NOTICE: 签到完成: 魔兽世界
NOTICE: 签到完成: 英雄联盟
NOTICE: 签到完成: 原神
NOTICE: 签到完成: 崩坏星穹铁道
NOTICE: 签到完成: 米哈游
NOTICE: 签到完成, 失败: 0/5
```

## 调试模式

启用调试模式可保存请求/响应日志：

```toml
[log]
level = "debug"
```

日志保存在 `~/.local/state/tieba-sign/log/` 目录下：

```
~/.local/state/tieba-sign/log/
├── forums.json                  # 关注贴吧列表响应
├── tbs.json                     # TBS 响应
└── sign-forum/
    └── 魔兽世界/
        ├── req-body             # 签到请求体
        └── resp-body.json       # 签到响应
```

## Arch Linux 安装

使用 AUR 助手（如 yay、paru）安装：

```bash
# 使用 yay
yay -S tieba-sign

# 使用 paru
paru -S tieba-sign
```

或手动构建：

```bash
cd packaging/archlinux/src
makepkg -s
```

### 配置服务定时运行

1. 创建配置文件：

```bash
mkdir -p ~/.config/tieba-sign
nano ~/.config/tieba-sign/config.toml
```

2. 启用用户定时器：

```bash
systemctl --user enable tieba-sign.timer
systemctl --user start tieba-sign.timer
```

3. 查看定时器状态：

```bash
systemctl --user list-timers
```

### 手动运行签到

```bash
systemctl --user run tieba-sign.service
```

## 代码说明

### 主程序流程 (main.go)

1. 加载配置文件
2. 设置日志级别
3. 创建 HTTP 客户端（携带 BDUSS、STOKEN Cookie）
4. 获取关注贴吧列表
5. 过滤出未签到的贴吧
6. 启动 4 个并发 goroutine 进行签到
7. 等待所有签到完成，输出结果

### HTTP 客户端 (client/client.go)

- 使用 cookiejar 管理 Cookie
- 预设 User-Agent 和 Host 请求头
- 提供 GET/POST 方法
- 记录调试日志

### 贴吧模块 (tieba/)

- `ForumsFetch`: 获取关注贴吧列表
- `TBSFetch`: 获取 TBS 防跨站请求令牌
- `ForumSign`: 执行签到请求
- `ForumsFilterUnsigned`: 过滤出未签到的贴吧

### 日志系统 (log/logger.go)

- 多级别日志：EMERG, ALERT, CRIT, ERR, WARNING, NOTICE, INFO, DEBUG
- ERROR 及以下级别输出到 stderr，其他输出到 stdout
- 支持写入文件保存请求/响应数据

## 常见问题

### 签到失败

1. 检查 BDUSS 和 STOKEN 是否有效
2. 启用 debug 模式查看详细日志
3. 查看 `~/.local/state/tieba-sign/log/` 下的响应文件

### 配置文件找不到

确保配置文件路径正确：

- 默认: `~/.config/tieba-sign/config.toml`
- 如设置 `XDG_CONFIG_HOME`: `$XDG_CONFIG_HOME/tieba-sign/config.toml`

### 日志文件找不到

调试日志保存在：

- 默认: `~/.local/state/tieba-sign/log/`
- 如设置 `XDG_STATE_HOME`: `$XDG_STATE_HOME/tieba-sign/log/`

## 相关链接

- [百度贴吧](https://tieba.baidu.com)
