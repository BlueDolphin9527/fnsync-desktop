# fnsync-desktop
FnSync非官方实现第三方 MacOS 桌面版

<img src="./doc/preview.jpg" width="720" >


官方APP：

https://www.coolapk.com/apk/269031

https://play.google.com/store/apps/details?id=holmium.fnsync

官方客户端(windows版)：

https://gitee.com/holmium/fnsync/releases

## feature
* 接收手机发送的文本
* 同步剪贴板文本到手机
* 发送文本到指定手机 

## prerequisite
* Go 1.6+
* Node.js 12+
* Xcode 12
## how to build

```
go get github.com/wailsapp/wails/v2/cmd/wails@v2.0.0-alpha.65
# for Go 1.7+
# go install github.com/wailsapp/wails/v2/cmd/wails@v2.0.0-alpha.65  
cd src
go mod vendor
wails build
```