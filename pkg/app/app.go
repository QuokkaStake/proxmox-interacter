package app

import (
	loggerPkg "main/pkg/logger"
	"main/pkg/proxmox"
	"main/pkg/templates"
	"main/pkg/types"
	"strings"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const MaxMessageSize = 4096

type App struct {
	Config          types.Config
	ProxmoxManager  *proxmox.Manager
	TemplateManager *templates.TemplateManager
	Logger          *zerolog.Logger
	Bot             *tele.Bot
	Version         string
}

func NewApp(config *types.Config, version string) *App {
	logger := loggerPkg.GetLogger(config.Log)
	templateManager := templates.NewTemplateManager()

	bot, err := tele.NewBot(tele.Settings{
		Token:  config.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, c tele.Context) {
			logger.Error().Err(err).Msg("Telebot error")
		},
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not start Telegram bot")
	}

	if len(config.Telegram.Admins) > 0 {
		logger.Debug().Msg("Using admins whitelist")
		bot.Use(middleware.Whitelist(config.Telegram.Admins...))
	}

	proxmoxManager := proxmox.NewManager(config, logger)

	return &App{
		Logger:          logger,
		ProxmoxManager:  proxmoxManager,
		TemplateManager: templateManager,
		Bot:             bot,
		Version:         version,
	}
}

func (a *App) Start() {
	a.Bot.Handle("/status", a.HandleStatus)
	a.Bot.Handle("/containers", a.HandleListContainers)
	a.Bot.Handle("/container", a.HandleContainerInfo)
	a.Bot.Handle("/node", a.HandleNodeInfo)
	a.Bot.Handle("/start", a.HandleStartContainer)
	a.Bot.Handle("/stop", a.HandleStopContainer)
	a.Bot.Handle("/restart", a.HandleRestartContainer)
	a.Bot.Handle("/disks", a.HandleListDisks)
	a.Bot.Handle("/help", a.HandleStartContainer)

	a.Bot.Handle(tele.OnCallback, a.HandleCallback)

	a.Logger.Info().Msg("Telegram bot listening")

	a.Bot.Start()
}

func (a *App) HandleCallback(c tele.Context) error {
	callback := c.Callback()

	callbacks := map[string]func(tele.Context, string) error{
		CallbackPrefixRestart:       a.HandleDoRestartContainer,
		CallbackPrefixCancelRestart: a.HandleDoCancelRestartContainer,
	}

	unique := strings.TrimSpace(callback.Data)

	for prefix, botCallback := range callbacks {
		if strings.HasPrefix(unique, prefix) {
			suffix, _ := strings.CutPrefix(unique, prefix)

			return botCallback(c, suffix)
		}
	}

	return nil
}

func (a *App) BotReply(c tele.Context, msg string, opts ...interface{}) error {
	msgsByNewline := strings.Split(msg, "\n")

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > MaxMessageSize {
			if err := c.Reply(sb.String(), tele.ModeHTML); err != nil {
				a.Logger.Error().Err(err).Msg("Could not send Telegram message")
				return err
			}

			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	opts = append(opts, tele.ModeHTML)

	if err := c.Reply(sb.String(), opts...); err != nil {
		a.Logger.Error().Err(err).Msg("Could not send Telegram message")
		return err
	}

	return nil
}
