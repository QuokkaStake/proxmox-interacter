package app

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleListDisks(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got list disks query")

	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	template, err := a.TemplateManager.Render("storages", clusters)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering storages template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
