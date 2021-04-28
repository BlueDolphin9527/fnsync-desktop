package clip

import (
	"time"

	"github.com/cxfksword/fnsync-desktop/clipboard/clip/ns"
)

var pb *ns.NSPasteboard
var lastChangeCount = 0
var notifyChangeCh chan string = make(chan string, 1)

func init() {
	go watchChange()
}

func Clear() {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}
	pb.ClearContents()
}

func Set(x string) bool {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}
	pb.ClearContents()
	return pb.SetString(x)
}

func Get() string {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}
	ret := pb.GetString()
	if ret.Ptr() == nil {
		return ""
	} else {
		return ret.String()
	}
}

func watchChange() {
	if pb == nil {
		pb = ns.NSPasteboardGeneralPasteboard()
	}

	timer := time.NewTimer(1 * time.Second)
	for {
		<-timer.C

		changeCount := pb.GetChangeCount()
		// 首次运行，初始化数据
		if lastChangeCount == 0 {
			lastChangeCount = changeCount
		}
		if lastChangeCount != changeCount {
			lastChangeCount = changeCount

			notifyChangeCh <- Get()
		}

		timer.Reset(1 * time.Second)
	}

}

func GetNotifyChangeCh() chan string {
	return notifyChangeCh
}
