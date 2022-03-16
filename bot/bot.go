package bot

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
	"github.com/andreapavoni/ttfm_bot/ttfm/commands"
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
	b.Bot.AddCommand("dj", commands.DjCommand())
	b.Bot.AddCommand("escort", commands.EscortCommand())
	b.Bot.AddCommand("escortme", commands.EscortMeCommand())
	b.Bot.AddCommand("autobop", commands.AutoBopCommand())
	b.Bot.AddCommand("autodj", commands.AutoDjCommand())
	b.Bot.AddCommand("autosnag", commands.AutoSnagCommand())
	b.Bot.AddCommand("bop", commands.BopCommand())
	b.Bot.AddCommand("dj", commands.DjCommand())
	b.Bot.AddCommand("snag", commands.SnagCommand())
	b.Bot.AddCommand("fan", commands.FanCommand())
	b.Bot.AddCommand("unfan", commands.UnfanCommand())
	b.Bot.AddCommand("props", commands.PropsCommand())
	b.Bot.AddCommand("skip", commands.SkipCommand())
	b.Bot.AddCommand("boot", commands.BootCommand())
	b.Bot.AddCommand("padd", commands.PlaylistCreateCommand())
	b.Bot.AddCommand("pdel", commands.PlaylistDeleteCommand())
	b.Bot.AddCommand("pls", commands.PlaylistListCommand())
	b.Bot.AddCommand("prm", commands.PlaylistRemoveSongCommand())
	b.Bot.AddCommand("pc", commands.PlaylistSwitchCommand())
	b.Bot.AddCommand("q", commands.QueueCommand())
	b.Bot.AddCommand("qadd", commands.QueueAddCommand())
	b.Bot.AddCommand("qrm", commands.QueueRemoveCommand())
	b.Bot.AddCommand("cfg", commands.SetConfigCommand())
	b.Bot.AddCommand("say", commands.SayCommand())
	b.Bot.AddCommand("help", commands.HelpCommand())
	b.Bot.AddCommand("r", commands.ReactionCommand())
	b.Bot.Start()
}
