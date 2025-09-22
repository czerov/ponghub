<div align="center">

# [![PongHub](imgs/band.png)](https://health.ch3nyang.top)

🌏 [Live Demo](https://health.ch3nyang.top) | 📖 [English](README.md)

</div>

## 简介

PongHub 是一个开源的服务状态监控网站，旨在帮助用户监控和验证服务的可用性。它支持

- **🕵️ 零侵入监控** - 无需改动代码即可实现全功能监控
- **🚀 一键部署** - 通过 Actions 自动构建，一键部署至 Github Pages
- **🌐 全平台支持** - 兼容 OpenAI 等公共服务及私有化部署
- **🔍 多端口探测** - 单服务支持同时监控多个端口状态
- **🤖 智能响应验证** - 精准匹配状态码及正则表达式校验响应体
- **🛠️ 自定义请求引擎** - 自由配置请求头/体、超时和重试策略
- **🔒 SSL 证书监控** - 自动检测 SSL 证书过期并发送通知
- **📊 实时状态展示** - 直观的服务响应时间、响应状态记录
- **⚠️ 异常告警通知** - 利用 GitHub Actions 实现异常告警通知

![浏览器截图](imgs/browser_CN.png)

## 快速开始

1. Star 并 Fork [PongHub](https://github.com/WCY-dt/ponghub)

2. 修改根目录下的 [`config.yaml`](config.yaml) 文件，配置你的服务检查项

3. 修改根目录下的 [`CNAME`](CNAME) 文件，配置你的自定义域名

   > 如果你不需要自定义域名，请删除 `CNAME` 文件

4. 提交修改并推送到你的仓库，GitHub Actions 将自动更新，无需干预

> [!TIP]
> 默认情况下，GitHub Actions 每 30 分钟运行一次。如果你需要更改运行频率，请修改 [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml) 文件中的 `cron` 表达式。
> 
> 请不要将频率设置过高，以免触发 GitHub 的限制。

> [!IMPORTANT]
> 如果 GitHub Actions 未正常自动触发，手动触发一次即可。
> 
> 请注意打开 GitHub Pages，并开启 GitHub Actions 的通知权限。

## 配置说明

### 基本配置

配置文件 `config.yaml` 的格式如下：

| 字段                                  | 类型  | 描述                        | 必填 | 备注                             |
|-------------------------------------|-----|---------------------------|----|--------------------------------|
| `display_num`                       | 整数  | 首页显示的服务数量                 | ✖️ | 默认 72 个                        |
| `timeout`                           | 整数  | 每次请求的超时时间，单位为秒            | ✖️ | 单位为秒，默认 5 秒                    |
| `max_retry_times`                   | 整数  | 请求失败时的重试次数                | ✖️ | 默认 2 次                         |
| `max_log_days`                      | 整数  | 日志保留天数，超过此天数的日志将被删除       | ✖️ | 默认 3 天                         |
| `cert_notify_days`                  | 整数  | SSL 证书过期前通知的天数            | ✖️ | 默认 7 天                         |
| `services`                          | 数组  | 服务列表                      | ✔️ |                                |
| `services.name`                     | 字符串 | 服务名称                      | ✔️ |                                |
| `services.endpoints`                | 数组  | 端口列表                      | ✔️ |                                |
| `services.endpoints.url`            | 字符串 | 请求的 URL                   | ✔️ |                                |
| `services.endpoints.method`         | 字符串 | 请求的 HTTP 方法               | ✖️ | 支持 `GET`/`POST`/`PUT`，默认 `GET` |
| `services.endpoints.headers`        | 对象  | 请求头内容                     | ✖️ | 键值对形式，支持自定义请求头                 |
| `services.endpoints.body`           | 字符串 | 请求体内容                     | ✖️ | 仅在 `POST`/`PUT` 请求时使用          |
| `services.endpoints.status_code`    | 整数  | 响应体期望的 HTTP 状态码（默认 `200`） | ✖️ | 默认 `200`                       |
| `services.endpoints.response_regex` | 字符串 | 响应体内容的正则表达式匹配             | ✖️ |                                |

下面是一个示例配置文件：

```yaml
display_num: 72
timeout: 5
max_retry_times: 2
max_log_days: 3
cert_notify_days: 7
services:
  - name: "GitHub API"
    endpoints:
      - url: "https://api.github.com"
      - url: "https://api.github.com/repos/wcy-dt/ponghub"
        method: "GET"
        headers:
          Content-Type: application/json
          Authorization: Bearer your_token
        status_code: 200
        response_regex: "full_name"
  - name: "Example Website"
    endpoints:
      - url: "https://example.com/health"
        response_regex: "status"
      - url: "https://example.com/status"
        method: "POST"
        body: '{"key": "value"}'
```

### 特殊参数

ponghub 现已支持强大的参数化配置功能，允许在配置文件中使用多种类型的动态变量，这些变量会在程序运行时实时生成和解析。

<details>
<summary>点击展开查看支持的参数类型</summary>

<div markdown="1">

#### 📅 日期时间参数
使用 `{{%格式}}` 格式定义日期时间参数：

- `{{%Y-%m-%d}}` - 当前日期，格式：2006-01-02（如：2025-09-22）
- `{{%H:%M:%S}}` - 当前时间，格式：15:04:05（如：17:30:45）
- `{{%s}}` - Unix时间戳（如：1727859600）
- `{{%Y}}` - 当前年份（如：2025）
- `{{%m}}` - 当前月份，格式：01-12
- `{{%d}}` - 当前日期，格式：01-31
- `{{%H}}` - 当前小时，格式：00-23
- `{{%M}}` - 当前分钟，格式：00-59
- `{{%S}}` - 当前秒数，格式：00-59
- `{{%B}}` - 完整月份名称（如：September）
- `{{%b}}` - 简短月份名称（如：Sep）
- `{{%A}}` - 完整星期名称（如：Monday）
- `{{%a}}` - 简短星期名称（如：Mon）

#### 🎲 随机数参数

- `{{rand}}` - 生成0-1000000范围的随机数
- `{{rand_int}}` - 生成大范围随机整数
- `{{rand(min,max)}}` - 生成指定范围的随机数
    - 示例：`{{rand(1,100)}}` - 生成1-100之间的随机数
    - 示例：`{{rand(1000,9999)}}` - 生成4位随机数

#### 🔤 随机字符串参数

- `{{rand_str}}` - 生成8位随机字符串（字母+数字）
- `{{rand_str(length)}}` - 生成指定长度的随机字符串
    - 示例：`{{rand_str(16)}}` - 生成16位随机字符串
- `{{rand_str_secure}}` - 生成16位加密安全的随机字符串
- `{{rand_hex(length)}}` - 生成指定长度的十六进制随机字符串
    - 示例：`{{rand_hex(8)}}` - 生成8位十六进制字符串
    - 示例：`{{rand_hex(32)}}` - 生成32位十六进制字符串

#### 🆔 UUID参数

- `{{uuid}}` - 生成标准UUID（带连字符）
    - 示例：`bf3655f7-8a93-4822-a458-2913a6fe4722`
- `{{uuid_short}}` - 生成短UUID（无连字符）
    - 示例：`14d44b7334014484bb81b015fb2401bf`

#### 🌍 环境变量参数

- `{{env(变量名)}}` - 读取环境变量的值
    - 示例：`{{env(API_KEY)}}` - 读取API_KEY环境变量
    - 示例：`{{env(VERSION)}}` - 读取VERSION环境变量
    - 如果环境变量不存在，返回空字符串

环境变量可通过 GitHub Actions 的 Repository Secrets 设置

#### 📊 序列号和哈希参数

- `{{seq}}` - 基于当前时间的序列号（6位数字）
- `{{seq_daily}}` - 每日序列号（自午夜起的秒数）
- `{{hash_short}}` - 短哈希值（6位十六进制）
- `{{hash_md5_like}}` - MD5风格的长哈希值（32位十六进制）

</div>
</details>

下面是一个示例配置文件：

```yaml
services:
  - name: "Parameterized Service"
    endpoints:
        - url: "https://api.example.com/data?date={{%Y-%m-%d}}&rand={{rand(1,100)}}"
        - url: "https://api.example.com/submit"
          method: "POST"
          headers:
            Content-Type: application/json
            X-Request-ID: "{{uuid}}"
          body: '{"session_id": "{{rand_str(16)}}", "timestamp": "{{%s}}"}'
```

### 自定义通知

PongHub 默认利用 GitHub Actions 报错实现异常告警通知。

如果需要自定义通知，可以在根目录下创建 `notify.sh` 脚本，脚本可以读取 `data/notify.txt` 文件中的内容，并通过邮件、短信或其他方式发送通知。如果脚本使用到了环境变量，请确保在 GitHub 仓库的 "Settings" -> "Secrets and variables" -> "Actions" 中正确设置这些变量。

## 本地开发

本项目使用 Makefile 进行本地开发和测试。你可以使用以下命令在本地运行项目：

```bash
make run
```

## 免责声明

[PongHub](https://github.com/WCY-dt/ponghub) 仅用于个人学习和研究，不对程序的使用行为或结果负责。请勿将其用于商业用途或非法活动。
