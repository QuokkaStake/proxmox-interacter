package app

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleHelp(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	template, err := a.TemplateManager.Render("help", a.Version)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering help template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
