# 程序员周报助手
公司要求提交周报，提交的时间窗口是周六全天，有时周六休息，容易搞忘，而且电脑不在身边，也不好处理。于是该项目诞生了。。。

注意：项目在 Windows11 上测试通过，其它系统未适配

## 使用说明
    程序运行需要指定几个参数：-paths "C:/dev/xx C:/dev/mine/" -command "log dev --since='1 week ago' --author='xxx' --no-merges --pretty=format:'%%s'"
    paths 为各个项目代码的根目录，多个目录空格分开；如果多个项目在同一目录下，可以只写父级目录，但路径需要以'/'结尾
    command 为 git 命令: dev 是代码分支，since 是拉取 git 日志的时间范围，author 是你 git 设置的用户名

## 最佳实践
    将项目打包后，创建 bat 脚本，win 下创建定时任务，执行 bat 脚本。
    脚本执行后会自动拉取日志，然后调用 AI 接口生成日志内容，并发送到推送到微信。

    日志内容的质量主要依赖于两点：
    1. git commit 的记录内容
    2. AI 提示词

    AI 接入的是智谱清言的自建智能体，总共 5000 次的免费调用，写周报够用了。
    env 文件中的 AI_KEY AI_SECRET ASSISTANT_ID，需要到智谱去申请。
    
    微信推送使用的 Server酱 https://sct.ftqq.com/
    env 文件中的 SEND_KEY，需要到 Server酱 去申请。

    请确保
        .env
        weekly-report.exe
        weekly-report.bat
    三个文件在同一目录

## MakeFile

打包应用
```bash
make build
```
 ## Todo
- [ ] 目前 git log dev... 固定为 dev 分支，考虑不同项目支持不同分支
- [ ] 支持更广泛的 AI 调用，支持切换
- [ ] 支持飞书app，使用飞书群机器人来推送消息，飞书支持更多的消息类型和更多的消息条数
- [ ] Mac 适配，但是没有 Mac，这个怎么弄啊...

## 资源
- [Server酱](https://sct.ftqq.com/) 微信推送
- [CURSOR](https://www.cursor.com/) The AI Code Editor

## 后续
- [x] 已支持 DeepSeek，可在 env 文件中切换，DeepSeek 需要自定义 prompt