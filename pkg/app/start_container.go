package app

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleStartContainer(c tele.Context) error {
	if len(strings.Split(c.Text(), " ")) == 1 {
		return a.HandleHelp(c)
	}

	return a.HandleContainerAction("start")(c)
}
