package app

import (
	"fmt"
	"strings"

)

func (a *App) HandleRestartContainer(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got restart container query")

	args := strings.SplitN(c.Text(), " ", 2)
	command, args := args[0], args[1:]

	if len(args) != 1 {
		return c.Reply(fmt.Sprintf("Usage: %s <container name or ID>", command))
	}

	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	container, _, found := clusters.FindContainer(args[0])
	if !found {
		return a.BotReply(c, "Container is not found!")
	}

	template, err := a.TemplateManager.Render("container_action", ContainerActionRender{
		Container: *container,
		Action:    "restart",
	})
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container_action template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	approveButton := menu.Data("✅Approve", fmt.Sprintf("%s%s", CallbackPrefixRestart, args[0]))
	cancelButton := menu.Data("❌Cancel", fmt.Sprintf("%s%s", CallbackPrefixCancelRestart, args[0]))

	menu.Inline(
		menu.Row(approveButton, cancelButton),
	)

	if err := a.BotReply(c, template, menu); err != nil {
		a.Logger.Error().Err(err).Msg("Could not delete message!")
	}

	return a.Bot.Delete(c.Message())

	// container, err := a.ProxmoxManager.RestartContainer(args[0])
	// if err != nil {
	//	return a.BotReply(c, fmt.Sprintf("Error starting container: %s", err))
	//}
}

func (a *App) HandleDoRestartContainer(c tele.Context, data string) error {
	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	container, _, found := clusters.FindContainer(data)
	if !found {
		return a.BotReply(c, "Container is not found!")
	}

	if _, err := a.Bot.Edit(c.Message(), c.Message().Text, &tele.ReplyMarkup{}); err != nil {
		a.Logger.Error().Err(err).Msg("Could not edit message!")
	}

	template, err := a.TemplateManager.Render("container_action_do", ContainerActionRender{
		Container: *container,
		Action:    "restarted",
	})
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container_action template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}

func (a *App) HandleDoCancelRestartContainer(c tele.Context, data string) error {
	if _, err := a.Bot.Edit(c.Message(), c.Message().Text, &tele.ReplyMarkup{}); err != nil {
		a.Logger.Error().Err(err).Msg("Could not edit message!")
	}
	return a.BotReply(c, "Restarting is cancelled!")
}
