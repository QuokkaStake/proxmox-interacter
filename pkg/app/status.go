package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"main/pkg/types"
)

func (a *App) HandleStatus(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got status query")

	resources, err := a.ProxmoxManager.GetResources()
	if err != nil {
		return a.BotReply(c, "Error fetching nodes status")
	}

	nodes := make([]types.Node, 0)

	for _, resource := range resources {
		resourceNodes, err := types.ParseNodesFromResponse(resource)
		if err != nil {
			return a.BotReply(c, "Error parsing nodes status")
		}

		nodes = append(nodes, resourceNodes...)
	}

	template, err := a.TemplateManager.Render("status", nodes)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering status template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
