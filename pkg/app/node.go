package app

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (a *App) HandleNodeInfo(c tele.Context) error {
	a.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got node info query")

	args := strings.SplitN(c.Text(), " ", 2)
	command, args := args[0], args[1:]

	if len(args) != 1 {
		return c.Reply(fmt.Sprintf("Usage: %s <node name>", command))
	}

	clusters, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	for _, cluster := range clusters {
		if cluster.Error != nil {
			continue
		}

		for _, node := range cluster.Nodes {
			if node.Node.Node == args[0] || node.Node.ID == args[0] {
				template, err := a.TemplateManager.Render("node", node)
				if err != nil {
					a.Logger.Error().Err(err).Msg("Error rendering node template")
					return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
				}

				return a.BotReply(c, template)
			}
		}
	}

	return a.BotReply(c, "Node is not found")
}
