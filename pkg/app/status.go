package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleStatus(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got status query")

	nodes, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, "Error fetching nodes status")
	}

	template, err := a.TemplateManager.Render("status", nodes)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering status template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
