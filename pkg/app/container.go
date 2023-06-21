package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
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

	containers, err := a.ProxmoxManager.GetContainers()
	if err != nil {
		return a.BotReply(c, fmt.Sprintf("Error fetching containers: %s", err))
	}

	container, found := utils.Find(containers, func(container types.Container) bool {
		// containers IDs are like lxc/XXXX or qemu/XXXX
		return fmt.Sprintf("%d", container.VMID) == args[0] ||
			container.ID == args[0] ||
			container.Name == args[0]
	})

	if !found {
		return a.BotReply(c, "Container is not found")
	}

	template, err := a.TemplateManager.Render("container", container)
	if err != nil {
		a.Logger.Error().Err(err).Msg("Error rendering container template")
		return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
	}

	return a.BotReply(c, template)
}
