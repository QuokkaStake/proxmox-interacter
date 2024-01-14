package app

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleContainerScale(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got container scale query")

	args := strings.SplitN(c.Text(), " ", 2)
	command, args := args[0], args[1:]

	if len(args) != 1 {
		return c.Reply(fmt.Sprintf("Usage: %s <container name or ID>", command))
	}

	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	container, cluster, scaleInfo, err := clusters.FindContainerToScale(args[0])
	if err != nil {
		template, err := a.TemplateManager.Render("container_error", ContainerErrorRender{
			Error:        err,
			ClusterInfos: clusters,
		})
		if err != nil {
			a.Logger.Error().Err(err).Msg("Error rendering container template")
			return c.Reply(fmt.Sprintf("Error rendering template when processing error: %s", err))
		}

		return a.BotReply(c, template)
	}

	config, err := a.ProxmoxManager.GetContainerConfig(*container, cluster)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error fetching container config")
		return a.BotReply(c, fmt.Sprintf("Error fetchign container config: %s", err))
	}

	template, err := a.TemplateManager.Render("container_scale", ContainerScaleRender{
		Container:   *container,
		ScaleParams: scaleInfo,
		Config:      config,
	})
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container_scale template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	if !scaleInfo.AnythingChanged(*container, config) {
		return a.BotReply(c, template)
	}

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	approveButton := menu.Data("✅Approve", fmt.Sprintf("%s%s", CallbackPrefixScale, args[0]))
	cancelButton := menu.Data("❌Cancel", CallbackPrefixCancelScale)

	menu.Inline(
		menu.Row(approveButton, cancelButton),
	)

	return a.BotReply(c, template, menu)
}

func (a *App) HandleDoContainerScale(c tele.Context, data string) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", data).
		Msg("Got container do scale query")

	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	container, cluster, scaleInfo, err := clusters.FindContainerToScale(data)
	if err != nil {
		template, err := a.TemplateManager.Render("container_error", ContainerErrorRender{
			Error:        err,
			ClusterInfos: clusters,
		})
		if err != nil {
			a.Logger.Error().Err(err).Msg("Error rendering container template")
			return c.Reply(fmt.Sprintf("Error rendering template when processing error: %s", err))
		}

		return a.BotReply(c, template)
	}

	config, err := a.ProxmoxManager.GetContainerConfig(*container, cluster)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error fetching container config")
		return a.BotReply(c, fmt.Sprintf("Error fetchign container config: %s", err))
	}

	if err := a.ProxmoxManager.ScaleContainer(cluster, *container, config, scaleInfo); err != nil {
		a.Logger.Error().Err(err).Msg("Error scaling container")
		return a.BotReply(c, fmt.Sprintf("Error scaling container: %s", err))
	}

	if _, err := a.Bot.Edit(c.Message(), c.Message().Text, &tele.ReplyMarkup{}); err != nil {
		a.Logger.Error().Err(err).Msg("Could not edit message!")
	}

	return a.BotReply(c, "Container is scaled!")
}

func (a *App) HandleDoCancelContainerScale(c tele.Context, data string) error {
	if _, err := a.Bot.Edit(c.Message(), c.Message().Text, &tele.ReplyMarkup{}); err != nil {
		a.Logger.Error().Err(err).Msg("Could not edit message!")
	}
	return a.BotReply(c, "Scaling is cancelled!")
}
