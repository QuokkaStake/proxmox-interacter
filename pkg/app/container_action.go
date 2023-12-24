package app

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type containerAction struct {
	action                   string
	doneAction               string
	shouldContainerBeStarted bool
	actionPrefix             string
	cancelActionPrefix       string
	function                 func(string) error
}

func (a *App) getAction(action string) containerAction {
	actions := map[string]containerAction{
		"restart": {
			action:                   "restart",
			doneAction:               "restarted",
			shouldContainerBeStarted: true,
			actionPrefix:             CallbackPrefixRestart,
			cancelActionPrefix:       CallbackPrefixCancelRestart,
			function:                 a.ProxmoxManager.RestartContainer,
		},
		"stop": {
			action:                   "stop",
			doneAction:               "stopped",
			shouldContainerBeStarted: true,
			actionPrefix:             CallbackPrefixStop,
			cancelActionPrefix:       CallbackPrefixCancelStop,
			function:                 a.ProxmoxManager.StopContainer,
		},
		"start": {
			action:                   "start",
			doneAction:               "started",
			shouldContainerBeStarted: false,
			actionPrefix:             CallbackPrefixStart,
			cancelActionPrefix:       CallbackPrefixCancelStart,
			function:                 a.ProxmoxManager.StartContainer,
		},
	}

	return actions[action]
}

func (a *App) HandleContainerAction(actionName string) func(c tele.Context) error {
	action := a.getAction(actionName)

	return func(c tele.Context) error {
		a.Logger.Info().
			Str("sender", c.Sender().Username).
			Str("text", c.Text()).
			Str("action", action.action).
			Msg("Got container action query")

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
			Action:    action.action,
		})
		if err != nil {
			a.Logger.Error().Err(err).Msg("Error rendering container_action template")
			return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
		}

		menu := &tele.ReplyMarkup{ResizeKeyboard: true}

		approveButton := menu.Data("✅Approve", fmt.Sprintf("%s%s", action.actionPrefix, args[0]))
		cancelButton := menu.Data("❌Cancel", fmt.Sprintf("%s%s", action.cancelActionPrefix, args[0]))

		menu.Inline(
			menu.Row(approveButton, cancelButton),
		)

		return a.BotReply(c, template, menu)

		// container, err := a.ProxmoxManager.RestartContainer(args[0])
		// if err != nil {
		//	return a.BotReply(c, fmt.Sprintf("Error starting container: %s", err))
		//}
	}
}

func (a *App) HandleDoContainerAction(actionName string) func(c tele.Context, data string) error {
	action := a.getAction(actionName)

	return func(c tele.Context, data string) error {
		clusters, err := a.ProxmoxManager.GetNodes()
		if err != nil {
			return a.BotReply(c, fmt.Sprintf("Error fetching nodes: %s", err))
		}

		container, _, found := clusters.FindContainer(data)
		if !found {
			return a.BotReply(c, "Container is not found!")
		}

		if container.IsRunning() && !action.shouldContainerBeStarted {
			return c.Reply("Container is already running!")
		} else if !container.IsRunning() && action.shouldContainerBeStarted {
			return c.Reply("Container is not running!")
		}

		if err := action.function(data); err != nil {
			return c.Reply(fmt.Sprintf("Error %s container: %s", action.doneAction, err))
		}

		if _, err := a.Bot.Edit(c.Message(), c.Message().Text, &tele.ReplyMarkup{}); err != nil {
			a.Logger.Error().Err(err).Msg("Could not edit message!")
		}

		template, err := a.TemplateManager.Render("container_action_do", ContainerActionRender{
			Container: *container,
			Action:    action.doneAction,
		})
		if err != nil {
			a.Logger.Error().Err(err).Msg("Error rendering container_action template")
			return c.Reply(fmt.Sprintf("Error rendering template: %s", err))
		}

		if err := a.BotReply(c, template); err != nil {
			a.Logger.Error().Err(err).Msg("Could not delete message!")
		}

		return a.Bot.Delete(c.Message())
	}
}

func (a *App) HandleDoCancelContainerAction(actionName string) func(c tele.Context, data string) error {
	action := a.getAction(actionName)

	return func(c tele.Context, data string) error {
		if _, err := a.Bot.Edit(c.Message(), c.Message().Text, &tele.ReplyMarkup{}); err != nil {
			a.Logger.Error().Err(err).Msg("Could not edit message!")
		}
		return a.BotReply(c, fmt.Sprintf(
			"%s%s is cancelled!",
			strings.ToUpper(action.action[:1]),
			action.action[1:],
		))
	}
}
