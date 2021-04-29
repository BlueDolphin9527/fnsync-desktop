package client

import (
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
	"github.com/rs/zerolog/log"
)

func (ui *sleepNotifier) startSubscribe() {
	notifierCh := notifier.GetInstance().Start()

	for {
		activity := <-notifierCh

		switch activity.Type {
		case notifier.Awake:
			log.Warn().Msg("Machine awake from sleep")
			go StartHandshake()
		case notifier.Sleep:
			log.Warn().Msg("Machine sleeping")
		}

	}
}
