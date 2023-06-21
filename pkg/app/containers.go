package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleListContainers(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list containers query")

	containers, err := a.ProxmoxManager.GetNodesWithContainers()
	if err != nil {
		return a.BotReply(c, "Error fetching containers")
	}

	template, err := a.TemplateManager.Render("containers", containers)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering containers template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
