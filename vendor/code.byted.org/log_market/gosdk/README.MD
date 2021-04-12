#### 概述
LogAgent SDK 将日志写入LogAgent.
#### DEMO-1 use sdk
```
defer gosdk.GracefulExit()
psm := "p.s.m"
m := &gosdk.Msg{
    Msg:[]byte("hello world!!!"),
    Tags:map[string]string{
        "tagKey":"tagValue",
    },
}
gosdk.Send(psm,m)
```
#### DEMO-2 use default sender
```
psm := "p.s.m"
sender := gosdk.NewDefaultSender(psm)

defer sender.Stop()

m := &gosdk.Msg{
    Msg: []byte("hello world!!!"),
    Tags: map[string]string{
        "tagKey": "tagValue",
    },
}
if err := sender.Send(psm, m); err != nil {
    fmt.Println(err.Error())
}
```
#### DEMO-3 use custom sender
```
psm := "p.s.m"
config := gosdk.NewDefaultSenderConfig(psm)
config.ChannelSize = 8 * 1024
config.SocketWriteTimeoutMs = 500
sender := gosdk.NewSender(config)

defer sender.Stop()

m := &gosdk.Msg{
    Msg: []byte("hello world!!!"),
    Tags: map[string]string{
        "tagKey": "tagValue",
    },
}
if err := sender.Send(psm, m); err != nil {
    fmt.Println(err.Error())
}
```