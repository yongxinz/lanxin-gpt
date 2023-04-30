# lanxin-gpt

蓝信机器人联动 ChatGPT

## 使用

```
POST http://127.0.0.1:8080/chat

{
    'msg': 'Hello'
}
```

程序会调用 ChatGPT API 查询结果，然后调用蓝信 WebHook API 将结果推送到蓝信群

## 开发

1、更新配置文件

```
$ cp .env.example .env
```

2、启动程序

```
$ go run main.go
```

## 扩展

如果接入蓝信智能机器人，可能使用场景会更丰富一些，但我不打算对这个项目做进一步开发了，感兴趣的小伙伴可以在此基础上继续去尝试
