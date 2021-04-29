package entity

import (
	"github.com/gen2brain/beeep"
	"github.com/rs/zerolog/log"
)

type UINotifyMsg struct {
	Title   string
	Message string
}

func (ui UINotifyMsg) Execute() {
	err := beeep.Notify(ui.Title, ui.Message, "../assets/appicon.png")
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}
