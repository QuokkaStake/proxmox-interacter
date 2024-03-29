package app

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleContainerInfo(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got container info query")

	args := strings.SplitN(c.Text(), " ", 2)
	command, args := args[0], args[1:]

	if len(args) != 1 {
		return c.Reply(fmt.Sprintf("Usage: %s <container name or ID>", command))
	}

	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	container, cluster, err := clusters.FindContainer(args[0])
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

	render := ContainerInfoRender{
		Container: *container,
	}

	if config, err := a.ProxmoxManager.GetContainerConfig(*container, cluster); err != nil {
		a.Logger.Error().Err(err).Msg("Error fetching container config")
		render.ConfigError = err
	} else {
		render.Config = config
	}

	template, err := a.TemplateManager.Render("container", render)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
