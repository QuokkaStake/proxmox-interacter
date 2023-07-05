package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func (a *App) HandleStartContainer(c tele.Context) error {
	if len(strings.Split(c.Text(), " ")) == 1 {
		return a.HandleHelp(c)
	}

	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got start container query")

	args := strings.SplitN(c.Text(), " ", 2)
	command, args := args[0], args[1:]

	if len(args) != 1 {
		return c.Reply(fmt.Sprintf("Usage: %s <container name or ID>", command))
	}

	container, err := a.ProxmoxManager.StartContainer(args[0])
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error starting container: %s", err))
	}

	template, err := a.TemplateManager.Render("container_action", ContainerActionRender{
		Container: *container,
		Action:    "started",
	})
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container_action template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
