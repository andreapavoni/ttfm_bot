package bot

import (
	"github.com/andreapavoni/ttfm_bot/commands"
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

type Bot struct {
	*ttfm.Bot
}

func New() *Bot {
	return &Bot{
		ttfm.New(),
	}
}

func (b *Bot) Start() {
	b.Bot.AddCommand("!dj", commands.DjCommandHandler)
	b.Bot.AddCommand("!escort", commands.EscortCommandHandler)
	b.Bot.AddCommand("!escortme", commands.EscortMeCommandHandler)
	b.Bot.AddCommand("!autobop", commands.AutoBopCommandHandler)
	b.Bot.AddCommand("!autodj", commands.AutoDjCommandHandler)
	b.Bot.AddCommand("!autosnag", commands.AutoSnagCommandHandler)
	b.Bot.AddCommand("!bop", commands.BopCommandHandler)
	b.Bot.AddCommand("!dj", commands.DjCommandHandler)
	b.Bot.AddCommand("!snag", commands.SnagCommandHandler)
	b.Bot.AddCommand("!fan", commands.FanCommandHandler)
	b.Bot.AddCommand("!unfan", commands.UnfanCommandHandler)
	b.Bot.AddCommand("!props", commands.PropsCommandHandler)
	b.Bot.AddCommand("!skip", commands.SkipCommandHandler)

	b.Bot.Start()
}
