# fnsync-desktop
非官方实现第三方 MacOS 桌面版FnSync

<img src="./doc/preview.jpg" width="720" >

## feature
* 授受手机发送的文本
* 同步剪贴板文本到手机
* 发送文本到指定同步手机 

## how to build

```
go get github.com/wailsapp/wails/v2/cmd/wails
cd src
go mod tidy
wails build
```