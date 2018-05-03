# 冰岩在线微信公用微服务
## 基本信息
+ 开发语言: golang
+ 来源：[子豪的微信服务](https://github.com/ZhihaoJun/wxtoken) 的修订版本


## 限制请求來源IP
AccessToken 和 JsApiTicket 接口限制来源 IP ，來自其他IP的请求将直接被拒绝：

主要是通过nginx配置拒绝其他IP请求


signature 接口不限制 IP ，需加上url参数

## API
### GET /service/resources/AccessToken
返回数据范例：
```json{"access_token":"OSWe_yetr-0_0duiXXjlpJJd_sbdvG3LZsvBmy_I8tJVjC5psjTPlyTthSpOhqekTRZ9OShBRq1VHRQ2gfY7v6g5MR2n26H0EO1wRZF_oxZPRNgoz1VVcWoj5wnJqOMQBGXcABAKPG"}
{"access_token":"ACCESS_TOKEN"}
```
| 键 | 类型 | 值 |
|:---:|:---:|:---:|
| access\_token | string | 公众号共用的微信access\_token。 |
### GET /service/resources/JsApiTicket
返回数据范例：
```json
{"jsapi_ticket":"HoagFKDcsGMVCIY2vOjf9kVPAFVbScgB_XXLPSUtSA2hLD6XtsAi3ajdMrzR9a5nGGiFvtLkyWHhkVEzyy7E0A"}
```
| 键 | 类型 | 值 |
|:---:|:---:|:---:|
| jsapi\_ticket | string | 公众号共用的微信jsapi\_ticket。 |

### GET /service/resources/signature?url=\<调用Js接口的完整URL>

返回数据范例：

```
{
    "timestamp": 1499742361,
    "nonce_str": "6161b71f1cf40740",
    "signature": "e11762530337be5e739703539a657c49c8c37711",
    "appId": "<appid>"   
}
```

| 键        | 类型   | 值                                     |
| --------- | ------ | -------------------------------------- |
| timestamp | 时间戳 | 签名形成时候时间戳。                   |
| nonce_str | string | 用于形成签名的随机字符串，共16个字符。 |
| signature | string | 生成的签名。                           |
| appId     | string | 公众号的appid                          |

## docker-compse.yml 简介

若 APPID 和 APP_SECRET修改，只需修改docker-compose.yml中的环境变量即可

```yaml
version: "3"

services:
    web:
      image: <your-image>
      volumes:
        - /etc/localtime:/etc/localtime:ro
      ports:
        - "6000:6000"
      environment:
        - TZ=Asia/Shanghai
        - WXTOKEN_ADDR=:6000
        - WXTOKEN_APPID=<appid>
        - WXTOKEN_APPSECRET=<appsecret>
      container_name: weixin_service_web
      entrypoint: ["bin/main"]

```

| 键        值          | 描述                           |
| --------- | -------------------------------------- |
| WXTOKEN_ADDR | 服务的docker容器中的本地地址，格式：\<host\>:\<port\> |
| WXTOKEN_APPID | 公众号的appid |
| WXTOKEN_APPSECRET | 公众号的app_secret           |

## 镜像制作

Dockerfile

```dockerfile
FROM golang:1.10

WORKDIR /app
COPY ./src src
ENV GOPATH "/app"
RUN go build -o bin/main src/main.go && \
    rm -r src
ENTRYPOINT [ "bin/main" ]

```

注：

+ 由于src下已经使用 dep 工具下载好需要的依赖，故 Dockerfile 中无需解决依赖问题，直接 build 之后删除源码即可

## 代码简介

#### 1. 缓存

+ 由于存储的需求不大，故此处没有使用redis存储而是直接存储于全局变量中，以加快速度，同事应对并发。
+ 通过读写锁来限制缓存访问。
+ 通过定时器，以及 go 的管道的组合来顺序更新 access_token 和 jsapi_ticket

#### 2. 视图

+ 使用 echo 框架

具体请看代码