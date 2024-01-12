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

	container, _, scaleInfo, err := clusters.FindContainerToScale(args[0])
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

	template, err := a.TemplateManager.Render("container_scale", ContainerScaleRender{
		Container:   *container,
		ScaleParams: scaleInfo,
	})
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container_scale template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	if !scaleInfo.AnythingChanged(*container) {
		return a.BotReply(c, template)
	}

	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	approveButton := menu.Data("✅Approve", fmt.Sprintf("%s%s", CallbackPrefixScale, args[0]))
	cancelButton := menu.Data("❌Cancel", fmt.Sprintf("%s%s", CallbackPrefixCancelScale, args[0]))

	menu.Inline(
		menu.Row(approveButton, cancelButton),
	)

	return a.BotReply(c, template, menu)
}
