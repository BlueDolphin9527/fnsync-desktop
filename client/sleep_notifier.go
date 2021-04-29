package client

var SleepNotifier *sleepNotifier = &sleepNotifier{}

type sleepNotifier struct {
}

func (ui *sleepNotifier) StartSubscribe() {
	ui.startSubscribe()
}
