
# 无忧秘书智脑Bill模块

无忧秘书智Bill相关模块。

![无忧秘书智脑](https://umi-intelligence.oss-cn-shenzhen.aliyuncs.com/static/website/screenshot-ai.umi6.com-2024.03.13-10_15_32.png)

## 功能特色
- RESTful API
- Gorm
- Swagger
- logging
- Jwt-go
- Gin
- Graceful restart or stop (fvbock/endless)
- App configurable
- Cron
- Redis

## 准备工作

请按照以下步骤准备环境：

- **配置项**：补全conf/app.ini下相关配置
- **端口放行**：确保以下端口已放行：28070。


## 开始
    拉取项目：git clone https://github.com/ymzn3820/umi_platform_bill_module.git
    构建：docker compose -f production.yaml build
    运行：docker compose -f production.yaml up -d
    检查：docker logs -f {containerId}, 如果没有报错信息，则运行成功。

## 导航
| 模块名称 | 链接 | 介绍|
| -------- | ---- |---- |
| 前端PC | [umi_platform_frontend](https://github.com/ymzn3820/umi_platform_frontend) | PC端前段代码仓库地址|
| 小程序端 | [umi_platform_mini_program](https://github.com/ymzn3820/umi_platform_mini_program) |小程序端代码仓库地址|
| H5端 | [umi_platform_h5](https://github.com/ymzn3820/umi_platform_h5) |H5端代码仓库地址|
| 支付模块 | [umi_platform_pay_module](https://github.com/ymzn3820/umi_platform_pay_module) |支付模块代码仓库地址|
| 用户模块 | [umi_platform_user_module](https://github.com/ymzn3820/umi_platform_user_module) |用户模块代码仓库地址|
| Chat模块 | [umi_platform_chat_module](https://github.com/ymzn3820/umi_platform_chat_module) |Chat模块代码仓库地址|

[返回引导页](https://github.com/ymzn3820/umi_platform_pay_module)

## License

BSD 3-Clause License

