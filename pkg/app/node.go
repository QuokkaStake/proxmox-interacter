package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
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

	nodes, err := a.ProxmoxManager.GetNodes()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
	}

	node, found := utils.Find(nodes, func(node types.Node) bool {
		return node.Node == args[0] || node.ID == args[0]
	})

	if !found {
		return a.BotReply(c, "Node is not found")
	}

	template, err := a.TemplateManager.Render("node", node)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering node template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
